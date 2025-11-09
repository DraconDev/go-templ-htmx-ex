package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/handlers"
	"github.com/DraconDev/go-templ-htmx-ex/templates"
)

var authHandler *handlers.AuthHandler

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create auth service
	authService := auth.NewService(cfg)

	// Create auth handler
	authHandler = handlers.NewAuthHandler(authService, cfg)

	// Create router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/health", healthHandler).Methods("GET")

	// OAuth login routes
	router.HandleFunc("/auth/google", authHandler.GoogleLoginHandler).Methods("GET")
	router.HandleFunc("/auth/github", authHandler.GitHubLoginHandler).Methods("GET")
	router.HandleFunc("/auth/callback", authHandler.AuthCallbackHandler).Methods("GET")

	// User profile page
	router.HandleFunc("/profile", profileHandler).Methods("GET")

	// Session management
	router.HandleFunc("/api/auth/validate", authHandler.ValidateSessionHandler).Methods("POST")
	router.HandleFunc("/api/auth/logout", authHandler.LogoutHandler).Methods("POST")
	router.HandleFunc("/api/auth/user", authHandler.GetUserHandler).Methods("GET")
	router.HandleFunc("/api/auth/set-session", authHandler.SetSessionHandler).Methods("POST")
	router.HandleFunc("/api/auth/health", authHealthCheckHandler).Methods("GET")

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.ServerPort)
		log.Printf("Visit http://localhost:%s to access the application", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// Helper function to get session token from cookie
func getSessionToken(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// hasSessionToken checks if user has a session token cookie (fast, no API call)
func hasSessionToken(r *http.Request) bool {
	_, err := r.Cookie("session_token")
	return err == nil
}

// HTTP Handlers

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Fast cookie check for navigation (Hybrid approach - no API call)
	var navigation templ.Component
	if hasSessionToken(r) {
		navigation = templates.NavigationLoggedIn(templates.UserInfo{Name: "User"})
	} else {
		navigation = templates.NavigationLoggedOut()
	}
	
	component := templates.Layout("Home", navigation, templates.HomeContent())
	component.Render(r.Context(), w)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Server-side validation for protected pages (SSR approach)
	var navigation templ.Component
	var userInfo templates.UserInfo
	
	if hasSessionToken(r) {
		// Get real user data from auth service via authHandler
		userInfo = authHandler.GetUserInfo(r)
		if userInfo.LoggedIn {
			navigation = templates.NavigationLoggedIn(userInfo)
		} else {
			// Token invalid, redirect to home
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	component := templates.Layout("Profile", navigation, templates.ProfileContent(userInfo.Name, userInfo.Email, userInfo.Picture))
	component.Render(r.Context(), w)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

func authHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Simple health check
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
		"status": "healthy", 
		"timestamp": "` + time.Now().Format(time.RFC3339) + `",
		"service": "main-app"
	}`))
}
