package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
)

// UserContextKey is the key used to store user info in request context
type UserContextKey string

const userContextKey UserContextKey = "user"

// Route categories for middleware configuration
var (
	// Protected routes that require authentication
	ProtectedRoutes = []string{
		"/profile", // User profile page
		"/admin",   // Admin dashboard
	}

	// Admin routes that require authentication
	AdminRoutes = []string{
		"/api/admin", // All admin API endpoints
	}

	// Public routes that don't require authentication
	PublicRoutes = []string{
		"/",       // Homepage
		"/health", // Health check endpoint
		"/login",  // Login page
		"/test",   // Test page for development
	}

	// OAuth routes that are part of authentication flow
	OAuthRoutes = []string{
		"/auth/google",
		"/auth/github",
		"/auth/discord",
		"/auth/microsoft",
		"/auth/callback",
	}

	// Auth API routes that handle session management
	AuthAPIRoutes = []string{
		"/api/auth/validate",
		"/api/auth/user",
		"/api/auth/logout",
		"/api/auth/set-session",
		"/api/auth/exchange-code",
		"/api/auth/test-session-create",
		"/api/auth/refresh",
	}
)

// RequiresAuthentication checks if a route requires authentication
func RequiresAuthentication(path string) bool {
	// Check exact matches first
	for _, route := range ProtectedRoutes {
		if path == route {
			return true
		}
	}

	// Check admin routes (prefix matching)
	for _, route := range AdminRoutes {
		if len(path) >= len(route) && path[:len(route)] == route {
			return true // All admin API routes are protected
		}
	}

	return false
}

// IsPublicRoute checks if a route is explicitly public
func IsPublicRoute(path string) bool {
	for _, route := range PublicRoutes {
		if path == route {
			return true
		}
	}

	// Check OAuth routes (prefix matching for /auth/*)
	for _, route := range OAuthRoutes {
		if len(path) >= len(route) && path[:len(route)] == route {
			return true
		}
	}

	return false
}

// GetRouteCategory returns the category of a route for debugging
func GetRouteCategory(path string) string {
	if RequiresAuthentication(path) {
		return "PROTECTED"
	}
	if IsPublicRoute(path) {
		return "PUBLIC"
	}
	if strings.HasPrefix(path, "/api/auth/") {
		return "AUTH_API"
	}
	return "UNKNOWN"
}

// AuthMiddleware validates server sessions for protected routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		category := GetRouteCategory(path)

		fmt.Printf("🔐 MIDDLEWARE: Processing route %s [Category: %s]\n", path, category)

		// Always validate session for all routes (to show logged-in status)
		userInfo := validateSession(r)
		ctx := context.WithValue(r.Context(), userContextKey, userInfo)

		// Check if this route requires authentication
		if RequiresAuthentication(path) {
			// If route requires auth but user is not logged in, redirect
			if !userInfo.LoggedIn {
				if r.URL.Path[:5] == "/api/" {
					// For API routes, return JSON error
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					if err := json.NewEncoder(w).Encode(map[string]interface{}{
						"error": "Authentication required",
					}); err != nil {
						fmt.Printf("🔐 MIDDLEWARE: Failed to encode error response: %v\n", err)
					}
					return
				}

				// For web routes, redirect to login
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext gets user info from request context
func GetUserFromContext(r *http.Request) layouts.UserInfo {
	userInfo, ok := r.Context().Value(userContextKey).(layouts.UserInfo)
	if !ok {
		return layouts.UserInfo{LoggedIn: false}
	}
	return userInfo
}
