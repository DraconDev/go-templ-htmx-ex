package auth

import (
	"net/http"
	"time"
)

// HTTPClient handles HTTP communication with auth service
type HTTPClient struct {
	client  *http.Client
	baseURL string
	timeout time.Duration
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client:  &http.Client{Timeout: timeout},
		baseURL: baseURL,
		timeout: timeout,
	}
}

// Do executes an HTTP request
func (c *HTTPClient) Do(req *http.Request) (*http.Response, []byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes := make([]byte, 0)
	if resp.Body != nil {
		bodyBytes = make([]byte, 4096) // Reasonable buffer size
		n, _ := resp.Body.Read(bodyBytes)
		bodyBytes = bodyBytes[:n]
	}

	return resp, bodyBytes, nil
}