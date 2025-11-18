package session

import (
	"encoding/json"
	"fmt"
	"net/http"

)

// ExchangeCodeHandler exchanges OAuth authorization code for tokens
// This handler is responsible ONLY for exchanging auth codes for session tokens
func (h *SessionHandler) ExchangeCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		AuthCode string `json:"auth_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleJSONError(w, "Invalid request body", err, nil) // Use nil for standard error for now
		return
	}

	if req.AuthCode == "" {
		handleJSONError(w, "Missing authorization code", nil, nil)
		return
	}

	// Create session from authorization code
	sessionData, err := h.AuthService.CreateSession(req.AuthCode)
	if err != nil {
		handleJSONError(w, err.Error(), err, nil)
		return
	}

	// Extract session_id from the response
	sessionID, err := extractSessionID(sessionData)
	if err != nil {
		handleJSONError(w, "No session_id received from auth service", nil, nil)
		return
	}

	// Use session utility to set the cookie
	sessionConfig := DefaultSessionCookieConfig()
	SetSessionCookie(w, sessionID, sessionConfig)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tokens exchanged successfully",
	})
}

// extractSessionID extracts session_id from session data map
func extractSessionID(sessionData map[string]interface{}) (string, error) {
	if sid, exists := sessionData["session_id"]; exists {
		if sidStr, ok := sid.(string); ok && sidStr != "" {
			return sidStr, nil
		}
	}
	return "", fmt.Errorf("no session_id in auth response")
}
