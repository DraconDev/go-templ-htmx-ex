package auth

import (
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Service handles communication with the auth microservice
type Service struct {
	config *config.Config
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// ValidateUser validates a user token (ACTUALLY USED by handlers)
func (s *Service) ValidateUser(token string) (*models.AuthResponse, error) {
	fmt.Printf("🔐 SERVICE: ValidateUser called with token: %s\n", token[:8]+"...")
	
	// TODO: Implement actual validation with auth microservice
	// Currently returns dummy response for compatibility
	return &models.AuthResponse{
		Success: true,
		UserID:  "demo-user",
		Email:   "demo@example.com",
		Name:    "Demo User",
	}, nil
}

// ValidateToken validates a token (ACTUALLY USED by ValidateSessionHandler)
func (s *Service) ValidateToken(token string) (*models.AuthResponse, error) {
	return s.ValidateUser(token)
}

// CreateSession creates a session from authorization code (ACTUALLY USED by /api/auth/test-session-create)
func (s *Service) CreateSession(code string) (map[string]interface{}, error) {
	fmt.Printf("🧪 SERVICE: CreateSession called with code: %s\n", code[:8]+"...")
	
	// TODO: Implement actual session creation
	return map[string]interface{}{
		"session_id": "demo-session-" + code,
		"user": map[string]interface{}{
			"id":    "demo-user",
			"email": "demo@example.com",
			"name":  "Demo User",
		},
	}, nil
}

// ExchangeCodeForTokens exchanges OAuth authorization code for server session (ACTUALLY USED by /api/auth/exchange-code)
func (s *Service) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	fmt.Printf("🔄 SERVICE: ExchangeCodeForTokens called with code: %s\n", code[:8]+"...")
	
	// TODO: Implement actual code exchange
	return &models.TokenExchangeResponse{
		Success: true,
		IdToken: "demo-session-" + code,
	}, nil
}
