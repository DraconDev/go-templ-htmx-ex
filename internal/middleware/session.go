package middleware

import (
	"fmt"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
)

// Global session cache instance - will be initialized in service.go
var sessionCache *SessionCache

// InitializeSessionCache initializes the global session cache
func InitializeSessionCache() {
	if sessionCache == nil {
		sessionCache = NewSessionCache()
	}
}

// validateSession validates server session from session_id cookie with 15-second caching
func validateSession(r *http.Request) layouts.UserInfo {
	// Ensure cache is initialized
	if sessionCache == nil {
		InitializeSessionCache()
	}

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
// validateSessionWithAuthService validates session by calling auth microservice
func validateSessionWithAuthService(sessionID string) (layouts.UserInfo, error) {
	fmt.Printf("🔐 MIDDLEWARE: Calling auth service to validate session %s\n", sessionID[:8]+"...")

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Prepare request to auth service
	reqBody := map[string]string{"session_id": sessionID}
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
}
