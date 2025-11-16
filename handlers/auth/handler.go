package auth

import (
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
)

// HTTPClient shared HTTP client for auth handlers
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a shared HTTP client
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  false,
				DisableKeepAlives:   false,
				MaxConnsPerHost:     10,
				ForceAttemptHTTP2:   true,
			},
		},
	}
}

// Do executes HTTP requests
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	Config      *config.Config // App configuration
	AuthService *auth.Service  // Auth service for session management
	HTTPClient  *HTTPClient    // Shared HTTP client
}

// NewAuthHandler creates a new authentication handler with shared HTTP client
func NewAuthHandler(config *config.Config) *AuthHandler {
	return &AuthHandler{
		Config:      config,
		AuthService: auth.NewService(config),
		HTTPClient:  NewHTTPClient(),
	}
}