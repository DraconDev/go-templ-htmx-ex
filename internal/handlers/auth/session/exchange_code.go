package session

import (
	"encoding/json"
	"fmt"
	"net/http"

)

// ExchangeCodeHandler exchanges OAuth authorization code for tokens
// This handler is responsible ONLY for exchanging auth codes for session tokens
func (h *SessionHandler) ExchangeCodeHandler(w http.ResponseWriter, r *http.Request) {
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

	// Exchange code for tokens via auth service (using the working reference logic)
	fmt.Printf("🔄 CODE: Calling auth service to exchange code for tokens...\n")
	authResp, err := h.AuthService.ExchangeCodeForTokens(req.AuthCode)
	if err != nil {
		fmt.Printf("❌ CODE: Auth service failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if !authResp.Success {
		fmt.Printf("❌ CODE: Auth service returned failure: %s\n", authResp.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": authResp.Error,
		})
		return
	}

	fmt.Printf("✅ CODE: Auth service returned success: %v\n", authResp.Success)
	fmt.Printf("🔄 CODE: Auth response: %+v\n", authResp)

	// Set session_id cookie for server sessions (same as reference)
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    authResp.UserID, // Using UserID as the session identifier
		Path:     "/",
		MaxAge:   2592000, // 30 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	http.SetCookie(w, sessionCookie)

	fmt.Printf("✅ CODE: Session token cookie set successfully (length: %d)\n", len(authResp.UserID))
	fmt.Printf("🔄 CODE: === Token exchange COMPLETED ===\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tokens exchanged successfully",
	})
}
