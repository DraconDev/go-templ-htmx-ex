package http

import (
	"net/http"
	"time"
)

// Client handles HTTP communication with auth service
type Client struct {
	timeout time.Duration
}

// NewClient creates a new HTTP client
func NewClient(timeout time.Duration) *Client {
	return &Client{
		timeout: timeout,
	}
}

// Do executes an HTTP request
func (c *Client) Do(req *http.Request) (*http.Response, []byte, error) {
	client := &http.Client{Timeout: c.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Read response body - placeholder for now
	bodyBytes := make([]byte, 0)
	if resp.Body != nil {
		bodyBytes = make([]byte, 1024)
	}

	return resp, bodyBytes, nil
}

// GetTimeout returns the configured timeout
func (c *Client) GetTimeout() time.Duration {
	return c.timeout
}
