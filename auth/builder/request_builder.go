package builder

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// RequestBuilder builds HTTP requests for auth service
type RequestBuilder struct {
	authSecret string
}

// NewRequestBuilder creates a new request builder
func NewRequestBuilder(authSecret string) *RequestBuilder {
	return &RequestBuilder{
		authSecret: authSecret,
	}
}

// BuildPOSTRequest creates a POST request with JSON data
func (rb *RequestBuilder) BuildPOSTRequest(endpoint string, params map[string]string) (*http.Request, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if rb.authSecret != "" {
		req.Header.Set("X-Auth-Secret", rb.authSecret)
	}

	return req, nil
}

// AddAuthSecret adds auth secret to params if not present
func (rb *RequestBuilder) AddAuthSecret(params map[string]string) {
	if _, exists := params["secret"]; !exists && rb.authSecret != "" {
		params["secret"] = rb.authSecret
	}
}