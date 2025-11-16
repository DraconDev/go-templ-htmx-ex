package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// HTTPClient handles HTTP communication with the auth service
type HTTPClient struct {
	client   *http.Client
	timeout  time.Duration
	secret   string
	baseURL  string
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(timeout time.Duration, secret, baseURL string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{Timeout: timeout},
		timeout: timeout,
		secret:  secret,
		baseURL: baseURL,
	}
}

// RequestBuilder builds HTTP requests for the auth service
type RequestBuilder struct {
	secret string
}

// NewRequestBuilder creates a new request builder
func NewRequestBuilder(secret string) *RequestBuilder {
	return &RequestBuilder{
		secret: secret,
	}
}

// Request represents a built HTTP request
type Request struct {
	*http.Request
}

// BuildPOSTRequest creates a POST request with JSON data
func (rb *RequestBuilder) BuildPOSTRequest(endpoint string, params map[string]string) (*Request, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, NewAuthServiceError("INVALID_PARAMS", "Failed to marshal request data", err)
	}

	url := fmt.Sprintf("%s%s", endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, NewAuthServiceError("REQUEST_BUILD_FAILED", "Failed to create request", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if rb.secret != "" {
		req.Header.Set("X-Auth-Secret", rb.secret)
	}

	return &Request{req}, nil
}

// ResponseParser handles parsing responses from the auth service
type ResponseParser struct{}

// NewResponseParser creates a new response parser
func NewResponseParser() *ResponseParser {
	return &ResponseParser{}
}

// ParseAuthResponse parses a response as AuthResponse
func (rp *ResponseParser) ParseAuthResponse(bodyBytes []byte) (*models.AuthResponse, error) {
	var authResp models.AuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, NewAuthServiceError("PARSE_FAILED", "Failed to parse auth response", err)
	}

	return &authResp, nil
}

// ParseTokenExchangeResponse parses a response as TokenExchangeResponse
func (rp *ResponseParser) ParseTokenExchangeResponse(bodyBytes []byte) (*models.TokenExchangeResponse, error) {
	var respData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &respData); err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to parse response: " + err.Error(),
		}, NewAuthServiceError("PARSE_FAILED", "Failed to parse token exchange response", err)
	}

	// Check for errors in response
	if errMsg, hasError := respData["error"]; hasError {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   fmt.Sprintf("%v", errMsg),
		}, NewAuthServiceError("AUTH_SERVICE_ERROR", fmt.Sprintf("Auth service error: %v", errMsg), nil)
	}

	// Extract session_id
	var sessionID string
	if sessionInterface, exists := respData["session_id"]; exists {
		if sessionStr, ok := sessionInterface.(string); ok {
			sessionID = sessionStr
		}
	}

	if sessionID == "" {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Missing session_id in auth service response",
		}, NewAuthServiceError("MISSING_SESSION_ID", "Missing session_id in response", nil)
	}

	return &models.TokenExchangeResponse{
		Success: true,
		IdToken: sessionID,
	}, nil
}

// ParseGenericResponse parses a response as generic map
func (rp *ResponseParser) ParseGenericResponse(bodyBytes []byte) (map[string]interface{}, error) {
	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, NewAuthServiceError("PARSE_FAILED", "Failed to parse generic response", err)
	}

	return response, nil
}

// Do executes an HTTP request
func (c *HTTPClient) Do(req *Request) (*http.Response, []byte, error) {
	resp, err := c.client.Do(req.Request)
	if err != nil {
		return nil, nil, NewAuthServiceError("REQUEST_FAILED", "HTTP request failed", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, NewAuthServiceError("BODY_READ_FAILED", "Failed to read response body", err)
	}

	return resp, bodyBytes, nil
}