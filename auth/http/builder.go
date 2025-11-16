package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AuthServiceError represents a basic auth service error for HTTP package
type AuthServiceError struct {
	Code    string
	Message string
	Err     error
}

func (e *AuthServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("auth service error: %s - %v", e.Message, e.Err)
	}
	return fmt.Sprintf("auth service error: %s", e.Message)
}

// Request represents a built HTTP request
type Request struct {
	*http.Request
}

// readBody reads the response body
func (r *Request) readBody() ([]byte, error) {
	return io.ReadAll(r.Body)
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

// BuildPOSTRequest creates a POST request with JSON data
func (rb *RequestBuilder) BuildPOSTRequest(endpoint string, params map[string]string) (*Request, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, &AuthServiceError{"INVALID_PARAMS", "Failed to marshal request data", err}
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &AuthServiceError{"REQUEST_BUILD_FAILED", "Failed to create request", err}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if rb.secret != "" {
		req.Header.Set("X-Auth-Secret", rb.secret)
	}

	return &Request{req}, nil
}

// GetSecret returns the configured secret
func (rb *RequestBuilder) GetSecret() string {
	return rb.secret
}