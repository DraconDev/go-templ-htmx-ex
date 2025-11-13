// AutoRefreshMiddleware detects expired tokens and automatically refreshes them
// This provides seamless token refresh without user intervention
func AutoRefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user has a session token
		cookie, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Check if token needs refresh (expires within 5 minutes)
		if middleware.CheckIfTokenNeedsRefresh(cookie.Value) {
			fmt.Printf("🔄 AUTO-REFRESH: Token expiring soon, attempting refresh...\n")
			
			// Automatically refresh the token
			refreshSuccess := attemptAutoTokenRefresh(w, r)
			if refreshSuccess {
				fmt.Printf("🔄 AUTO-REFRESH: ✅ Token refreshed automatically!\n")
				// Token has been refreshed, continue with the request
				next.ServeHTTP(w, r)
				return
			} else {
				fmt.Printf("🔄 AUTO-REFRESH: ❌ Auto-refresh failed\n")
				// Redirect to login if refresh fails
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}

		// Token is valid, proceed normally
		next.ServeHTTP(w, r)
	})
}

// attemptAutoTokenRefresh attempts to refresh the session token automatically
func attemptAutoTokenRefresh(w http.ResponseWriter, r *http.Request) bool {
	// Get refresh token cookie
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Printf("🔄 AUTO-REFRESH: No refresh token available\n")
		return false
	}

	// Call auth service to get new tokens
	cfg := config.LoadConfig()
	authService := auth.NewService(cfg)
	refreshResp, err := authService.RefreshToken(refreshCookie.Value)
	
	if err != nil || !refreshResp.Success || refreshResp.Token == "" {
		fmt.Printf("🔄 AUTO-REFRESH: Auth service failed\n")
		return false
	}

	// Set new session token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    refreshResp.Token,
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production
	})

	fmt.Printf("🔄 AUTO-REFRESH: New token set: %d chars\n", len(refreshResp.Token))
	return true
}

// CheckIfTokenNeedsRefresh checks if a token expires within the specified timeframe
func CheckIfTokenNeedsRefresh(token string) bool {
	if token == "" {
		return true
	}

	// Parse JWT to get expiration
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return true
	}

	// Decode payload
	payload, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return true
	}

	// Parse expiration
	var claims struct {
		Exp int64 `json:"exp"`
	}
	
	if err := json.Unmarshal(payload, &claims); err != nil {
		return true
	}

	// Check if token expires within next 5 minutes
	now := time.Now().Unix()
	refreshThreshold := now + (5 * 60) // 5 minutes

	return claims.Exp < refreshThreshold
}