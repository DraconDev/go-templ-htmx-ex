package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
        
        // Test the complete token refresh flow with detailed logging
        function testTokenRefresh() {
            console.log('🔄 TOKEN REFRESH TEST: === STARTED ===');
            const resultDiv = document.getElementById('refresh-result');
            
            // STEP 1: Log that we're starting the test
            console.log('🔄 TOKEN REFRESH TEST: Step 1 - Starting refresh test...');
            resultDiv.innerHTML = '<p class="text-blue-600">🔄 Starting token refresh test...</p>';
            
            // STEP 2: Check if refresh_token cookie exists
            const cookies = document.cookie.split(';');
            const refreshCookie = cookies.find(cookie => cookie.trim().startsWith('refresh_token='));
            console.log('🔄 TOKEN REFRESH TEST: Step 2 - Checking for refresh_token cookie...');
            console.log('🔄 TOKEN REFRESH TEST: All cookies:', cookies);
            console.log('🔄 TOKEN REFRESH TEST: refresh_token cookie found:', !!refreshCookie);
            
            if (!refreshCookie) {
                console.log('❌ TOKEN REFRESH TEST: No refresh_token cookie found - user may not be logged in');
                resultDiv.innerHTML = '<p class="text-red-600">❌ No refresh_token cookie found. User may not be logged in.</p>';
                return;
            }
            
            // STEP 3: Call the refresh endpoint
            console.log('🔄 TOKEN REFRESH TEST: Step 3 - Calling /api/auth/refresh...');
            resultDiv.innerHTML = '<p class="text-blue-600">🔄 Calling refresh endpoint...</p>';
            
            fetch('/api/auth/refresh', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => {
                console.log('🔄 TOKEN REFRESH TEST: Step 4 - Got response from server');
                console.log('🔄 TOKEN REFRESH TEST: Response status:', response.status);
                console.log('🔄 TOKEN REFRESH TEST: Response headers:', response.headers);
                
                return response.json().then(data => {
                    console.log('🔄 TOKEN REFRESH TEST: Response data:', data);
                    
                    if (response.ok && data.success) {
                        console.log('✅ TOKEN REFRESH TEST: SUCCESS - Token refreshed successfully!');
                        console.log('🔄 TOKEN REFRESH TEST: New session token should be set in cookies');
                        
                        // Check if new session_token cookie was set
                        const newSessionCookie = document.cookie.split(';').find(cookie => cookie.trim().startsWith('session_token='));
                        console.log('🔄 TOKEN REFRESH TEST: New session_token cookie set:', !!newSessionCookie);
                        
                        resultDiv.innerHTML = '<p class="text-green-600">✅ SUCCESS: Token refreshed!</p>' +
                            '<p class="text-sm text-gray-600 mt-2">New session_token cookie should now be set.</p>' +
                            '<p class="text-sm text-gray-600">Check browser console for detailed logs.</p>';
                    } else {
                        console.log('❌ TOKEN REFRESH TEST: FAILED - Server returned error');
                        resultDiv.innerHTML = '<p class="text-red-600">❌ ERROR: ' + (data.error || 'Unknown error') + '</p>';
                    }
                });
            })
            .catch(error => {
                console.log('❌ TOKEN REFRESH TEST: NETWORK ERROR');
                console.log('❌ TOKEN REFRESH TEST: Error:', error);
                resultDiv.innerHTML = '<p class="text-red-600">❌ NETWORK ERROR: ' + error.message + '</p>';
            });
            
            console.log('🔄 TOKEN REFRESH TEST: === TEST INITIATED ===');
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
                    <a href="/auth/google" class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">Login with Google</a>
                    <a href="/auth/github" class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded">Login with GitHub</a>
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
                <ol class="text-sm text-yellow-700 space-y-1">
                    <li>1. Open browser console (F12) to see detailed logs</li>
                    <li>2. First login with Google or GitHub</li>
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
		Token string `json:"token"`
		// NOTE: refresh_token should NEVER be sent from client-side JavaScript
		// The auth service should set refresh_token cookie directly
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
		fmt.Printf("🔐 SESSION: Missing access token in request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing access token",
		})
		return
	}

	fmt.Printf("🔐 SESSION: Access token received, length: %d\n", len(req.Token))
	fmt.Printf("🔐 SESSION: IMPORTANT: Setting ONLY session_token cookie from access_token")
	fmt.Printf("🔐 SESSION: Refresh token should be set by auth service directly as HTTP-only cookie")

	// Set session cookie with the access token (JWT)
	sessionCookie := &http.Cookie{
		Name:     "session_token",
		Value:    req.Token,
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	http.SetCookie(w, sessionCookie)
	fmt.Printf("🔐 SESSION: Session cookie set - Name: %s, Length: %d\n", sessionCookie.Name, len(sessionCookie.Value))

	fmt.Printf("🔐 SESSION: SUCCESS: Session established with access token")
	fmt.Printf("🔐 SESSION: === Set session COMPLETED ===\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Session set successfully with access token",
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

// RefreshTokenHandler handles token refresh requests
// Flow: Frontend calls when JWT expires ->
//
//	Server reads refresh_token cookie ->
//	Calls auth service for new JWT ->
//	**Sets new session_token cookie automatically**
func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔄 REFRESH: === Token refresh STARTED ===\n")
	fmt.Printf("🔄 REFRESH: Request URL: %s\n", r.URL.String())
	fmt.Printf("🔄 REFRESH: Request headers: %v\n", r.Header)
	
	w.Header().Set("Content-Type", "application/json")

	// STEP 1: Get refresh token from HTTP-only cookie (automatically sent by browser)
	fmt.Printf("🔄 REFRESH: Looking for refresh_token cookie...\n")
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Printf("❌ REFRESH: No refresh_token cookie found: %v\n", err)
		fmt.Printf("🔄 REFRESH: All cookies: %v\n", r.Cookies())
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "No refresh token found",
		})
		return
	}

	fmt.Printf("✅ REFRESH: Found refresh_token cookie, length: %d\n", len(cookie.Value))
	fmt.Printf("🔄 REFRESH: Cookie details - Name: %s, Domain: %s, Path: %s, MaxAge: %d\n",
		cookie.Name, cookie.Domain, cookie.Path, cookie.MaxAge)

	// STEP 2: Call auth service to refresh token using the refresh token
	fmt.Printf("🔄 REFRESH: Calling auth service to refresh token...\n")
	userResp, err := h.AuthService.RefreshToken(cookie.Value)
	if err != nil {
		fmt.Printf("❌ REFRESH: Auth service failed: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("✅ REFRESH: Auth service returned success: %v\n", userResp.Success)
	fmt.Printf("🔄 REFRESH: New token length: %d\n", len(userResp.Token))

	// STEP 3: **CRITICAL** - Set the new JWT cookie for the user
	// This replaces the expired session_token with a fresh one
	fmt.Printf("🔄 REFRESH: Setting new session_token cookie...\n")
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    userResp.Token, // NEW JWT from auth service
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production
	})

	fmt.Printf("✅ REFRESH: New session_token cookie set successfully\n")
	fmt.Printf("🔄 REFRESH: === Token refresh COMPLETED ===\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Token refreshed successfully",
	})
}
