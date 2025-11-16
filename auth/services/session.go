package auth

import (
	"fmt"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// SessionService handles session validation
type SessionService struct {
	config   *config.Config
	http     *HTTPClient
	builder  *RequestBuilder
	parser   *ResponseParser
	timeout  time.Duration
}

// NewSessionService creates a new session service
func NewSessionService(cfg *config.Config) *SessionService {
	return &SessionService{
		config:  cfg,
		http:    NewHTTPClient(10 * time.Second),
		builder: NewRequestBuilder(cfg.AuthSecret),
		parser:  NewResponseParser(),
		timeout: 10 * time.Second,
	}
}

// ValidateSession validates a session token
func (s *SessionService) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/refresh", s.config.AuthServiceURL)
	params := map[string]string{
		"session_id": sessionID,
	}
	return s.CallAuthService(endpoint, params)
}

// CallAuthService makes a request to the auth microservice
func (s *SessionService) CallAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
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