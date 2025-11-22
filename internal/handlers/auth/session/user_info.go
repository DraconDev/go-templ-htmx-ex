package session

import (
	"net/http"
)

// GetUserInfo retrieves user information from the session cookie
// Returns UserInfo for template rendering
func (h *SessionHandler) GetUserInfo(r *http.Request) UserInfo {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return UserInfo{LoggedIn: false}
	}

	sessionID := cookie.Value
	if sessionID == "" {
		return UserInfo{LoggedIn: false}
	}

	// Get user context from auth service (via session refresh)
	userContext, err := h.AuthService.GetUserInfo(sessionID)
	if err != nil {
		return UserInfo{LoggedIn: false}
	}

	return UserInfo{
		LoggedIn: true,
		Name:     userContext.Name,
		Email:    userContext.Email,
		Picture:  userContext.Picture,
	}
}

// IsUserLoggedIn checks if a user is currently logged in
// This function is responsible ONLY for checking login status
func (h *SessionHandler) IsUserLoggedIn(r *http.Request) bool {
	return IsSessionValid(r)
}
