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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

func microserviceTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	component := templates.Layout("Microservice Testing", templates.MicroserviceTestContent())
	component.Render(r.Context(), w)
}

func serviceTestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]

	w.Header().Set("Content-Type", "text/html")
	component := templates.Layout(fmt.Sprintf("Testing %s", serviceName), templates.ServiceTestContent(serviceName))
	component.Render(r.Context(), w)
}

// API Handlers for HTMX

func servicesAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Simple response for now - we'll implement real service discovery later
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
		"services": [
			{"name": "User Service", "url": "http://localhost:8001", "status": "unknown"},
			{"name": "Order Service", "url": "http://localhost:8002", "status": "unknown"},
			{"name": "Payment Service", "url": "http://localhost:8003", "status": "unknown"},
			{"name": "Notification Service", "url": "http://localhost:8004", "status": "unknown"}
		]
	}`))
}

func runTestAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	serviceURL := r.FormValue("service_url")
	testType := r.FormValue("test_type")

	// Create test result component
	component := templates.TestResult(serviceURL, testType, "success", "Test completed successfully!")
	component.Render(r.Context(), w)
}

// Auth Service Handlers

func authHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Use the real auth client to check health
	resp, err := authClient.HealthCheck()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "unavailable",
			"error":      err.Error(),
			"url":        "cerberus-auth-ms-548010171143.europe-west1.run.app:443",
			"service":    "auth",
			"full_error": fmt.Sprintf("%+v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    resp.Status,
		"url":       "https://cerberus-auth-ms-548010171143.europe-west1.run.app",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   resp.Message,
		"service":   "auth",
	})
}

func authLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse JSON request body instead of form
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Use the real auth client to login
	authResp, err := authClient.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Return successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       authResp.UserID,
		"session_token": authResp.SessionToken,
		"email":         authResp.Email,
		"status":        "success",
	})
}

func authRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse JSON request body
	var registerReq struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		ProjectID string `json:"project_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Use the real auth client to register
	authResp, err := authClient.Register(registerReq.Email, registerReq.Password, registerReq.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Return successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       authResp.UserID,
		"session_token": authResp.SessionToken,
		"status":        "success",
	})
}

func authValidateSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse JSON request body
	var validateReq struct {
		SessionToken string `json:"session_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&validateReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Use the real auth client to validate session
	authResp, err := authClient.ValidateSession(validateReq.SessionToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
			"valid": false,
		})
		return
	}

	// Return validation result
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":     authResp.UserID,
		"valid":       authResp.Valid,
		"project_ids": authResp.ProjectIDs,
		"status":      "validated",
	})
}

func authGetUserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse JSON request body
	var userDetailsReq struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userDetailsReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Use the auth client to get user details
	authResp, err := authClient.GetUserDetails(userDetailsReq.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Return user details
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": authResp.UserID,
		"email":   authResp.Email,
		"status":  "success",
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
