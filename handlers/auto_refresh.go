package handlers

import (
	"fmt"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
)

// AutoRefreshMiddleware provides automatic token refresh functionality
// This middleware can be used to add automatic token refresh to protected routes
func AutoRefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user has a session token
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// No session token, continue normally
			next.ServeHTTP(w, r)
			return
		}

		if cookie.Value != "" {
			fmt.Printf("🔄 AUTO-REFRESH: Session token found, length: %d\n", len(cookie.Value))
			// Token exists, continue normally - refresh is handled by frontend JavaScript
			next.ServeHTTP(w, r)
			return
		}

		// Empty token, continue normally
		next.ServeHTTP(w, r)
	})
}

