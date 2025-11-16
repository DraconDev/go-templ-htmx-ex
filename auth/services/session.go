package services

import (
	"fmt"
	"time"

	httpclient "github.com/DraconDev/go-templ-htmx-ex/auth/http"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// SessionService handles session-related operations
type SessionService struct {
	client    *httpclient.Client
	authURL   string
}

// NewSessionService creates a new session service
func NewSessionService(authURL string, timeout time.Duration) *SessionService {
	return &SessionService{
		client:  httpclient.NewClient(timeout),
		authURL: authURL,
	}
}

// ValidateSession validates a session token
func (s *SessionService) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// ExchangeCodeForTokens exchanges OAuth authorization code for server session
func (s *SessionService) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}