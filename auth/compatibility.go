package auth

import (
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
)

// CompatibilityService provides backward compatibility methods
type CompatibilityService struct {
	config  *config.Config
	http    *HTTPClient
	builder *RequestBuilder
	parser  *ResponseParser
	timeout time.Duration
}

// NewCompatibilityService creates a new compatibility service
func NewCompatibilityService(cfg *config.Config) *CompatibilityService {
	timeout := 10 * time.Second
	httpClient := NewHTTPClient(timeout)
	
	return &CompatibilityService{
		config:  cfg,
		http:    httpClient,
		builder: NewRequestBuilder(cfg.AuthSecret),
		parser:  NewResponseParser(),
		timeout: timeout,
	}
}

// CallAuthService makes a request to the auth microservice
func (s *CompatibilityService) CallAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
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

// Logout logs out a user
func (s *CompatibilityService) Logout(token string) error {
	return nil
}