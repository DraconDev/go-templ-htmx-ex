package auth

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient handles HTTP communication
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

// Do executes an HTTP request
func (c *HTTPClient) Do(req *http.Request) (*http.Response, []byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return resp, bodyBytes, nil
}

// GetTimeout returns the configured timeout
func (c *HTTPClient) GetTimeout() time.Duration {
	return c.timeout
}