package auth

import (
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/auth/http"
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Service handles communication with the auth microservice for session management
type Service struct {
	config  *config.Config
	http    *http.Client
	builder *Builder
	parser  *Parser
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	httpClient := http.NewClient()
	
	return &Service{
		config:  cfg,
		http:    httpClient,
		builder: NewBuilder(cfg.AuthSecret),
		parser:  NewParser(),
	}
}

// CallAuthService makes a request to the auth microservice
func (s *Service) CallAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	// Build request
	req, err := s.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	// Execute request
	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}

	// Parse response
	authResp, err := s.parser.ParseAuthResponse(bodyBytes)
	if err != nil {
		return nil, err
	}

	return authResp, nil
}

// ValidateSession validates a session_id with the auth service
func (s *Service) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/refresh", s.config.AuthServiceURL)
	params := map[string]string{
		"session_id": sessionID,
	}
	return s.CallAuthService(endpoint, params)
}

// GetUserInfo retrieves user information from auth service using session_id
func (s *Service) GetUserInfo(sessionID string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/userinfo", s.config.AuthServiceURL)
	params := map[string]string{
		"session_id": sessionID,
	}
	return s.CallAuthService(endpoint, params)
}

// Logout logs out a user using session_id
func (s *Service) Logout(sessionID string) error {
	// Log the logout for debugging purposes
	fmt.Printf("User logged out with session_id: %s\n", sessionID)
	return nil
}

// ValidateUser validates a user using session_id (alias for GetUserInfo)
func (s *Service) ValidateUser(sessionID string) (*models.AuthResponse, error) {
	return s.GetUserInfo(sessionID)
}

// ValidateToken validates a token (alias for ValidateSession) - kept for compatibility
func (s *Service) ValidateToken(sessionID string) (*models.AuthResponse, error) {
	return s.ValidateSession(sessionID)
}

// CreateSession creates a session from OAuth authorization code
func (s *Service) CreateSession(code string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/auth/session/create", s.config.AuthServiceURL)
	params := map[string]string{"code": code}

	// Build request
	req, err := s.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	// Execute request
	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}

	// Parse response as generic map
	response, err := s.parser.ParseGenericResponse(bodyBytes)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ExchangeCodeForTokens exchanges OAuth authorization code for session_id
func (s *Service) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/create", s.config.AuthServiceURL)
	params := map[string]string{
		"auth_code": code,
	}

	// Build request
	req, err := s.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to build request",
		}, err
	}

	// Execute request
	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to call auth service: " + err.Error(),
		}, err
	}

	// Parse response
	tokenResp, err := s.parser.ParseTokenExchangeResponse(bodyBytes)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to parse token exchange response: " + err.Error(),
		}, err
	}

	return tokenResp, nil
}