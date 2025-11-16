package auth

import (
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// SessionValidator handles session validation operations
type SessionValidator struct {
	Service *Service
}

// NewSessionValidator creates a new session validator
func NewSessionValidator(service *Service) *SessionValidator {
	return &SessionValidator{Service: service}
}

// ValidateSession validates a session token
func (sv *SessionValidator) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/refresh", sv.Service.config.AuthServiceURL)
	params := map[string]string{
		"session_id": sessionID,
	}
	return sv.Service.CallAuthService(endpoint, params)
}

// ExchangeCodeForTokens exchanges OAuth authorization code for server session
func (sv *SessionValidator) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/create", sv.Service.config.AuthServiceURL)
	params := map[string]string{
		"auth_code": code,
	}

	req, err := sv.Service.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to build request",
		}, err
	}

	_, bodyBytes, err := sv.Service.http.Do(req)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to call auth service: " + err.Error(),
		}, err
	}

	return sv.Service.parser.ParseTokenExchangeResponse(bodyBytes)
}