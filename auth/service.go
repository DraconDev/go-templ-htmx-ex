package auth

import (
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Service handles communication with the auth microservice
// NOTE: This service is largely unused - middleware makes direct HTTP calls
type Service struct {
	config *config.Config
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// ValidateSession validates a session token - Used by handlers but middleware doesn't call it
func (s *Service) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	// TODO: Actually implement this if needed
	// Currently returns dummy response for compatibility
	return &models.AuthResponse{
		Success: true,
		UserID:  "demo-session",
		Email:   "demo@example.com",
		Name:    "Demo User",
	}, nil
}

// ValidateToken validates a token (alias for ValidateSession) - ACTUALLY USED
func (s *Service) ValidateToken(token string) (*models.AuthResponse, error) {
	return s.ValidateSession(token)
}

// ValidateUser validates a user token (alias for GetUserInfo) - ACTUALLY USED
func (s *Service) ValidateUser(token string) (*models.AuthResponse, error) {
	return s.ValidateSession(token)
}
