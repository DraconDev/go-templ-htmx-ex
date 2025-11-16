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
	return &Service{
		config:  cfg,
		http:    NewHTTPClient(10 * time.Second),
		builder: NewRequestBuilder(cfg.AuthSecret),
		parser:  NewResponseParser(),
		timeout: 10 * time.Second,
	}
}

// CallAuthService makes a request to the auth microservice
func (s *Service) CallAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
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
