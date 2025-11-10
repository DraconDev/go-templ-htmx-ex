package handlers

import (
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/templates"
)

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Get user authentication status for dynamic navigation
	var userInfo templates.UserInfo
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// User has a session token - try to validate with auth service
		userInfo = templates.UserInfo{LoggedIn: true}
		// In a real implementation, we would validate the token here
		// For now, having a session token indicates login
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

	// Get session token from cookie
	_, err := r.Cookie("session_token")
	if err != nil {
		// Redirect to home if not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Get user data from auth service
	// For now, pass the token data - the template and JavaScript will handle the rest
	// This maintains the working behavior from the original code
	component := templates.Layout("Profile", templates.NavigationLoggedOut(), templates.ProfileContent("", "", ""))
	component.Render(r.Context(), w)
}

// LoginHandler handles the login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	component := templates.Layout("Login", templates.NavigationLoggedOut(), templates.LoginContent())
	component.Render(r.Context(), w)
}
