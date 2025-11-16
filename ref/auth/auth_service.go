package auth

import (
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// AuthService represents the main auth service that orchestrates all auth operations
type AuthService struct {
	sessionService SessionService
	userService    UserService
	tokenService   TokenService
}

// NewAuthService creates a new auth service with all dependencies
func NewAuthService(cfg *config.Config) *AuthService {
	timeout := 10 * time.Second
	httpClient := NewHTTPClient(timeout)
	builder := NewRequestBuilder(cfg.AuthSecret)
	parser := NewResponseParser()
	baseURL := cfg.AuthServiceURL

	// Create individual services
	sessionSvc := NewSessionServiceImpl(httpClient, builder, parser, baseURL, timeout)
	userSvc := NewUserServiceImpl(httpClient, builder, parser, baseURL)
	tokenSvc := NewTokenServiceImpl(sessionSvc)

	return &AuthService{
		sessionService: sessionSvc,
		userService:    userSvc,
		tokenService:   tokenSvc,
	}
}

// Implement SessionService interface
func (a *AuthService) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	return a.sessionService.ValidateSession(sessionID)
}

func (a *AuthService) CreateSession(code string) (map[string]interface{}, error) {
	return a.sessionService.CreateSession(code)
}

func (a *AuthService) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	return a.sessionService.ExchangeCodeForTokens(code)
}

// Implement UserService interface
func (a *AuthService) GetUserInfo(token string) (*models.AuthResponse, error) {
	return a.userService.GetUserInfo(token)
}

func (a *AuthService) ValidateUser(token string) (*models.AuthResponse, error) {
	return a.userService.ValidateUser(token)
}

func (a *AuthService) Logout(token string) error {
	return a.userService.Logout(token)
}

// Implement TokenService interface
func (a *AuthService) ValidateToken(token string) (*models.AuthResponse, error) {
	return a.tokenService.ValidateToken(token)
}

// Convenience methods for backward compatibility
func (a *AuthService) CallAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	// Delegate to session service for this specific endpoint
	return a.sessionService.ValidateSession("")
}