package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/config"
	"github.com/DraconDev/go-templ-htmx-ex/internal/models"
)

// AuthService handles session management with the auth microservice
type AuthService struct {
	config  *config.Config
	client  *http.Client
	timeout time.Duration
}

// NewAuthService creates a new auth service instance
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		config:  cfg,
		client:  &http.Client{Timeout: 10 * time.Second},
		timeout: 10 * time.Second,
	}
}

// CreateSession exchanges OAuth authorization code for session_id and user info
func (s *AuthService) CreateSession(auth_code string) (map[string]interface{}, error) {
	return s.callAuthServiceGeneric("/auth/session/create", map[string]string{
		"auth_code": auth_code,
	})
}

// RefreshSession refreshes an existing session_id
func (s *AuthService) RefreshSession(session_id string) (*models.AuthResponse, error) {
	return s.callAuthService("/auth/session/refresh", map[string]string{
		"session_id": session_id,
	})
}

// GetUserInfo retrieves user information using session_id
func (s *AuthService) GetUserInfo(session_id string) (*models.AuthResponse, error) {
	return s.callAuthService("/auth/userinfo", map[string]string{
		"session_id": session_id,
	})
// ValidateSession validates a session_id and returns user information for middleware
func (s *AuthService) ValidateSession(session_id string) (*models.AuthResponse, error) {
	return s.callAuthService("/auth/session/refresh", map[string]string{
		"session_id": session_id,
	})
}
}

// Logout logs out a user using session_id
func (s *AuthService) Logout(session_id string) error {
	fmt.Printf("User logged out with session_id: %s\n", session_id)
	return nil
}

// callAuthService makes a request to the auth microservice
func (s *AuthService) callAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	bodyBytes, err := s.makeRequest(endpoint, params)
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
func (s *AuthService) callAuthServiceGeneric(endpoint string, params map[string]string) (map[string]interface{}, error) {
	bodyBytes, err := s.makeRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// makeRequest handles the HTTP request to auth microservice
func (s *AuthService) makeRequest(endpoint string, params map[string]string) ([]byte, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.config.AuthServiceURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}