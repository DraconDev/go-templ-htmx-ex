package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/DraconDev/go-templ-htmx-ex/templates"
)

// UserSession represents a logged-in user session
type UserSession struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	LoggedIn bool   `json:"logged_in"`
}

// Config holds application configuration
type Config struct {
	ServerPort    string
	AuthServiceURL string
	RedirectURL    string
}

var (
	config = &Config{
		ServerPort:     getEnvOrDefault("PORT", "8081"),
		AuthServiceURL: getEnvOrDefault("AUTH_SERVICE_URL", "http://localhost:8080"),
		RedirectURL:    getEnvOrDefault("REDIRECT_URL", "http://localhost:8081"),
	}
)

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Create router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/health", healthHandler).Methods("GET")

	// Google OAuth login routes
	router.HandleFunc("/auth/google", googleLoginHandler).Methods("GET")
	router.HandleFunc("/auth/callback", authCallbackHandler).Methods("GET")
	
	// User profile page
	router.HandleFunc("/profile", profileHandler).Methods("GET")
	
	// Session management
	router.HandleFunc("/api/auth/validate", authValidateSessionHandler).Methods("POST")
	router.HandleFunc("/api/auth/logout", authLogoutHandler).Methods("POST")
	router.HandleFunc("/api/auth/user", authGetUserHandler).Methods("GET")
	router.HandleFunc("/api/auth/health", authHealthCheckHandler).Methods("GET")

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.ServerPort),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", config.ServerPort)
		log.Printf("Visit http://localhost:%s to access the application", config.ServerPort)
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

// HTTP Handlers

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	component := templates.Layout("Home", templates.HomeContent())
	component.Render(r.Context(), w)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Get current user session
	resp, err := callAuthService(fmt.Sprintf("%s/auth/userinfo", config.AuthServiceURL), map[string]string{
		"token": getSessionToken(r),
	})
	if err != nil {
		// Redirect to home if not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	if !resp.Success {
		// Redirect to home if not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Create profile content with user data
	component := templates.Layout("Profile", templates.ProfileContent(resp.Name, resp.Email, resp.Picture))
	component.Render(r.Context(), w)
}

// Helper function to get session token from cookie
func getSessionToken(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

// Google OAuth Handlers

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to the auth microservice's Google OAuth endpoint
	authURL := fmt.Sprintf("%s/auth/google?redirect_uri=%s/auth/callback", 
		config.AuthServiceURL, config.RedirectURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the callback from the auth microservice
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// Call the auth microservice to exchange code for token
	tokenResp, err := callAuthService(fmt.Sprintf("%s/auth/google/callback", config.AuthServiceURL), map[string]string{
		"code": code,
	})
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}

	// Set session cookie with the JWT token
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    tokenResp.Token,
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusFound)
}

func authValidateSessionHandler(w http.ResponseWriter, r *http.Request) {
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
	resp, err := callAuthService(fmt.Sprintf("%s/auth/validate", config.AuthServiceURL), map[string]string{
		"token": cookie.Value,
	})
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
		"valid":    resp.Success,
		"user_id":  resp.UserID,
		"email":    resp.Email,
		"name":     resp.Name,
		"picture":  resp.Picture,
		"status":   "validated",
	})
}

func authLogoutHandler(w http.ResponseWriter, r *http.Request) {
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

func authGetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"logged_in": false,
		})
		return
	}

	// Get user info from auth microservice
	resp, err := callAuthService(fmt.Sprintf("%s/auth/userinfo", config.AuthServiceURL), map[string]string{
		"token": cookie.Value,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"logged_in": false,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logged_in": resp.Success,
		"user_id":   resp.UserID,
		"email":     resp.Email,
		"name":      resp.Name,
		"picture":   resp.Picture,
	})
}

func authHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Simple health check
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "main-app",
	})
}

// Helper function to call auth service
func callAuthService(endpoint string, params map[string]string) (*AuthResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	// Create form data
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, value)
	}
	
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}
	
	return &authResp, nil
}

// AuthResponse represents the response from the auth service
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
	Error   string `json:"error,omitempty"`
}
