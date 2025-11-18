package session

import (
	"encoding/json"
	"net/http"
)

// LogoutHandler handles user logout
// This handler is responsible ONLY for clearing session cookies
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Use session utility to clear the cookie
	sessionConfig := DefaultSessionCookieConfig()
	ClearSessionCookie(w, sessionConfig)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}
