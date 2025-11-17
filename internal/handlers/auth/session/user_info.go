package auth

import (
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers/auth"
	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
)

// GetUserInfo returns current user information for server-side rendering
// This function is responsible ONLY for extracting user info from session cookies
func (h *AuthHandler) GetUserInfo(r *http.Request) layouts.UserInfo {
	// Use session utility to get session cookie
	sessionID, err := GetSessionCookie(r)
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}
	}

	// Get user info from auth microservice
	userResp, err := h.AuthService.GetUserInfo(sessionID)
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}
	}

	return layouts.UserInfo{
		LoggedIn: userResp.Success,
		Name:     userResp.Name,
		Email:    userResp.Email,
		Picture:  userResp.Picture,
	}
}

// IsUserLoggedIn checks if a user is currently logged in
// This function is responsible ONLY for checking login status
func (h *auth.AuthHandler) IsUserLoggedIn(r *http.Request) bool {
	return IsSessionValid(r)
}
