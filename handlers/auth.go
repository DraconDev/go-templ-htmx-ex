package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/templates"
)

// =============================================================================
// AUTHENTICATION HANDLER
// =============================================================================
// This handler manages the complete OAuth + JWT authentication flow for the app:
// 1. OAuth redirects to external providers (Google, GitHub)
// 2. Callback processing to extract JWT tokens from URL fragments
// 3. Session management with HTTP-only cookies
// 4. Token validation and refresh logic
// =============================================================================

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	AuthService *auth.Service // Communication with auth microservice
	Config      *config.Config // App configuration
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *auth.Service, config *config.Config) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Config:      config,
	}
}

// =============================================================================
// TEST ROUTES
// =============================================================================
// These routes help test the authentication flow during development

// TestTokenRefreshHandler serves a test page with token refresh button
func (h *AuthHandler) TestTokenRefreshHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	testHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Auth Test - Token Refresh</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-blue-300 min-h-screen">
    <div class="container mx-auto py-8 px-4">
        <h1 class="text-3xl font-bold text-center mb-8">Authentication Test Page</h1>
        
        <div class="max-w-2xl mx-auto space-y-6">
            <!-- Test Token Refresh -->
            <div class="bg-white rounded-lg shadow p-6">
                <h2 class="text-xl font-semibold mb-4">Test Token Refresh</h2>
                <p class="text-gray-600 mb-4">This button will test the token refresh flow.</p>
                <button
                    hx-post="/api/auth/refresh"
                    hx-target="#refresh-result"
                    hx-swap="innerHTML"
                    class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
                >
                    Test Token Refresh
                </button>
                <div id="refresh-result" class="mt-4 p-3 bg-gray-100 rounded"></div>
            </div>
            
            <!-- Check Current User -->
            <div class="bg-white rounded-lg shadow p-6">
                <h2 class="text-xl font-semibold mb-4">Check Current User</h2>
                <p class="text-gray-600 mb-4">Check if user is currently logged in.</p>
                <button
                    hx-get="/api/auth/user"
                    hx-target="#user-result"
                    hx-swap="innerHTML"
                    class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                >
                    Check User Status
                </button>
                <div id="user-result" class="mt-4 p-3 bg-gray-100 rounded"></div>
            </div>
            
            <!-- OAuth Login Buttons -->
            <div class="bg-white rounded-lg shadow p-6">
                <h2 class="text-xl font-semibold mb-4">OAuth Login</h2>
                <div class="space-x-4">
                    <a href="/auth/google" class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">Login with Google</a>
                    <a href="/auth/github" class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded">Login with GitHub</a>
                </div>
            </div>
            
            <!-- Callback Test -->
            <div class="bg-white rounded-lg shadow p-6">
                <h2 class="text-xl font-semibold mb-4">Test Callback</h2>
                <p class="text-gray-600 mb-4">Test the callback page that processes JWT tokens.</p>
                <a href="/auth/callback#access_token=test-jwt-token&token_type=Bearer"
                   class="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded">
                    Test Callback with Fake Token
                </a>
            </div>
        </div>
    </div>
</body>
</html>
`
	w.Write([]byte(testHTML))
}
// OAUTH LOGIN FLOWS
// =============================================================================

// GoogleLoginHandler handles Google OAuth login
// Flow: User clicks "Login with Google" -> Redirect to our auth service ->
//       Auth service handles Google OAuth -> Returns to our callback with JWT
func (h *AuthHandler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 GOOGLE LOGIN: Starting Google OAuth flow\n")
	fmt.Printf("🔐 GOOGLE LOGIN: AuthServiceURL = %s\n", h.Config.AuthServiceURL)
	fmt.Printf("🔐 GOOGLE LOGIN: RedirectURL = %s\n", h.Config.RedirectURL)

	// STEP 1: Redirect to our auth microservice with redirect_uri parameter
	// The auth service will handle the actual Google OAuth flow
	authURL := fmt.Sprintf("%s/auth/google?redirect_uri=%s/auth/callback",
		h.Config.AuthServiceURL, h.Config.RedirectURL)

	fmt.Printf("🔐 GOOGLE LOGIN: Redirecting to: %s\n", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *AuthHandler) GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 GITHUB LOGIN: Starting GitHub OAuth flow\n")
	fmt.Printf("🔐 GITHUB LOGIN: AuthServiceURL = %s\n", h.Config.AuthServiceURL)
	fmt.Printf("🔐 GITHUB LOGIN: RedirectURL = %s\n", h.Config.RedirectURL)

	// OAuth endpoints are public - just redirect
	authURL := fmt.Sprintf("%s/auth/github?redirect_uri=%s/auth/callback",
		h.Config.AuthServiceURL, h.Config.RedirectURL)

	fmt.Printf("🔐 GITHUB LOGIN: Redirecting to: %s\n", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// AuthCallbackHandler handles the OAuth callback
// Flow: Google redirects here with JWT in URL fragment (#access_token=...)
//       Client-side JS extracts token and calls /api/auth/set-session
func (h *AuthHandler) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 CALLBACK: === OAuth callback STARTED ===\n")
	fmt.Printf("🔐 CALLBACK: URL = %s\n", r.URL.String())
	fmt.Printf("🔐 CALLBACK: Query params = %v\n", r.URL.Query())
	fmt.Printf("🔐 CALLBACK: Fragment = %s\n", r.URL.Fragment)

	fmt.Printf("🔐 CALLBACK: Setting content type and rendering template...\n")
	w.Header().Set("Content-Type", "text/html")
	
	// STEP 2: Render callback page with JavaScript to extract JWT from URL fragment
	// The fragment (#access_token=...) is not sent to server, so JS must handle it
	component := templates.Layout("Authenticating", templates.NavigationLoggedOut(), templates.AuthCallbackContent())

	fmt.Printf("🔐 CALLBACK: About to render component...\n")
	component.Render(r.Context(), w)
	fmt.Printf("🔐 CALLBACK: Component rendered successfully\n")
	fmt.Printf("🔐 CALLBACK: === OAuth callback COMPLETED ===\n")
}

// SetSessionHandler sets the user session from client-side JavaScript
func (h *AuthHandler) SetSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 SESSION: === Set session STARTED ===\n")
	fmt.Printf("🔐 SESSION: Content-Type: %s\n", r.Header.Get("Content-Type"))

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("🔐 SESSION: Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	if req.Token == "" {
		fmt.Printf("🔐 SESSION: Missing token in request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing token",
		})
		return
	}

	fmt.Printf("🔐 SESSION: Token received, length: %d\n", len(req.Token))

	// Set session cookie with the JWT token
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    req.Token,
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	http.SetCookie(w, cookie)
	fmt.Printf("🔐 SESSION: Cookie set with name: %s, value length: %d\n", cookie.Name, len(cookie.Value))

	fmt.Printf("🔐 SESSION: Sending success response\n")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Session set successfully",
	})
	fmt.Printf("🔐 SESSION: === Set session COMPLETED ===\n")
}

// GetUserHandler returns current user information
func (h *AuthHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 GETUSER: === GetUser STARTED ===\n")
	w.Header().Set("Content-Type", "application/json")

	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Printf("🔐 GETUSER: No session cookie found: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"logged_in": false,
		})
		return
	}

	fmt.Printf("🔐 GETUSER: Session cookie found, value length: %d\n", len(cookie.Value))

	// Get user info from auth microservice
	fmt.Printf("🔐 GETUSER: Calling auth service to validate user...\n")
	userResp, err := h.AuthService.ValidateUser(cookie.Value)
	if err != nil {
		fmt.Printf("🔐 GETUSER: Auth service failed: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"logged_in": false,
		})
		return
	}

	fmt.Printf("🔐 GETUSER: Auth service response - Success: %v, Name: %s\n", userResp.Success, userResp.Name)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logged_in": userResp.Success,
		"user_id":   userResp.UserID,
		"email":     userResp.Email,
		"name":      userResp.Name,
		"picture":   userResp.Picture,
	})

	fmt.Printf("🔐 GETUSER: === GetUser COMPLETED ===\n")
}

// ValidateSessionHandler validates the current session
func (h *AuthHandler) ValidateSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
			"error": "No session token",
		})
		return
	}

	// Validate token with auth microservice
	userResp, err := h.AuthService.ValidateToken(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   userResp.Success,
		"user_id": userResp.UserID,
		"email":   userResp.Email,
		"name":    userResp.Name,
		"picture": userResp.Picture,
		"status":  "validated",
	})
}

// LogoutHandler handles user logout
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

// GetUserInfo returns current user information for server-side rendering
func (h *AuthHandler) GetUserInfo(r *http.Request) templates.UserInfo {
	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return templates.UserInfo{LoggedIn: false}
	}

	// Get user info from auth microservice
	userResp, err := h.AuthService.ValidateUser(cookie.Value)
	if err != nil {
		return templates.UserInfo{LoggedIn: false}
	}

	return templates.UserInfo{
		LoggedIn: userResp.Success,
		Name:     userResp.Name,
		Email:    userResp.Email,
		Picture:  userResp.Picture,
	}
}

// RefreshTokenHandler handles token refresh requests
// Flow: Frontend calls when JWT expires ->
//       Server reads refresh_token cookie ->
//       Calls auth service for new JWT ->
//       **Sets new session_token cookie automatically**
func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// STEP 1: Get refresh token from HTTP-only cookie (automatically sent by browser)
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "No refresh token found",
		})
		return
	}

	// STEP 2: Call auth service to refresh token using the refresh token
	userResp, err := h.AuthService.RefreshToken(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// STEP 3: **CRITICAL** - Set the new JWT cookie for the user
	// This replaces the expired session_token with a fresh one
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    userResp.Token, // NEW JWT from auth service
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Token refreshed successfully",
	})
}
