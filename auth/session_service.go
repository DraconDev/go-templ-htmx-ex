package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
)

// SessionServiceImpl implements session-related operations
type SessionServiceImpl struct {
	http     *HTTPClient
	builder  *RequestBuilder
	parser   *ResponseParser
	baseURL  string
	timeout  time.Duration
}

// NewSessionServiceImpl creates a new session service
func NewSessionServiceImpl(httpClient *HTTPClient, builder *RequestBuilder, parser *ResponseParser, baseURL string, timeout time.Duration) *SessionServiceImpl {
	return &SessionServiceImpl{
		http:     httpClient,
		builder:  builder,
		parser:   parser,
		baseURL:  baseURL,
		timeout:  timeout,
	}
}

// ValidateSession validates a session token
func (s *SessionServiceImpl) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/refresh", s.baseURL)
	params := map[string]string{
		"session_id": sessionID,
	}
	return s.makeRequest(endpoint, params)
}

// CreateSession creates a session from authorization code
func (s *SessionServiceImpl) CreateSession(code string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/auth/session/create", s.baseURL)
	params := map[string]string{"code": code}

	req, err := s.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}

	return s.parser.ParseGenericResponse(bodyBytes)
}

// ExchangeCodeForTokens exchanges OAuth authorization code for server session
func (s *SessionServiceImpl) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/create", s.baseURL)
	params := map[string]string{
		"auth_code": code,
	}

	req, err := s.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to build request",
		}, err
	}

	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to call auth service: " + err.Error(),
		}, err
	}

	tokenResp, err := s.parser.ParseTokenExchangeResponse(bodyBytes)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to parse token exchange response: " + err.Error(),
		}, err
	}

	return tokenResp, nil
}

// makeRequest is a helper method for session operations
func (s *SessionServiceImpl) makeRequest(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	req, err := s.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}

	authResp, err := s.parser.ParseAuthResponse(bodyBytes)
	if err != nil {
		return nil, err
	}

	return authResp, s.parser.ValidateResponseSuccess(authResp)
}