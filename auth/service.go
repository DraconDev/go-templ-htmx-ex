package auth

import (
	"fmt"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Service handles communication with the auth microservice
type Service struct {
	config   *config.Config
	http     *HTTPClient
	builder  *RequestBuilder
	parser   *ResponseParser
	timeout  time.Duration
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	timeout := 10 * time.Second
	httpClient := NewHTTPClient(timeout, cfg.AuthSecret, cfg.AuthServiceURL)
	
	return &Service{
		config:  cfg,
		http:    httpClient,
		builder: NewRequestBuilder(cfg.AuthSecret),
		parser:  NewResponseParser(),
		timeout: timeout,
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

	// Validate response
	if !authResp.Success {
		return authResp, NewAuthServiceError("AUTH_FAILED", "Authentication failed", nil)
	}

	return authResp, nil
}

// ValidateSession validates a session token
func (s *Service) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/refresh", s.config.AuthServiceURL)
	params := map[string]string{
		"session_id": sessionID,
	}
	return s.CallAuthService(endpoint, params)
}

// GetUserInfo retrieves user information from auth service
func (s *Service) GetUserInfo(token string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/userinfo", s.config.AuthServiceURL)
	params := map[string]string{
		"token": token,
	}
	return s.CallAuthService(endpoint, params)
}

// Logout logs out a user (通知auth service)
func (s *Service) Logout(token string) error {
	// Since this is a server session system, we log it
	// In a more complex system, you might want to blacklist the token
	fmt.Printf("User logged out with token: %s\n", token)
	return nil
}

// ValidateUser validates a user token (alias for GetUserInfo)
func (s *Service) ValidateUser(token string) (*models.AuthResponse, error) {
	return s.GetUserInfo(token)
}

// ValidateToken validates a token (alias for ValidateSession)
func (s *Service) ValidateToken(token string) (*models.AuthResponse, error) {
	return s.ValidateSession(token)
}

// CreateSession exchanges OAuth authorization code for session creation
// This is a test function to see what /session/create returns
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

// ExchangeCodeForTokens exchanges OAuth authorization code for server session
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
