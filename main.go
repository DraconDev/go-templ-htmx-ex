package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

// getUserFromJWT gets user info using local JWT validation (5-10ms, no API call)
func getUserFromJWT(r *http.Request) templates.UserInfo {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return templates.UserInfo{LoggedIn: false}
	}
	
	return validateJWTWithRealData(cookie.Value)
}

// hasSessionToken checks if user has a session token cookie (fast, no API call)
func hasSessionToken(r *http.Request) bool {
	_, err := r.Cookie("session_token")
	return err == nil
}

// HTTP Handlers

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Use local JWT validation for EVERY page (fast + consistent)
	userInfo := getUserFromJWT(r)
	
	var navigation templ.Component
	if userInfo.LoggedIn {
		// Fast local validation: 5-10ms for real user data
		navigation = templates.NavigationLoggedIn(userInfo)
	} else {
		navigation = templates.NavigationLoggedOut()
	}
	
	component := templates.Layout("Home", navigation, templates.HomeContent())
	component.Render(r.Context(), w)
}

// validateJWTWithRealData validates JWT and returns real user data
func validateJWTWithRealData(token string) templates.UserInfo {
	log.Printf("🔍 JWT: Validating token, length: %d", len(token))
	
	if token == "" {
		log.Printf("🔍 JWT: Empty token")
		return templates.UserInfo{LoggedIn: false}
	}
	
	// Parse JWT to get real user data
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		log.Printf("🔍 JWT: Invalid token format, parts: %d", len(parts))
		return templates.UserInfo{LoggedIn: false}
	}
	
	log.Printf("🔍 JWT: Token split successful")
	
	// Decode payload (the middle part)
	payload, err := jwtBase64URLDecode(parts[1])
	if err != nil {
		log.Printf("🔍 JWT: Base64 decode failed: %v", err)
		return templates.UserInfo{LoggedIn: false}
	}
	
	log.Printf("🔍 JWT: Base64 decode successful, payload length: %d", len(payload))
	
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
		log.Printf("🔍 JWT: JSON unmarshal failed: %v", err)
		return templates.UserInfo{LoggedIn: false}
	}
	
	log.Printf("🔍 JWT: Claims parsed - Name: %s, Email: %s, Iss: %s", claims.Name, claims.Email, claims.Iss)
	
	// Check if token is still valid (not expired)
	if claims.Exp < time.Now().Unix() {
		log.Printf("🔍 JWT: Token expired - Exp: %d, Now: %d", claims.Exp, time.Now().Unix())
		return templates.UserInfo{LoggedIn: false}
	}
	
	log.Printf("🔍 JWT: Token not expired")
	
	// Check issuer to make sure it's from our auth service
	if claims.Iss != "auth-ms" {
		log.Printf("🔍 JWT: Invalid issuer: %s (expected: auth-ms)", claims.Iss)
		return templates.UserInfo{LoggedIn: false}
	}
	
	log.Printf("🔍 JWT: Validation successful - User: %s", claims.Name)
	
	// Return real user data!
	return templates.UserInfo{
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

func profileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Use local JWT validation for consistency (5-10ms everywhere)
	userInfo := getUserFromJWT(r)
	if !userInfo.LoggedIn {
		// Redirect to home if not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Create profile content with real user data
	navigation := templates.NavigationLoggedIn(userInfo)
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
