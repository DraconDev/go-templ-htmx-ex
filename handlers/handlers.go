package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/DraconDev/go-templ-htmx-ex/middleware"
	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
	"github.com/DraconDev/go-templ-htmx-ex/templates/pages"
)

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

// GetUserFromJWT gets user info using local JWT validation (5-10ms, no API call)
func GetUserFromJWT(r *http.Request) layouts.UserInfo {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}
	}

	return validateJWTWithRealData(cookie.Value)
}

// validateJWTWithRealData validates JWT and returns real user data
func validateJWTWithRealData(token string) layouts.UserInfo {
	if token == "" {
		return layouts.UserInfo{LoggedIn: false}
	}

	// Parse JWT to get real user data
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return layouts.UserInfo{LoggedIn: false}
	}

	// Decode payload (the middle part)
	payload, err := jwtBase64URLDecode(parts[1])
	if err != nil {
		return layouts.UserInfo{LoggedIn: false}
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
		return layouts.UserInfo{LoggedIn: false}
	}

	// Check if token is still valid (not expired)
	if claims.Exp < time.Now().Unix() {
		return layouts.UserInfo{LoggedIn: false}
	}

	// Check issuer to make sure it's from our auth service
	if claims.Iss != "auth-ms" {
		return layouts.UserInfo{LoggedIn: false}
	}

	// Return real user data!
	return layouts.UserInfo{
		LoggedIn: true,
		Name:     claims.Name,
		Email:    claims.Email,
		Picture:  claims.Picture,
	}
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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Get user info from middleware context
	userInfo := middleware.GetUserFromContext(r)
	
	fmt.Printf("🏠 HOME: User info - LoggedIn: %v, Name: %s, Email: %s\n",
		userInfo.LoggedIn, userInfo.Name, userInfo.Email)

	var navigation templ.Component
	if userInfo.LoggedIn {
		fmt.Printf("🏠 HOME: Rendering NavigationLoggedIn\n")
		navigation = layouts.NavigationLoggedIn(userInfo)
	} else {
		fmt.Printf("🏠 HOME: Rendering NavigationLoggedOut\n")
		navigation = layouts.NavigationLoggedOut()
	}

	component := layouts.Layout("Home", "Production-ready startup platform with Google OAuth, PostgreSQL database, and admin dashboard. Built with Go + HTMX + Templ.", navigation, pages.HomeContent())
	component.Render(r.Context(), w)
}

// ProfileHandler handles the user profile page
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Use local JWT validation for consistency (5-10ms everywhere)
	userInfo := GetUserFromJWT(r)
	if !userInfo.LoggedIn {
		// Redirect to home if not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Create profile content with real user data
	navigation := layouts.NavigationLoggedIn(userInfo)
	component := layouts.Layout("Profile", "User profile page with authentication details and account management.", navigation, pages.ProfileContent(userInfo.Name, userInfo.Email, userInfo.Picture))
	component.Render(r.Context(), w)
}

// LoginHandler handles the login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	component := layouts.Layout("Login", "Secure authentication page with Google OAuth integration for user access.", layouts.NavigationLoggedOut(), pages.LoginContent())
	component.Render(r.Context(), w)
}
