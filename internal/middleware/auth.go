package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/internal/routes"
	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
)

// UserContextKey is the key used to store user info in request context
type UserContextKey string

const userContextKey UserContextKey = "user"

// AuthMiddleware validates server sessions for protected routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		category := routes.GetRouteCategory(path)

		fmt.Printf("🔐 MIDDLEWARE: Processing route %s [Category: %s]\n", path, category)

		// Always validate session for all routes (to show logged-in status)
		userInfo := validateSession(r)
		ctx := context.WithValue(r.Context(), userContextKey, userInfo)

		// Check if this route requires authentication
		if routes.RequiresAuthentication(path) {
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
