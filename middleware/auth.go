package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
		fmt.Printf("🔐 MIDDLEWARE: No session cookie found: %v\n", err)
		return templates.UserInfo{LoggedIn: false}
	}

	if cookie.Value == "" {
		fmt.Printf("🔐 MIDDLEWARE: Empty session token\n")
		return templates.UserInfo{LoggedIn: false}
	}

	fmt.Printf("🔐 MIDDLEWARE: Validating JWT token, length: %d\n", len(cookie.Value))

	// Parse JWT to get real user data
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 3 {
		fmt.Printf("🔐 MIDDLEWARE: Invalid JWT format, parts: %d\n", len(parts))
		return templates.UserInfo{LoggedIn: false}
	}

	// Decode payload (the middle part)
	payload, err := jwtBase64URLDecode(parts[1])
	if err != nil {
		fmt.Printf("🔐 MIDDLEWARE: Failed to decode JWT payload: %v\n", err)
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
		fmt.Printf("🔐 MIDDLEWARE: Failed to parse JWT claims: %v\n", err)
		fmt.Printf("🔐 MIDDLEWARE: Payload: %s\n", string(payload))
		return templates.UserInfo{LoggedIn: false}
	}

	fmt.Printf("🔐 MIDDLEWARE: JWT claims - Name: %s, Email: %s, Issuer: %s\n",
		claims.Name, claims.Email, claims.Iss)

	// Check if token is still valid (not expired)
	if claims.Exp < time.Now().Unix() {
		fmt.Printf("🔐 MIDDLEWARE: JWT expired. Exp: %d, Now: %d\n", claims.Exp, time.Now().Unix())
		return templates.UserInfo{LoggedIn: false}
	}

	// Check issuer to make sure it's from our auth service
	if claims.Iss != "auth-ms" {
		fmt.Printf("🔐 MIDDLEWARE: Invalid JWT issuer: %s (expected: auth-ms)\n", claims.Iss)
		return templates.UserInfo{LoggedIn: false}
	}

	fmt.Printf("🔐 MIDDLEWARE: JWT validation successful for %s (%s)\n", claims.Name, claims.Email)

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
	switch len(data) % 4 {
	case 2:
		data += "=="
	case 3:
		data += "="
	case 1:
		return nil, fmt.Errorf("invalid base64url length")
	}

	return base64.URLEncoding.DecodeString(data)
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