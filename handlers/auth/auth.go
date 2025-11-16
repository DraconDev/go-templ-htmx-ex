package auth

import (
	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	Config      *config.Config // App configuration
	AuthService *auth.Service  // Auth service for session management (includes HTTP client)
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(config *config.Config) *AuthHandler {
	return &AuthHandler{
		Config:      config,
		AuthService: auth.NewService(config),
	}
}