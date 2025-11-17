package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
)

// =============================================================================
// SESSION MANAGEMENT HANDLERS
// =============================================================================
// These handlers manage authentication sessions:
// - Set and validate session cookies
// - Handle user logout
// - Exchange authorization codes for sessions
// =============================================================================

// SetSessionHandler sets a new session cookie
func (h *AuthHandler) SetSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 SESSION: === Set session STARTED ===\n")
	fmt.Printf("🔐 SESSION: Content-Type: %s\n", r.Header.Get("Content-Type"))

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("🔐 SESSION: Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	if req.SessionID == "" {
		fmt.Printf("🔐 SESSION: Missing session_id\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing session_id",
		})
		return
	}

	fmt.Printf("🔐 SESSION: Session ID received, length: %d\n", len(req.SessionID))

	// Set session_id cookie (replaces session_token)
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    req.SessionID,
		Path:     "/",
		MaxAge:   2592000, // 30 days (server-side validation handles real security)
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	// Set session_id cookie
	http.SetCookie(w, sessionCookie)

	fmt.Printf("🔐 SESSION: Session ID cookie set successfully:")
	fmt.Printf("🔐 SESSION: - session_id cookie, Length: %d\n", len(sessionCookie.Value))

	fmt.Printf("🔐 SESSION: SUCCESS: Server session established")
	fmt.Printf("🔐 SESSION: === Set session COMPLETED ===\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Server session established successfully",
	})
}

// LogoutHandler handles user logout
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

// ExchangeCodeHandler exchanges OAuth authorization code for tokens
func (h *AuthHandler) ExchangeCodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔄 CODE: === Exchange authorization code STARTED ===\n")
	fmt.Printf("🔄 CODE: Request URL: %s\n", r.URL.String())

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		AuthCode string `json:"auth_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("🔄 CODE: Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	if req.AuthCode == "" {
		fmt.Printf("🔄 CODE: Missing authorization code\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing authorization code",
		})
		return
	}

	fmt.Printf("🔄 CODE: Authorization code received, length: %d\n", len(req.AuthCode))

	// Create session from authorization code (returns JSON with all info)
	fmt.Printf("🔄 CODE: Creating session from authorization code...\n")
	sessionData, err := h.AuthService.CreateSession(req.AuthCode)
	if err != nil {
		fmt.Printf("❌ CODE: Auth service failed: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("✅ CODE: Auth service returned session data: %+v\n", sessionData)

	// Extract session_id from the response
	var session_id string
	if sid, exists := sessionData["session_id"]; exists {
		if sidStr, ok := sid.(string); ok {
			session_id = sidStr
		}
	}

	if session_id == "" {
		fmt.Printf("❌ CODE: No session_id in auth response\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "No session_id received from auth service",
		})
		return
	}

	// Set session_id cookie for server sessions
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    session_id,
		Path:     "/",
		MaxAge:   2592000, // 30 days (server-side validation handles real security)
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	// Set session_id cookie
	http.SetCookie(w, sessionCookie)

	fmt.Printf("✅ CODE: Session token cookie set successfully (length: %d)\n", len(session_id))
	fmt.Printf("🔄 CODE: === Token exchange COMPLETED ===\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tokens exchanged successfully",
	})
}
