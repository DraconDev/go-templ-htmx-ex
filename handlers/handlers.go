package handlers

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/DraconDev/go-templ-htmx-ex/templates"
)

// Handlers contains all the handlers
type Handlers struct {
	AuthHandler *AuthHandler
}

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Check for JWT token to determine authentication status
	var userInfo templates.UserInfo
	cookie, err := r.Cookie("session_token")
	if err == nil && cookie.Value != "" {
		// User has a JWT token - get real user info from auth service
		userInfo = templates.UserInfo{LoggedIn: true, Name: "User", Picture: "https://via.placeholder.com/40"}
	} else {
		userInfo = templates.UserInfo{LoggedIn: false}
	}
	
	// Choose navigation based on auth status
	var nav templ.Component
	if userInfo.LoggedIn {
		nav = templates.NavigationLoggedIn(userInfo)
	} else {
		nav = templates.NavigationLoggedOut()
	}
	
	component := templates.Layout("Home", nav, templates.HomeContent())
	component.Render(r.Context(), w)
}

// ProfileHandler handles the user profile page
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Get JWT token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		// Redirect to home if not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// User has valid JWT token - show logged-in navigation
	userInfo := templates.UserInfo{LoggedIn: true}
	component := templates.Layout("Profile", templates.NavigationLoggedIn(userInfo), templates.ProfileContent("", "", ""))
	component.Render(r.Context(), w)
}

// LoginHandler handles the login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	component := templates.Layout("Login", templates.NavigationLoggedOut(), templates.LoginContent())
	component.Render(r.Context(), w)
}
