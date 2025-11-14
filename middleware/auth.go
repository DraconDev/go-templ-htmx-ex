package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
)

// UserContextKey is the key used to store user info in request context
type UserContextKey string

const userContextKey UserContextKey = "user"

// SessionCache stores validation results with 15-second TTL
type SessionCache struct {
	sync.RWMutex
	entries map[string]*cacheEntry
}

type cacheEntry struct {
	userInfo layouts.UserInfo
	expiresAt time.Time
}

// NewSessionCache creates a new session cache
func NewSessionCache() *SessionCache {
	return &SessionCache{
		entries: make(map[string]*cacheEntry),
	}
}

// Get retrieves cached user info if not expired
func (c *SessionCache) Get(sessionID string) (layouts.UserInfo, bool) {
	c.RLock()
	defer c.RUnlock()
	
	entry, exists := c.entries[sessionID]
	if !exists {
		return layouts.UserInfo{LoggedIn: false}, false
	}
	
	if time.Now().After(entry.expiresAt) {
		// Expired entry, clean up
		delete(c.entries, sessionID)
		return layouts.UserInfo{LoggedIn: false}, false
	}
	
	return entry.userInfo, true
}

// Set caches user info with 15-second TTL
func (c *SessionCache) Set(sessionID string, userInfo layouts.UserInfo) {
	c.Lock()
	defer c.Unlock()
	
	c.entries[sessionID] = &cacheEntry{
		userInfo:  userInfo,
		expiresAt: time.Now().Add(15 * time.Second),
	}
}

// Global session cache instance
var sessionCache = NewSessionCache()

// AuthMiddleware validates server sessions for protected routes
func AuthMiddleware(next http.Handler) http.Handler {
	// Protected routes that require authentication
	protectedPaths := map[string]bool{
		"/profile":   true, // User profile page
		"/admin":     true, // Admin dashboard
		"/api/admin": true, // Admin API routes
	}

	// API routes that require authentication
	apiProtectedPaths := map[string]bool{
		"/api/admin": true, // Admin API routes
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always validate session and add to context (for UI purposes)
		userInfo := validateSession(r)
		ctx := context.WithValue(r.Context(), userContextKey, userInfo)

		// Check if this route requires authentication
		var requiresAuth bool
		if r.URL.Path[:5] == "/api/" {
			requiresAuth = apiProtectedPaths[r.URL.Path]
		} else {
			requiresAuth = protectedPaths[r.URL.Path]
		}

		// If route requires auth but user is not logged in, redirect
		if requiresAuth && !userInfo.LoggedIn {
			if r.URL.Path[:5] == "/api/" {
				// For API routes, return JSON error
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": "Authentication required",
				})
				return
			}

			// For web routes, redirect to home
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateSession validates server session from session_id cookie with 15-second caching
func validateSession(r *http.Request) layouts.UserInfo {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Printf("🔐 MIDDLEWARE: No session cookie found: %v\n", err)
		return layouts.UserInfo{LoggedIn: false}
	}

	if cookie.Value == "" {
		fmt.Printf("🔐 MIDDLEWARE: Empty session ID\n")
		return layouts.UserInfo{LoggedIn: false}
	}

	fmt.Printf("🔐 MIDDLEWARE: Validating session, ID length: %d\n", len(cookie.Value))

	// Check cache first (15-second TTL)
	if cached, found := sessionCache.Get(cookie.Value); found {
		fmt.Printf("🔐 MIDDLEWARE: Cache hit for session %s\n", cookie.Value[:8]+"...")
		return cached
	}

	fmt.Printf("🔐 MIDDLEWARE: Cache miss - calling auth service for session %s\n", cookie.Value[:8]+"...")

	// Cache miss - call auth service to validate session
	userInfo, err := validateSessionWithAuthService(cookie.Value)
	if err != nil {
		fmt.Printf("🔐 MIDDLEWARE: Auth service validation failed: %v\n", err)
		return layouts.UserInfo{LoggedIn: false}
	}

	// Cache result for 15 seconds
	sessionCache.Set(cookie.Value, userInfo)
	
	return userInfo
}

// validateSessionWithAuthService validates session by calling auth microservice
func validateSessionWithAuthService(sessionID string) (layouts.UserInfo, error) {
	// This would call the auth service to validate the session
	// For now, we'll implement this in the auth service package
	// Since the auth service has Redis and manages sessions
	
	// TODO: Implement actual auth service call
	// This is a placeholder that would call auth service /auth/validate-session
	
	fmt.Printf("🔐 MIDDLEWARE: Would call auth service to validate session %s\n", sessionID[:8]+"...")
	
	// Return invalid for now - this needs to be implemented
	return layouts.UserInfo{LoggedIn: false}, fmt.Errorf("session validation not yet implemented")
}

// GetUserFromContext gets user info from request context
func GetUserFromContext(r *http.Request) layouts.UserInfo {
	userInfo, ok := r.Context().Value(userContextKey).(layouts.UserInfo)
	if !ok {
		return layouts.UserInfo{LoggedIn: false}
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
