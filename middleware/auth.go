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
	fmt.Printf("🔐 MIDDLEWARE: Calling auth service to validate session %s\n", sessionID[:8]+"...")
	
	// Make HTTP request to auth microservice
	// Since we don't have the auth service instance here, we'll call the API directly
	client := &http.Client{Timeout: 5 * time.Second}
	
	// Prepare request to validate session
	req, err := http.NewRequest("POST", "http://localhost:8081/api/auth/validate-session", nil)
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}, fmt.Errorf("failed to create request: %v", err)
	}
	
	// Add session ID in request body
	type SessionRequest struct {
		SessionID string `json:"session_id"`
	}
	
	sessionReq := SessionRequest{SessionID: sessionID}
	reqData, err := json.Marshal(sessionReq)
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}, fmt.Errorf("failed to marshal request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Body = http.NoBody
	
	// Note: In a real implementation, this would forward the session cookie
	// For now, we'll use a placeholder implementation
	
	fmt.Printf("🔐 MIDDLEWARE: Session validation API call would go here\n")
	
	// Placeholder: Return invalid session for now
	// This would be implemented to call the auth service and return user info
	return layouts.UserInfo{LoggedIn: false}, fmt.Errorf("session validation API not implemented")
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
