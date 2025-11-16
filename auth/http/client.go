package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/models"
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

// HTTPClient handles HTTP communication with the auth service
type HTTPClient struct {
	client  *http.Client
	timeout time.Duration
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client:  &http.Client{Timeout: timeout},
		timeout: timeout,
	}
}

// Do executes an HTTP request and returns response
func (c *HTTPClient) Do(req *Request) (*http.Response, []byte, error) {
	resp, err := c.client.Do(req.Request)
	if err != nil {
		return nil, nil, &AuthServiceError{"REQUEST_FAILED", "HTTP request failed", err}
	}
	defer resp.Body.Close()

	bodyBytes, err := req.readBody()
	if err != nil {
		return nil, nil, &AuthServiceError{"BODY_READ_FAILED", "Failed to read response body", err}
	}

	return resp, bodyBytes, nil
}

// Close closes the HTTP client
func (c *HTTPClient) Close() {
	// Nothing to close for default http.Client
}

// GetTimeout returns the configured timeout
func (c *HTTPClient) GetTimeout() time.Duration {
	return c.timeout
}

// ParseAuthResponse parses response into AuthResponse
func (c *HTTPClient) ParseAuthResponse(bodyBytes []byte) (*models.AuthResponse, error) {
	// This is a simple passthrough to avoid circular imports
	// In a real implementation, this would be handled by the parser package
	return nil, fmt.Errorf("use http.Parser for response parsing")
}