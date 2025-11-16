package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
	"github.com/DraconDev/go-templ-htmx-ex/templates/pages"
)

// =============================================================================
// AUTHENTICATION HANDLER
// =============================================================================
// This handler manages the complete OAuth + Server Session authentication flow for the app:
// 1. OAuth redirects to external providers (Google, GitHub, Discord, Microsoft)
// 2. Callback processing to extract authorization codes from URL
// 3. Server session management with HTTP-only cookies
// 4. Session validation through middleware
// =============================================================================

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	Config      *config.Config // App configuration
	AuthService *auth.Service  // Auth service for session management
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(config *config.Config) *AuthHandler {
	return &AuthHandler{
		Config: config,
	}
}



// OAUTH LOGIN FLOWS
// =============================================================================

// LoginHandler handles OAuth login for any provider
// Flow: User clicks "Login with [Provider]" -> Redirect to our auth service ->
//
//	Auth service handles OAuth -> Returns to our callback with session token
//
// Usage: /auth/login?provider=google|github|discord|microsoft
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get provider from query parameter
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		fmt.Printf("🔐 LOGIN ERROR: Missing provider parameter\n")
		http.Redirect(w, r, "/login?error=missing_provider", http.StatusFound)
		return
	}

	// Validate provider
	validProviders := map[string]bool{
		"google":    true,
		"github":    true,
		"discord":   true,
		"microsoft": true,
	}

	if !validProviders[provider] {
		fmt.Printf("🔐 LOGIN ERROR: Invalid provider '%s'\n", provider)
		http.Redirect(w, r, "/login?error=invalid_provider", http.StatusFound)
		return
	}

	fmt.Printf("🔐 LOGIN: Starting %s OAuth flow\n", provider)
	fmt.Printf("🔐 LOGIN: AuthServiceURL = %s\n", h.Config.AuthServiceURL)
	fmt.Printf("🔐 LOGIN: RedirectURL = %s\n", h.Config.RedirectURL)

	// Redirect to our auth microservice with redirect_uri parameter
	// The auth service will handle the actual OAuth flow for the specified provider
	authURL := fmt.Sprintf("%s/auth/%s?redirect_uri=%s/auth/callback",
		h.Config.AuthServiceURL, provider, h.Config.RedirectURL)

	fmt.Printf("🔐 LOGIN: Redirecting to: %s\n", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// AuthCallbackHandler handles the OAuth callback
// Flow: OAuth provider redirects here with authorization code in URL
//
//	Client-side JS extracts token and calls /api/auth/set-session
func (h *AuthHandler) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 CALLBACK: === OAuth callback STARTED ===\n")
	fmt.Printf("🔐 CALLBACK: URL = %s\n", r.URL.String())
	fmt.Printf("🔐 CALLBACK: Query params = %v\n", r.URL.Query())
	fmt.Printf("🔐 CALLBACK: Fragment = %s\n", r.URL.Fragment)

	fmt.Printf("🔐 CALLBACK: Setting content type and rendering template...\n")
	w.Header().Set("Content-Type", "text/html")

	// STEP 2: Render callback page with JavaScript to extract session token from URL fragment
	// The fragment (#access_token=...) is not sent to server, so JS must handle it
	component := layouts.Layout("Authenticating", "Authentication processing page for OAuth callback and session establishment.", layouts.NavigationLoggedOut(), pages.AuthCallbackContent())

	fmt.Printf("🔐 CALLBACK: About to render component...\n")
	if err := component.Render(r.Context(), w); err != nil {
		fmt.Printf("🚨 CALLBACK: Error rendering component: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
		MaxAge:   2592000, // 30 days (server-side validation handles real security)
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

// LogoutHandler handles user logout
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
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

// ExchangeCodeHandler exchanges OAuth authorization code for tokens
func (h *AuthHandler) ExchangeCodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔄 CODE: === Exchange authorization code STARTED ===\n")
	fmt.Printf("🔄 CODE: Request URL: %s\n", r.URL.String())

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		AuthCode string `json:"auth_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("🔄 CODE: Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	if req.AuthCode == "" {
		fmt.Printf("🔄 CODE: Missing authorization code\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing authorization code",
		})
		return
	}

	fmt.Printf("🔄 CODE: Authorization code received, length: %d\n", len(req.AuthCode))

	// Create session from authorization code (returns JSON with all info)
	fmt.Printf("🔄 CODE: Creating session from authorization code...\n")
	sessionData, err := h.AuthService.CreateSession(req.AuthCode)
	if err != nil {
		fmt.Printf("❌ CODE: Auth service failed: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("✅ CODE: Auth service returned session data: %+v\n", sessionData)

	// Extract session_id from the response
	var session_id string
	if sid, exists := sessionData["session_id"]; exists {
		if sidStr, ok := sid.(string); ok {
			session_id = sidStr
		}
	}

	if session_id == "" {
		fmt.Printf("❌ CODE: No session_id in auth response\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "No session_id received from auth service",
		})
		return
	}

	// Set session_id cookie for server sessions
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    session_id,
		Path:     "/",
		MaxAge:   2592000, // 30 days (server-side validation handles real security)
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	// Set session_id cookie
	http.SetCookie(w, sessionCookie)

	fmt.Printf("✅ CODE: Session token cookie set successfully (length: %d)\n", len(session_id))
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
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}
	}

	// Get user info from auth microservice
	userResp, err := h.AuthService.GetUserInfo(cookie.Value)
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
