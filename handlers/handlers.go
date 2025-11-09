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
	component := templates.Layout("Home", templates.HomeContent())
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

	// For now, show placeholder - the actual user data is handled by the template and JavaScript
	// This allows the frontend JavaScript to fetch user data dynamically
	component := templates.Layout("Profile", templates.ProfileContent("", "", ""))
	component.Render(r.Context(), w)
}