package auth

import (
	"log"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// AuthClient handles communication with auth service
type AuthClient struct {
	service *Service
	logger  *log.Logger
}

// NewAuthClient creates a new auth client
func NewAuthClient(s *Service) *AuthClient {
	return &AuthClient{
		service: s,
		logger:  log.New(log.Writer(), "[auth] ", log.LstdFlags),
	}
}

// Logout logs out a user using session_id
func (ac *AuthClient) Logout(sessionID string) error {
	ac.logger.Printf("Logging out session: %s", sessionID)
	
	// Call auth service to properly invalidate session
	_, err := ac.service.callAuthServiceGeneric("/auth/logout", map[string]string{
		"session_id": sessionID,
	})
	
	if err != nil {
		ac.logger.Printf("Error logging out session %s: %v", sessionID, err)
		return err
	}
	
	ac.logger.Printf("Successfully logged out session: %s", sessionID)
	return nil
}

// CheckHealth checks auth service health
func (ac *AuthClient) CheckHealth() error {
	ac.logger.Println("Checking auth service health")
	
	_, err := ac.service.callAuthServiceGeneric("/health", map[string]string{})
	return err
}

// AuthError represents authentication errors
type AuthError struct {
	Code    string
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}