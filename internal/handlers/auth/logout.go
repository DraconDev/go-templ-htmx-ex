package auth

import (
	"encoding/json"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers/auth/session_utils"
)

// LogoutHandler handles user logout
// This handler is responsible ONLY for clearing session cookies
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Use session utility to clear the cookie
	sessionConfig := session.DefaultSessionCookieConfig()
	session.ClearSessionCookie(w, sessionConfig)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}