package auth

import (
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// SessionService handles session-related operations
type SessionService interface {
	ValidateSession(sessionID string) (*models.AuthResponse, error)
	CreateSession(code string) (map[string]interface{}, error)
	ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error)
}

// UserService handles user-related operations
type UserService interface {
	GetUserInfo(token string) (*models.AuthResponse, error)
	ValidateUser(token string) (*models.AuthResponse, error)
	Logout(token string) error
}

// TokenService handles token-related operations
type TokenService interface {
	ValidateToken(token string) (*models.AuthResponse, error)
}

// AuthOrchestrator coordinates all auth operations
type AuthOrchestrator interface {
	SessionService
	UserService
	TokenService
}