package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
	"github.com/DraconDev/go-templ-htmx-ex/templates/pages"
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
	AuthService *auth.Service  // Communication with auth microservice
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
    <script>
        // =============================================================================
        // TOKEN REFRESH TESTING FUNCTIONS
        // =============================================================================
        
        // Simple token refresh test - browser automatically handles HTTP-only cookies
        function testTokenRefresh() {
            console.log('🔄 TOKEN REFRESH TEST: Testing refresh token flow...');
            const resultDiv = document.getElementById('refresh-result');
            
            resultDiv.innerHTML = '<p class="text-blue-600">🔄 Refreshing token...</p>';
            
            // Browser automatically sends refresh_token cookie (HTTP-only, inaccessible to JS)
            fetch('/api/auth/refresh', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => {
                console.log('🔄 TOKEN REFRESH TEST: Response status:', response.status);
                return response.json().then(data => {
                    console.log('🔄 TOKEN REFRESH TEST: Response data:', data);
                    return { ok: response.ok, data: data };
                });
            })
            .then(result => {
                if (result.ok && result.data.success) {
                    console.log('✅ TOKEN REFRESH TEST: SUCCESS');
                    resultDiv.innerHTML = '<p class="text-green-600">✅ Token refreshed successfully!</p>';
                } else {
                    console.log('❌ TOKEN REFRESH TEST: FAILED');
                    resultDiv.innerHTML = '<p class="text-red-600">❌ Token refresh failed: ' + (result.data.error || 'Unknown error') + '</p>';
                }
            })
            .catch(error => {
                console.log('❌ TOKEN REFRESH TEST: Error:', error);
                resultDiv.innerHTML = '<p class="text-red-600">❌ Network error: ' + error.message + '</p>';
            });
        }
        
        // Helper function to check current user status
        function checkUserStatus() {
            console.log('👤 USER STATUS CHECK: === STARTED ===');
            const resultDiv = document.getElementById('user-result');
            
            fetch('/api/auth/user')
            .then(response => response.json())
            .then(data => {
                console.log('👤 USER STATUS CHECK: Response:', data);
                
                if (data.logged_in) {
                    resultDiv.innerHTML = '<p class="text-green-600">✅ Logged in as: ' + data.name + '</p>' +
                        '<p class="text-sm text-gray-600">Email: ' + data.email + '</p>';
                } else {
                    resultDiv.innerHTML = '<p class="text-red-600">❌ Not logged in</p>';
                }
            })
            .catch(error => {
                console.log('👤 USER STATUS CHECK: Error:', error);
                resultDiv.innerHTML = '<p class="text-red-600">❌ Error: ' + error.message + '</p>';
            });
        }
        
        // Log when page loads
        document.addEventListener('DOMContentLoaded', function() {
            console.log('🧪 AUTH TEST PAGE: Loaded - Check browser console for detailed testing logs');
        });
    </script>
</head>
<body class="bg-blue-300 min-h-screen">
    <div class="container mx-auto py-8 px-4">
        <h1 class="text-3xl font-bold text-center mb-8">Authentication Test Page</h1>
        
        <div class="max-w-2xl mx-auto space-y-6">
            <!-- Test Token Refresh -->
            <div class="bg-white rounded-lg shadow p-6">
                <h2 class="text-xl font-semibold mb-4">Test Token Refresh</h2>
                <p class="text-gray-600 mb-4">This button will test the complete token refresh flow with detailed console logging.</p>
                <button
                    onclick="testTokenRefresh()"
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
                    onclick="checkUserStatus()"
                    class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                >
                    Check User Status
                </button>
                <div id="user-result" class="mt-4 p-3 bg-gray-100 rounded"></div>
            </div>
            
            <!-- OAuth Login Buttons -->
            <div class="bg-white rounded-lg shadow p-6">
                <h2 class="text-xl font-semibold mb-4">OAuth Login</h2>
                <p class="text-gray-600 mb-4">Use these to test the full OAuth flow.</p>
                <div class="space-x-4">
                    <a href="/auth/google" class="bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white font-bold py-3 px-6 rounded-lg transition-all duration-300 shadow-lg hover:shadow-red-500/25 border border-red-400/30">Login with Google</a>
                    <a href="/auth/github" class="bg-gradient-to-r from-gray-600 to-gray-700 hover:from-gray-700 hover:to-gray-800 text-white font-bold py-3 px-6 rounded-lg transition-all duration-300 shadow-lg hover:shadow-gray-500/25 border border-gray-400/30">Login with GitHub</a>
                    <a href="/auth/discord" class="bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-600 hover:to-indigo-700 text-white font-bold py-3 px-6 rounded-lg transition-all duration-300 shadow-lg hover:shadow-indigo-500/25 border border-indigo-400/30">Login with Discord</a>
                    <a href="/auth/microsoft" class="bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white font-bold py-3 px-6 rounded-lg transition-all duration-300 shadow-lg hover:shadow-blue-500/25 border border-blue-400/30">Login with Microsoft</a>
                </div>
            </div>
            
            <!-- Callback Test -->
            <div class="bg-white rounded-lg shadow p-6">
                <h2 class="text-xl font-semibold mb-4">Test Callback</h2>
                <p class="text-gray-600 mb-4">Test the callback page that processes JWT tokens from URL fragments.</p>
                <a href="/auth/callback#access_token=test-jwt-token&token_type=Bearer"
                   class="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded">
                    Test Callback with Fake Token
                </a>
            </div>
            
            <!-- Instructions -->
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                <h3 class="text-lg font-semibold text-yellow-800 mb-2">🔍 Testing Instructions</h3>
                <ol class="text-sm text-sm text-yellow-700 space-y-1">
                    <li>1. Open browser console (F12) to see detailed logs</li>
                    <li>2. First login with Google, GitHub, Discord, or Microsoft</li>
                    <li>3. Then click "Test Token Refresh" to see the flow</li>
                    <li>4. Check console logs for every step of the process</li>
                </ol>
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
//
//	Auth service handles Google OAuth -> Returns to our callback with JWT
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

// DiscordLoginHandler handles Discord OAuth login
// Flow: User clicks "Login with Discord" -> Redirect to our auth service ->
//
//	Auth service handles Discord OAuth -> Returns to our callback with JWT
func (h *AuthHandler) DiscordLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 DISCORD LOGIN: Starting Discord OAuth flow\n")
	fmt.Printf("🔐 DISCORD LOGIN: AuthServiceURL = %s\n", h.Config.AuthServiceURL)
	fmt.Printf("🔐 DISCORD LOGIN: RedirectURL = %s\n", h.Config.RedirectURL)

	// Redirect to our auth microservice with redirect_uri parameter
	// The auth service will handle the actual Discord OAuth flow
	authURL := fmt.Sprintf("%s/auth/discord?redirect_uri=%s/auth/callback",
		h.Config.AuthServiceURL, h.Config.RedirectURL)

	fmt.Printf("🔐 DISCORD LOGIN: Redirecting to: %s\n", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// MicrosoftLoginHandler handles Microsoft OAuth login
// Flow: User clicks "Login with Microsoft" -> Redirect to our auth service ->
//
//	Auth service handles Microsoft OAuth -> Returns to our callback with JWT
func (h *AuthHandler) MicrosoftLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 MICROSOFT LOGIN: Starting Microsoft OAuth flow\n")
	fmt.Printf("🔐 MICROSOFT LOGIN: AuthServiceURL = %s\n", h.Config.AuthServiceURL)
	fmt.Printf("🔐 MICROSOFT LOGIN: RedirectURL = %s\n", h.Config.RedirectURL)

	// Redirect to our auth microservice with redirect_uri parameter
	// The auth service will handle the actual Microsoft OAuth flow
	authURL := fmt.Sprintf("%s/auth/microsoft?redirect_uri=%s/auth/callback",
		h.Config.AuthServiceURL, h.Config.RedirectURL)

	fmt.Printf("🔐 MICROSOFT LOGIN: Redirecting to: %s\n", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// AuthCallbackHandler handles the OAuth callback
// Flow: Google redirects here with JWT in URL fragment (#access_token=...)
//
//	Client-side JS extracts token and calls /api/auth/set-session
func (h *AuthHandler) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 CALLBACK: === OAuth callback STARTED ===\n")
	fmt.Printf("🔐 CALLBACK: URL = %s\n", r.URL.String())
	fmt.Printf("🔐 CALLBACK: Query params = %v\n", r.URL.Query())
	fmt.Printf("🔐 CALLBACK: Fragment = %s\n", r.URL.Fragment)

	fmt.Printf("🔐 CALLBACK: Setting content type and rendering template...\n")
	w.Header().Set("Content-Type", "text/html")

	// STEP 2: Render callback page with JavaScript to extract JWT from URL fragment
	// The fragment (#access_token=...) is not sent to server, so JS must handle it
	component := layouts.Layout("Authenticating", "Authentication processing page for OAuth callback and session establishment.", layouts.NavigationLoggedOut(), pages.AuthCallbackContent())

	fmt.Printf("🔐 CALLBACK: About to render component...\n")
	component.Render(r.Context(), w)
	fmt.Printf("🔐 CALLBACK: Component rendered successfully\n")
	fmt.Printf("🔐 CALLBACK: === OAuth callback COMPLETED ===\n")
}

func (h *AuthHandler) SetSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 SESSION: === Set session STARTED ===\n")
	fmt.Printf("🔐 SESSION: Content-Type: %s\n", r.Header.Get("Content-Type"))

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("🔐 SESSION: Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	if req.SessionID == "" {
		fmt.Printf("🔐 SESSION: Missing session_id\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing session_id",
		})
		return
	}

	fmt.Printf("🔐 SESSION: Session ID received, length: %d\n", len(req.SessionID))

	// Set session_id cookie (replaces session_token)
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    req.SessionID,
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	// Set session_id cookie
	http.SetCookie(w, sessionCookie)

	fmt.Printf("🔐 SESSION: Session ID cookie set successfully:")
	fmt.Printf("🔐 SESSION: - session_id cookie, Length: %d\n", len(sessionCookie.Value))

	fmt.Printf("🔐 SESSION: SUCCESS: Server session established")
	fmt.Printf("🔐 SESSION: === Set session COMPLETED ===\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Server session established successfully",
	})
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

// TestCreateSessionHandler tests the session creation endpoint
func (h *AuthHandler) TestCreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🧪 TEST: Testing /auth/session/create endpoint...\n")
	
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("🧪 TEST: Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	if req.Code == "" {
		fmt.Printf("🧪 TEST: Missing authorization code\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing authorization code",
		})
		return
	}

	fmt.Printf("🧪 TEST: Authorization code received: %s\n", req.Code)

	// Test the new CreateSession function
	response, err := h.AuthService.CreateSession(req.Code)
	if err != nil {
		fmt.Printf("🧪 TEST: CreateSession failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("🧪 TEST: Session creation response: %+v\n", response)
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"response": response,
	})
}

// ExchangeCodeHandler exchanges OAuth authorization code for tokens
func (h *AuthHandler) ExchangeCodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔄 CODE: === Exchange authorization code STARTED ===\n")
	fmt.Printf("🔄 CODE: Request URL: %s\n", r.URL.String())

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("🔄 CODE: Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	if req.Code == "" {
		fmt.Printf("🔄 CODE: Missing authorization code\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing authorization code",
		})
		return
	}

	fmt.Printf("🔄 CODE: Authorization code received, length: %d\n", len(req.Code))

	// Exchange code for tokens via auth service
	fmt.Printf("🔄 CODE: Calling auth service to exchange code for tokens...\n")
	tokensResp, err := h.AuthService.ExchangeCodeForTokens(req.Code)
	if err != nil {
		fmt.Printf("❌ CODE: Auth service failed: %v\n", err)
		fmt.Printf("❌ CODE: Error type: %T\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
			"error_type": fmt.Sprintf("%T", err),
		})
		return
	}

	if !tokensResp.Success {
		fmt.Printf("❌ CODE: Token exchange failed: %s\n", tokensResp.Error)
		fmt.Printf("❌ CODE: Response: %+v\n", tokensResp)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": tokensResp.Error,
			"success": tokensResp.Success,
		})
		return
	}

	fmt.Printf("✅ CODE: Auth service returned success: %v\n", tokensResp.Success)
	fmt.Printf("🔄 CODE: Auth response: %+v\n", tokensResp)

	// Generate session ID for server session (in real app, this would come from auth service)
	sessionID := fmt.Sprintf("sess_%d_%x", time.Now().UnixNano(), time.Now().Unix())

	// Set session_id cookie for server sessions
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	// Set session_id cookie
	http.SetCookie(w, sessionCookie)

	fmt.Printf("✅ CODE: Session ID cookie set successfully: %s\n", sessionID)
	fmt.Printf("🔄 CODE: === Token exchange COMPLETED ===\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tokens exchanged successfully",
	})
}

// GetUserInfo returns current user information for server-side rendering
func (h *AuthHandler) GetUserInfo(r *http.Request) layouts.UserInfo {
	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}
	}

	// Get user info from auth microservice
	userResp, err := h.AuthService.ValidateUser(cookie.Value)
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}
	}

	return layouts.UserInfo{
		LoggedIn: userResp.Success,
		Name:     userResp.Name,
		Email:    userResp.Email,
		Picture:  userResp.Picture,
	}
}

