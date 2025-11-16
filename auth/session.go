package auth

import (
	"log"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// SessionManager handles session lifecycle management
type SessionManager struct {
	service *Service
	logger  *log.Logger
}

// NewSessionManager creates a new session manager
func NewSessionManager(s *Service) *SessionManager {
	return &SessionManager{
		service: s,
		logger:  log.New(log.Writer(), "[session] ", log.LstdFlags),
	}
}

// CreateSession exchanges OAuth authorization code for session_id and user info
func (sm *SessionManager) CreateSession(authCode string) (map[string]interface{}, error) {
	sm.logger.Printf("Creating session for auth code: %s", authCode)
	
	return sm.service.callAuthServiceGeneric("/auth/session/create", map[string]string{
		"auth_code": authCode,
	})
}

// RefreshSession refreshes an existing session_id
func (sm *SessionManager) RefreshSession(sessionID string) (*models.AuthResponse, error) {
	sm.logger.Printf("Refreshing session: %s", sessionID)
	
	return sm.service.callAuthService("/auth/session/refresh", map[string]string{
		"session_id": sessionID,
	})
}

// ValidateSession checks if a session is still valid
func (sm *SessionManager) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	sm.logger.Printf("Validating session: %s", sessionID)
	
	return sm.service.callAuthService("/auth/session/validate", map[string]string{
		"session_id": sessionID,
	})
}

// ExpireSession marks a session as expired
func (sm *SessionManager) ExpireSession(sessionID string) error {
	sm.logger.Printf("Expiring session: %s", sessionID)
	
	// Call auth service to invalidate session
	_, err := sm.service.callAuthService("/auth/session/expire", map[string]string{
		"session_id": sessionID,
	})
	
	return err
}

// SessionInfo represents session information
type SessionInfo struct {
	SessionID   string    `json:"session_id"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	LastAccess  time.Time `json:"last_access"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	IsActive    bool      `json:"is_active"`
}

// GetSessionInfo retrieves detailed session information
func (sm *SessionManager) GetSessionInfo(sessionID string) (*SessionInfo, error) {
	authResp, err := sm.service.GetUserInfo(sessionID)
	if err != nil {
		return nil, err
	}

	return &SessionInfo{
		SessionID:  sessionID,
		UserID:     authResp.UserID,
		CreatedAt:  time.Now(),
		LastAccess: time.Now(),
		IsActive:   authResp.Success,
	}, nil
}