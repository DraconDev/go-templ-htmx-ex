package auth

import (
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Service handles communication with the auth microservice
type Service struct {
	config          *config.Config
	http            *HTTPClient
	builder         *RequestBuilder
	parser          *ResponseParser
	timeout         time.Duration
	sessionValidator *SessionValidator
	userInfoService  *UserInfoService
	tokenService     *TokenService
	compatibility    *CompatibilityService
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	timeout := 10 * time.Second
	httpClient := NewHTTPClient(timeout)
	builder := NewRequestBuilder(cfg.AuthSecret)
	parser := NewResponseParser()
	baseURL := cfg.AuthServiceURL

	// Create specialized services
	sessionValidator := NewSessionValidator(&Service{config: cfg, http: httpClient, builder: builder, parser: parser, timeout: timeout})
	userInfoService := NewUserInfoService(&Service{config: cfg, http: httpClient, builder: builder, parser: parser, timeout: timeout})
	tokenService := NewTokenService(sessionValidator)
	compatibility := NewCompatibilityService(&Service{config: cfg, http: httpClient, builder: builder, parser: parser, timeout: timeout})

	return &Service{
		config:          cfg,
		http:            httpClient,
		builder:         builder,
		parser:          parser,
		timeout:         timeout,
		sessionValidator: sessionValidator,
		userInfoService:  userInfoService,
		tokenService:     tokenService,
		compatibility:    compatibility,
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

// Delegate methods for compatibility with existing handlers
func (s *Service) ValidateUser(token string) (*models.AuthResponse, error) {
	return s.userInfoService.ValidateUser(token)
}

func (s *Service) ValidateToken(token string) (*models.AuthResponse, error) {
	return s.tokenService.ValidateToken(token)
}

func (s *Service) CreateSession(code string) (map[string]interface{}, error) {
	return s.sessionValidator.Service.CallAuthService("", nil) // Simplified for now
}

func (s *Service) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	return s.sessionValidator.ExchangeCodeForTokens(code)
}
