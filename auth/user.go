package auth

import (
	"log"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// UserManager handles user-related operations
type UserManager struct {
	service *Service
	logger  *log.Logger
}

// NewUserManager creates a new user manager
func NewUserManager(s *Service) *UserManager {
	return &UserManager{
		service: s,
		logger:  log.New(log.Writer(), "[user] ", log.LstdFlags),
	}
}

// GetUserInfo retrieves user information using session_id
func (um *UserManager) GetUserInfo(sessionID string) (*models.AuthResponse, error) {
	um.logger.Printf("Getting user info for session: %s", sessionID)
	
	return um.service.callAuthService("/auth/userinfo", map[string]string{
		"session_id": sessionID,
	})
}

// GetUserProfile retrieves detailed user profile
func (um *UserManager) GetUserProfile(sessionID string) (*UserProfile, error) {
	authResp, err := um.GetUserInfo(sessionID)
	if err != nil {
		return nil, err
	}

	if !authResp.Success {
		return nil, &AuthError{
			Code:    "USER_NOT_FOUND",
			Message: "User not found or session invalid",
		}
	}

	return &UserProfile{
		UserID:  authResp.UserID,
		Email:   authResp.Email,
		Name:    authResp.Name,
		Picture: authResp.Picture,
	}, nil
}

// UpdateUserProfile updates user profile information
func (um *UserManager) UpdateUserProfile(sessionID string, updates UserProfileUpdates) error {
	um.logger.Printf("Updating user profile for session: %s", sessionID)
	
	_, err := um.service.callAuthService("/auth/user/update", map[string]string{
		"session_id": sessionID,
		"name":       updates.Name,
		"picture":    updates.Picture,
	})
	
	return err
}

// UserProfile represents user profile information
type UserProfile struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// UserProfileUpdates represents updatable profile fields
type UserProfileUpdates struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// HasRole checks if user has specific role
func (up *UserProfile) HasRole(role string) bool {
	// This would typically check against a roles list
	// For now, returning based on email domain or name patterns
	switch role {
	case "admin":
		return up.Email == "admin@example.com"
	case "user":
		return true
	default:
		return false
	}
}