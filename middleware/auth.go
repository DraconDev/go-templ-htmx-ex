package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
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
	userInfo  layouts.UserInfo
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		category := GetRouteCategory(path)

		fmt.Printf("🔐 MIDDLEWARE: Processing route %s [Category: %s]\n", path, category)

		// Always validate session for all routes (to show logged-in status)
		userInfo := validateSession(r)
		ctx := context.WithValue(r.Context(), userContextKey, userInfo)

		// Check if this route requires authentication
		if requiresAuthentication(path) {
			// If route requires auth but user is not logged in, redirect
			if !userInfo.LoggedIn {
				if r.URL.Path[:5] == "/api/" {
					// For API routes, return JSON error
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error": "Authentication required",
					})
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

// validateSession validates server session from session_id cookie with 15-second caching
func validateSession(r *http.Request) layouts.UserInfo {
	// Get session_id cookie for server sessions
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
		// Return unauthenticated instead of crashing
		return layouts.UserInfo{LoggedIn: false}
	}

	// Cache result for 15 seconds
	sessionCache.Set(cookie.Value, userInfo)

	return userInfo
}

// validateSessionWithAuthService validates session by calling auth microservice
func validateSessionWithAuthService(sessionID string) (layouts.UserInfo, error) {
	fmt.Printf("🔐 MIDDLEWARE: Calling auth service to validate session %s\n", sessionID[:8]+"...")

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Prepare request to auth service
	reqBody := map[string]string{"session_token": sessionID}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("🔐 MIDDLEWARE: Failed to marshal request: %v\n", err)
		return layouts.UserInfo{LoggedIn: false}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/session/refresh", config.Current.AuthServiceURL), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("🔐 MIDDLEWARE: Failed to create request: %v\n", err)
		return layouts.UserInfo{LoggedIn: false}, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Add auth secret if configured
	if config.Current.AuthSecret != "" {
		req.Header.Set("X-Auth-Secret", config.Current.AuthSecret)
	}

	fmt.Printf("🔐 MIDDLEWARE: Sending validation request to auth service\n")

	// Send request to auth service
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("🔐 MIDDLEWARE: Failed to call auth service: %v\n", err)
		// Don't fail the request if auth service is unavailable
		return layouts.UserInfo{LoggedIn: false}, nil
	}
	defer resp.Body.Close()

	fmt.Printf("🔐 MIDDLEWARE: Auth service response status: %s\n", resp.Status)

	// Parse response
	var respData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		fmt.Printf("🔐 MIDDLEWARE: Failed to parse response: %v\n", err)
		return layouts.UserInfo{LoggedIn: false}, nil
	}

	fmt.Printf("🔐 MIDDLEWARE: Auth service response: %v\n", respData)

	// Check if session is valid by looking for user_context
	if userContext, ok := respData["user_context"].(map[string]interface{}); ok && userContext != nil {
		// Session is valid - extract user info from user_context
		userInfo := layouts.UserInfo{
			LoggedIn: true,
		}

		if name, ok := userContext["name"].(string); ok && name != "" {
			userInfo.Name = name
		}
		if email, ok := userContext["email"].(string); ok && email != "" {
			userInfo.Email = email
		}
		if picture, ok := userContext["picture"].(string); ok && picture != "" {
			userInfo.Picture = picture
		}
		if userID, ok := userContext["user_id"].(string); ok && userID != "" {
			fmt.Printf("🔐 MIDDLEWARE: Session valid for user: %s (%s)\n", userInfo.Name, userInfo.Email)
		}

		return userInfo, nil
	}

	fmt.Printf("🔐 MIDDLEWARE: Session validation failed\n")
	return layouts.UserInfo{LoggedIn: false}, nil
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
