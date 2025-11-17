package auth

import (
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
)

// =============================================================================
// SESSION MANAGEMENT HANDLERS
// =============================================================================
// These handlers manage authentication sessions:
// - Get user information from session cookies
// =============================================================================




// GetUserInfo returns current user information for server-side rendering
func (h *AuthHandler) GetUserInfo(r *http.Request) layouts.UserInfo {
	// Get session token from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}
	}

	// Get user info from auth microservice
	userResp, err := h.AuthService.GetUserInfo(cookie.Value)
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
