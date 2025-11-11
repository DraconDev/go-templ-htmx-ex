package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/templates"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	AuthService *auth.Service
	Config      *config.Config
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *auth.Service, config *config.Config) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Config:      config,
	}
}

// GoogleLoginHandler handles Google OAuth login
func (h *AuthHandler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 GOOGLE LOGIN: Starting Google OAuth flow\n")
	
	// Check if auth secret is configured
	if h.Config.AuthSecret == "" {
		fmt.Printf("🔐 GOOGLE LOGIN: Auth secret not configured\n")
		http.Error(w, "Auth secret not configured", http.StatusInternalServerError)
		return
	}
	
	// Create the auth service URL with secret parameter
	authURL := fmt.Sprintf("%s/auth/google?redirect_uri=%s&secret=%s",
		h.Config.AuthServiceURL, 
		url.QueryEscape(fmt.Sprintf("%s/auth/callback", h.Config.RedirectURL)),
		h.Config.AuthSecret)
	
	fmt.Printf("🔐 GOOGLE LOGIN: Redirecting to: %s\n", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// GitHubLoginHandler handles GitHub OAuth login
func (h *AuthHandler) GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 GITHUB LOGIN: Starting GitHub OAuth flow\n")
	
	// Check if auth secret is configured
	if h.Config.AuthSecret == "" {
		fmt.Printf("🔐 GITHUB LOGIN: Auth secret not configured\n")
		http.Error(w, "Auth secret not configured", http.StatusInternalServerError)
		return
	}
	
	// Create the auth service URL with secret parameter
	authURL := fmt.Sprintf("%s/auth/github?redirect_uri=%s&secret=%s",
		h.Config.AuthServiceURL, 
		url.QueryEscape(fmt.Sprintf("%s/auth/callback", h.Config.RedirectURL)),
		h.Config.AuthSecret)
	
	fmt.Printf("🔐 GITHUB LOGIN: Redirecting to: %s\n", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// AuthCallbackHandler handles the OAuth callback
func (h *AuthHandler) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 CALLBACK: === OAuth callback STARTED ===\n")
	fmt.Printf("🔐 CALLBACK: URL = %s\n", r.URL.String())
	
	w.Header().Set("Content-Type", "text/html")
	component := templates.Layout("Authenticating", templates.NavigationLoggedOut(), templates.AuthCallbackContent())
	component.Render(r.Context(), w)
	fmt.Printf("🔐 CALLBACK: === OAuth callback COMPLETED ===\n")
}

// SetSessionHandler sets the user session from client-side JavaScript
func (h *AuthHandler) SetSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔐 SESSION: === Set session STARTED ===\n")
	
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
		"valid":    userResp.Success,
		"user_id":  userResp.UserID,
		"email":    userResp.Email,
		"name":     userResp.Name,
		"picture":  userResp.Picture,
		"status":   "validated",
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
