package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
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

// parseAndValidateJWT parses and validates JWT claims
func parseAndValidateJWT(token string) (layouts.UserInfo, bool) {
	// Parse JWT to get real user data
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		fmt.Printf("🔐 MIDDLEWARE: Invalid JWT format, parts: %d\n", len(parts))
		return layouts.UserInfo{LoggedIn: false}, false
	}

	// Decode payload (the middle part)
	payload, err := jwtBase64URLDecode(parts[1])
	if err != nil {
		fmt.Printf("🔐 MIDDLEWARE: Failed to decode JWT payload: %v\n", err)
		return layouts.UserInfo{LoggedIn: false}, false
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
		return layouts.UserInfo{LoggedIn: false}, false
	}

	fmt.Printf("🔐 MIDDLEWARE: JWT claims - Name: %s, Email: %s, Issuer: %s\n",
		claims.Name, claims.Email, claims.Iss)

	// Check if token is still valid (not expired)
	if claims.Exp < time.Now().Unix() {
		fmt.Printf("🔐 MIDDLEWARE: JWT expired. Exp: %d, Now: %d\n", claims.Exp, time.Now().Unix())
		return layouts.UserInfo{LoggedIn: false}, false
	}

	// Check issuer to make sure it's from our auth service
	if claims.Iss != "auth-ms" {
		fmt.Printf("🔐 MIDDLEWARE: Invalid JWT issuer: %s (expected: auth-ms)\n", claims.Iss)
		return layouts.UserInfo{LoggedIn: false}, false
	}

	fmt.Printf("🔐 MIDDLEWARE: JWT validation successful for %s (%s)\n", claims.Name, claims.Email)

	// Return real user data!
	return layouts.UserInfo{
		LoggedIn: true,
		Name:     claims.Name,
		Email:    claims.Email,
		Picture:  claims.Picture,
	}, true
}

// CheckIfTokenNeedsRefresh checks if the session token is expired and needs refresh
func CheckIfTokenNeedsRefresh(token string) bool {
	if token == "" {
		return true
	}
	
	// Parse JWT to get expiration
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return true
	}
	
	// Decode payload (the middle part)
	payload, err := jwtBase64URLDecode(parts[1])
	if err != nil {
		return true
	}
	
	// Parse expiration from JWT payload
	var claims struct {
		Exp int64 `json:"exp"`
	}
	
	if err := json.Unmarshal(payload, &claims); err != nil {
		return true
	}
	
	// Check if token expires within next 5 minutes (to refresh proactively)
	now := time.Now().Unix()
	expiryThreshold := now + (5 * 60) // 5 minutes from now
	
	return claims.Exp < expiryThreshold
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
