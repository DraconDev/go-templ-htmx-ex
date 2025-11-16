package auth

import (
	"net/http"
	"time"
)

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
		return nil, nil, NewAuthServiceError("REQUEST_FAILED", "HTTP request failed", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := req.readBody()
	if err != nil {
		return nil, nil, NewAuthServiceError("BODY_READ_FAILED", "Failed to read response body", err)
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