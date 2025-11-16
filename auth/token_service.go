package auth

import (
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// TokenServiceImpl implements token-related operations
type TokenServiceImpl struct {
	sessionService SessionService
}

// NewTokenServiceImpl creates a new token service
func NewTokenServiceImpl(sessionService SessionService) *TokenServiceImpl {
	return &TokenServiceImpl{
		sessionService: sessionService,
	}
}

// ValidateToken validates a token (alias for ValidateSession)
func (t *TokenServiceImpl) ValidateToken(token string) (*models.AuthResponse, error) {
	return t.sessionService.ValidateSession(token)
}