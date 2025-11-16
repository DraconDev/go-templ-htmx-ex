package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// HTTPClient handles HTTP communication with auth service
type HTTPClient struct {
	client   *http.Client
	baseURL  string
	timeout  time.Duration
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

// Service handles session management with the auth microservice
type Service struct {
	config  *config.Config
	http    *HTTPClient
	timeout time.Duration
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	return &Service{
		config:  cfg,
		http:    NewHTTPClient(cfg.AuthServiceURL, 10*time.Second),
		timeout: 10 * time.Second,
	}
}

// CreateSession exchanges OAuth authorization code for session_id and user info
func (s *Service) CreateSession(authCode string) (map[string]interface{}, error) {
	return s.callAuthServiceGeneric("/auth/session/create", map[string]string{
		"auth_code": authCode,
	})
}

// RefreshSession refreshes an existing session_id
func (s *Service) RefreshSession(sessionID string) (*models.AuthResponse, error) {
	return s.callAuthService("/auth/session/refresh", map[string]string{
		"session_id": sessionID,
	})
}

// GetUserInfo retrieves user information using session_id
func (s *Service) GetUserInfo(sessionID string) (*models.AuthResponse, error) {
	return s.callAuthService("/auth/userinfo", map[string]string{
		"session_id": sessionID,
	})
}

// Logout logs out a user using session_id
func (s *Service) Logout(sessionID string) error {
	fmt.Printf("User logged out with session_id: %s\n", sessionID)
	return nil
}

// callAuthService makes a request to the auth microservice
func (s *Service) callAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.http.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}

	var authResp models.AuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

// callAuthServiceGeneric makes a request and returns generic response
func (s *Service) callAuthServiceGeneric(endpoint string, params map[string]string) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.http.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}
