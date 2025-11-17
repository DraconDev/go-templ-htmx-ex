package auth

import (
	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/config"
	"github.com/DraconDev/go-templ-htmx-ex/internal/services"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	Config      *config.Config // App configuration
	AuthService *services.AuthService  // Auth service for session management (includes HTTP client)
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(config *config.Config) *AuthHandler {
	return &AuthHandler{
		Config:      config,
		AuthService: services.NewAuthService(config),
	}
}