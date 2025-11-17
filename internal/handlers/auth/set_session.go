package auth

import (
	"encoding/json"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/errors"
	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers/auth/session_utils"
)

// SetSessionHandler handles setting a new session cookie
// This handler is responsible ONLY for setting session cookies from a provided session ID
func (h *AuthHandler) SetSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleJSONError(w, "Invalid request body", err, errors.NewBadRequestError)
		return
	}

	if req.SessionID == "" {
		handleJSONError(w, "Missing session_id", nil, errors.NewBadRequestError)
		return
	}

	// Use session utility to set the cookie
	sessionConfig := session.DefaultSessionCookieConfig()
	session.SetSessionCookie(w, req.SessionID, sessionConfig)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Server session established successfully",
	})
}

// handleJSONError is a helper to standardize error responses
func handleJSONError(w http.ResponseWriter, message string, err error, errorType func(string) *errors.AppError) {
	if err != nil {
		// Log the error for debugging
		// This would typically use a proper logger
	}
	w.WriteHeader(errorType(message).Code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
	})
}