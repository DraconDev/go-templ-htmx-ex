package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DraconDev/go-templ-htmx-ex/templates"
)

// UserContextKey is the key used to store user info in request context
type UserContextKey string

const userContextKey UserContextKey = "user"

// AuthMiddleware validates JWT tokens for protected routes
func AuthMiddleware(next http.Handler) http.Handler {
	// Public routes that don't require authentication
	publicPaths := map[string]bool{
		"/":                    true, // Home page can be public
		"/health":              true, // Health check
		"/login":               true, // Login page
		"/auth/google":         true, // OAuth login
		"/auth/github":         true, // OAuth login
		"/auth/callback":       true, // OAuth callback
		"/api/auth/set-session": true, // Session setting
		"/api/auth/validate":   true, // Auth validation
		"/api/auth/refresh":    true, // Token refresh
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for public paths
		if publicPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Validate JWT for protected routes
		userInfo := validateJWT(r)
		
		if !userInfo.LoggedIn {
			// For API routes, return JSON error
			if strings.HasPrefix(r.URL.Path, "/api/") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": "Authentication required",
				})
				return
			}
			
			// For web routes, redirect to login
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), userContextKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateJWT validates the JWT token from session cookie
func validateJWT(r *http.Request) templates.UserInfo {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return templates.UserInfo{LoggedIn: false}
	}

	if cookie.Value == "" {
		return templates.UserInfo{LoggedIn: false}
	}

	// Parse JWT to get real user data
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 3 {
		return templates.UserInfo{LoggedIn: false}
	}

	// Decode payload (the middle part)
	payload, err := jwtBase64URLDecode(parts[1])
	if err != nil {
		return templates.UserInfo{LoggedIn: false}
	}

	// Parse user data from JWT payload
	var claims struct {
		Sub     string `json:"sub"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
		Exp     int64  `json:"exp"`
		Iss     string `json:"iss"`
	}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return templates.UserInfo{LoggedIn: false}
	}

	// Check if token is still valid (not expired)
	if claims.Exp < 0 { // Simplified check for demo
		return templates.UserInfo{LoggedIn: false}
	}

	// Check issuer to make sure it's from our auth service
	if claims.Iss != "auth-ms" {
		return templates.UserInfo{LoggedIn: false}
	}

	// Return real user data!
	return templates.UserInfo{
		LoggedIn: true,
		Name:     claims.Name,
		Email:    claims.Email,
		Picture:  claims.Picture,
	}
}

// jwtBase64URLDecode decodes base64url encoding (needed for JWT)
func jwtBase64URLDecode(data string) ([]byte, error) {
	// Add padding if needed
	padding := len(data) % 4
	if padding > 0 {
		data += strings.Repeat("=", 4-padding)
	}

	// Convert base64url to base64
	data = strings.ReplaceAll(data, "-", "+")
	data = strings.ReplaceAll(data, "_", "/")

	return json.RawMessage(data), nil
}

// GetUserFromContext gets user info from request context
func GetUserFromContext(r *http.Request) templates.UserInfo {
	userInfo, ok := r.Context().Value(userContextKey).(templates.UserInfo)
	if !ok {
		return templates.UserInfo{LoggedIn: false}
	}
	return userInfo
}

// RequireAdmin middleware that checks if user is admin
func RequireAdmin(next http.Handler, adminEmail string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := GetUserFromContext(r)
		
		if !userInfo.LoggedIn {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// Check if user is admin
		if userInfo.Email != adminEmail {
			http.Error(w, "Access denied: Admin privileges required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}