package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	ServerPort string
}

var (
	config = &Config{
		ServerPort: getEnvOrDefault("PORT", "8081"),
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

// Google OAuth Handlers

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to the auth microservice's Google OAuth endpoint
	authURL := "http://localhost:8080/auth/google?redirect_uri=http://localhost:8081/auth/callback"
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
	tokenResp, err := callAuthService("http://localhost:8080/auth/google/callback", map[string]string{
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
	resp, err := callAuthService("http://localhost:8080/auth/validate", map[string]string{
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
	resp, err := callAuthService("http://localhost:8080/auth/userinfo", map[string]string{
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

// Test Login Handlers

func testLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	component := templates.Layout("Test Login", templates.AuthTestContent())
	component.Render(r.Context(), w)
}

func authTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Determine action based on the URL path
	var action string
	switch r.URL.Path {
	case "/api/auth/login":
		action = "login"
	case "/api/auth/register":
		action = "register"
	case "/api/auth/validate":
		action = "validate"
	case "/api/auth/user-details":
		action = "user-details"
	default:
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	// Handle different authentication scenarios
	switch action {
	case "login":
		// Parse login credentials
		var loginReq struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Demo login - accepts any credentials
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":       "demo-user-123",
			"session_token": "demo-session-token-456",
			"email":         loginReq.Email,
			"status":        "success",
			"message":       "Test login successful (demo mode)",
		})

	case "register":
		// Parse registration credentials
		var registerReq struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			ProjectID string `json:"project_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Demo registration - accepts any credentials
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":       "demo-user-" + fmt.Sprintf("%d", time.Now().Unix()),
			"session_token": "demo-session-token-" + fmt.Sprintf("%d", time.Now().Unix()),
			"email":         registerReq.Email,
			"status":        "success",
			"message":       "Test registration successful (demo mode)",
		})

	case "validate":
		// Parse validation request
		var validateReq struct {
			SessionToken string `json:"session_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&validateReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Demo session validation
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":     "demo-user-123",
			"valid":       true,
			"project_ids": []string{"demo-project-1", "demo-project-2"},
			"status":      "validated",
			"message":     "Test session validation successful (demo mode)",
		})

	case "user-details":
		// Parse user details request
		var userDetailsReq struct {
			UserID string `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userDetailsReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Demo user details
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": userDetailsReq.UserID,
			"email":   "demo@example.com",
			"status":  "success",
			"message": "Test user details retrieved (demo mode)",
		})

	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}
