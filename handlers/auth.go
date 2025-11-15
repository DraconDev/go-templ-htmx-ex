// RefreshSessionHandler refreshes an existing session
func (h *AuthHandler) RefreshSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔄 REFRESH: === Session refresh STARTED ===\n")

	var req struct {
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("🔄 REFRESH: Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	if req.SessionID == "" {
		fmt.Printf("🔄 REFRESH: Missing session ID\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing session ID",
		})
		return
	}

	fmt.Printf("🔄 REFRESH: Session ID received, length: %d\n", len(req.SessionID))

	// Refresh session with auth service
	fmt.Printf("🔄 REFRESH: Calling auth service to refresh session...\n")
	refreshResp, err := h.AuthService.RefreshSession(req.SessionID)
	if err != nil {
		fmt.Printf("❌ REFRESH: Auth service failed: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// If refresh is successful and returns new session info, update cookie
	if refreshResp.Success && refreshResp.SessionID != "" {
		fmt.Printf("✅ REFRESH: Session refreshed successfully, new session ID: %s\n", refreshResp.SessionID[:8]+"...")
		
		// Update session cookie with new session ID
		sessionCookie := &http.Cookie{
			Name:     "session_id",
			Value:    refreshResp.SessionID,
			Path:     "/",
			MaxAge:   3600, // 1 hour
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
		}

		http.SetCookie(w, sessionCookie)
		fmt.Printf("✅ REFRESH: Updated session cookie successfully\n")
	}

	fmt.Printf("🔄 REFRESH: === Session refresh COMPLETED ===\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"session_id":   refreshResp.SessionID,
		"user_context": refreshResp.UserContext,
		"message":      "Session refreshed successfully",
	})
}
