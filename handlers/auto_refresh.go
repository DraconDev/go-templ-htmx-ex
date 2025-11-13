package handlers

import (
	"fmt"
	"net/http"

)

// AutoRefreshMiddleware provides automatic token refresh functionality
// This middleware can be used to add automatic token refresh to protected routes
func AutoRefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user has a session token and if it needs refresh
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			// No session token, continue normally
			next.ServeHTTP(w, r)
			return
		}

		// Check if token expires within next 5 minutes
		if needsTokenRefresh(cookie.Value) {
			fmt.Printf("🔄 AUTO-REFRESH: Token expiring soon, attempting refresh...\n")
			
			// Automatically refresh the token
			refreshSuccess := attemptAutoTokenRefresh(w, r)
			if refreshSuccess {
				fmt.Printf("✅ AUTO-REFRESH: Token refreshed successfully!\n")
			} else {
				fmt.Printf("❌ AUTO-REFRESH: Refresh failed\n")
			}
		}

		// Continue with the request
		next.ServeHTTP(w, r)
	})
}

// needsTokenRefresh checks if a JWT token expires within the next 5 minutes
func needsTokenRefresh(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return true // Invalid token format
	}

	// Decode the payload (middle part)
	payload, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return true // Can't decode payload
	}

	// Parse expiration
	var claims struct {
		Exp int64 `json:"exp"`
	}
	
	if err := json.Unmarshal(payload, &claims); err != nil {
		return true // Can't parse claims
	}

	// Check if token expires within next 5 minutes
	now := time.Now().Unix()
	refreshThreshold := now + (5 * 60) // 5 minutes

	return claims.Exp < refreshThreshold
}

// attemptAutoTokenRefresh attempts to refresh the session token automatically
func attemptAutoTokenRefresh(w http.ResponseWriter, r *http.Request) bool {
	// Get refresh token cookie
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Printf("❌ AUTO-REFRESH: No refresh token cookie found\n")
		return false
	}

	// Call auth service to get new tokens
	cfg := config.LoadConfig()
	authService := auth.NewService(cfg)
	refreshResp, err := authService.RefreshToken(refreshCookie.Value)
	
	if err != nil || !refreshResp.Success || refreshResp.Token == "" {
		fmt.Printf("❌ AUTO-REFRESH: Auth service failed: %v\n", err)
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

	fmt.Printf("✅ AUTO-REFRESH: New session token set: %d chars\n", len(refreshResp.Token))
	return true
}

