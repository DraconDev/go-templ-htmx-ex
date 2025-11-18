package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/internal/services"
	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/config"
	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
	"github.com/DraconDev/go-templ-htmx-ex/templates/pages"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	AuthService *services.AuthService // Communication with auth microservice
	Config      *config.Config        // App configuration
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Config:      config,
	}
}

// =============================================================================
// OAUTH LOGIN FLOW
// =============================================================================

// LoginHandler handles OAuth login for any provider
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		fmt.Printf("🔐 LOGIN ERROR: Missing provider parameter\n")
		http.Redirect(w, r, "/login?error=missing_provider", http.StatusFound)
		return
	}

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

	authURL := fmt.Sprintf("%s/auth/%s?redirect_uri=%s/auth/callback",
		h.Config.AuthServiceURL, provider, h.Config.RedirectURL)

	fmt.Printf("🔐 LOGIN: Redirecting to: %s\n", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// AuthCallbackHandler handles the OAuth callback
func (h *AuthHandler) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 CALLBACK: === OAuth callback STARTED ===\n")
	fmt.Printf("🔐 CALLBACK: URL = %s\n", r.URL.String())
	fmt.Printf("🔐 CALLBACK: Query params = %v\n", r.URL.Query())
	fmt.Printf("🔐 CALLBACK: Fragment = %s\n", r.URL.Fragment)

	fmt.Printf("🔐 CALLBACK: Setting content type and rendering template...\n")
	w.Header().Set("Content-Type", "text/html")

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

// ExchangeCodeHandler exchanges OAuth authorization code for session tokens
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

	// Exchange code for tokens via auth service
	fmt.Printf("🔄 CODE: Calling auth service to exchange code for tokens...\n")
	authResp, err := h.AuthService.ExchangeCodeForTokens(req.AuthCode)
	if err != nil {
		fmt.Printf("❌ CODE: Auth service failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if !authResp.Success {
		fmt.Printf("❌ CODE: Auth service returned failure: %s\n", authResp.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": authResp.Error,
		})
		return
	}

	fmt.Printf("✅ CODE: Auth service returned success: %v\n", authResp.Success)
	fmt.Printf("🔄 CODE: Auth response: %+v\n", authResp)

	// Set session_id cookie for server sessions
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    authResp.UserID, // Using UserID as the session identifier
		Path:     "/",
		MaxAge:   2592000, // 30 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	http.SetCookie(w, sessionCookie)

	fmt.Printf("✅ CODE: Session token cookie set successfully (length: %d)\n", len(authResp.UserID))
	fmt.Printf("🔄 CODE: === Token exchange COMPLETED ===\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tokens exchanged successfully",
	})
}

// SetSessionHandler handles setting a new session cookie
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

	// Set session_id cookie
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    req.SessionID,
		Path:     "/",
		MaxAge:   2592000, // 30 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	http.SetCookie(w, sessionCookie)

	fmt.Printf("✅ SESSION: Session ID cookie set successfully")
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