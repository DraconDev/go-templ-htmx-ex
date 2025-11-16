package http

import (
	"net/http"
	"time"
)

// Client handles HTTP communication
type Client struct {
	client  *http.Client
	timeout time.Duration
}

// NewClient creates a new HTTP client
func NewClient(timeout time.Duration) *Client {
	return &Client{
		client:  &http.Client{Timeout: timeout},
		timeout: timeout,
	}
}

// Do executes an HTTP request
func (c *Client) Do(req *http.Request) (*http.Response, []byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	bodyBytes := make([]byte, 0)
	return resp, bodyBytes, nil
}

// GetTimeout returns the configured timeout
func (c *Client) GetTimeout() time.Duration {
	return c.timeout
}
