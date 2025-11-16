package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/auth/http"
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Service handles communication with the auth microservice for session management
type Service struct {
	config  *config.Config
	http    *http.Client
	timeout time.Duration
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	return &Service{
		config:  cfg,
		http:    http.NewClient(10 * time.Second),
		timeout: 10 * time.Second,
	}
}

// ValidateSession validates a session_id with the auth service
func (s *Service) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	return s.callAuthService("/auth/session/refresh", map[string]string{
		"session_id": sessionID,
	})
}

// GetUserInfo retrieves user information from auth service using session_id
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

// CreateSession creates a session from OAuth authorization code
func (s *Service) CreateSession(code string) (map[string]interface{}, error) {
	return s.callAuthServiceGeneric("/auth/session/create", map[string]string{
		"code": code,
	})
}

// ExchangeCodeForTokens exchanges OAuth authorization code for session_id
func (s *Service) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	return s.callAuthServiceToken("/auth/session/create", map[string]string{
		"auth_code": code,
	})
}

// callAuthService makes a request to the auth microservice
func (s *Service) callAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.config.AuthServiceURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.config.AuthSecret != "" {
		req.Header.Set("X-Auth-Secret", s.config.AuthSecret)
	}

	resp, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

	req, err := http.NewRequest("POST", s.config.AuthServiceURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.config.AuthSecret != "" {
		req.Header.Set("X-Auth-Secret", s.config.AuthSecret)
	}

	resp, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// callAuthServiceToken makes a request and returns token exchange response
func (s *Service) callAuthServiceToken(endpoint string, params map[string]string) (*models.TokenExchangeResponse, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return &models.TokenExchangeResponse{Success: false, Error: "Failed to marshal request"}, err
	}

	req, err := http.NewRequest("POST", s.config.AuthServiceURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return &models.TokenExchangeResponse{Success: false, Error: "Failed to create request"}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.config.AuthSecret != "" {
		req.Header.Set("X-Auth-Secret", s.config.AuthSecret)
	}

	resp, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return &models.TokenExchangeResponse{Success: false, Error: "Failed to call auth service"}, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return &models.TokenExchangeResponse{Success: false, Error: "Failed to parse response"}, err
	}

	// Extract session_id for session management
	var sessionID string
	if sid, exists := response["session_id"]; exists {
		if sidStr, ok := sid.(string); ok {
			sessionID = sidStr
		}
	}

	return &models.TokenExchangeResponse{
		Success: true,
		IdToken: sessionID,
	}, nil
}