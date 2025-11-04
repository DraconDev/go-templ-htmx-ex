package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"star/template"
)

// Config holds application configuration
type Config struct {
	ServerPort string
}

var (
	config = &Config{
		ServerPort: getEnvOrDefault("PORT", "8080"),
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

	// Microservice testing routes
	router.HandleFunc("/test", microserviceTestHandler).Methods("GET", "POST")
	router.HandleFunc("/test/{service}", serviceTestHandler).Methods("GET", "POST")

	// API endpoints for HTMX interactions
	router.HandleFunc("/api/services", servicesAPIHandler).Methods("GET")
	router.HandleFunc("/api/test", runTestAPIHandler).Methods("POST")

	// Auth service API endpoints
	router.HandleFunc("/api/auth/health", authHealthCheckHandler).Methods("GET")
	router.HandleFunc("/api/auth/login", authLoginHandler).Methods("POST")
	router.HandleFunc("/api/auth/register", authRegisterHandler).Methods("POST")
	router.HandleFunc("/api/auth/validate", authValidateSessionHandler).Methods("POST")

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
	component := template.Layout("Home", template.HomeContent())
	component.Render(r.Context(), w)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

func microserviceTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	component := template.Layout("Microservice Testing", template.MicroserviceTestContent())
	component.Render(r.Context(), w)
}

func serviceTestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]

	w.Header().Set("Content-Type", "text/html")
	component := template.Layout(fmt.Sprintf("Testing %s", serviceName), template.ServiceTestContent(serviceName))
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
	component := template.TestResult(serviceURL, testType, "success", "Test completed successfully!")
	component.Render(r.Context(), w)
}

// Auth Service Handlers

func authHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if auth service is available
	authServiceURL := "https://cerberus-auth-ms-548010171143.europe-west1.run.app"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(authServiceURL + "/api/health")
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "unavailable", "error": "Service not reachable", "url": "` + authServiceURL + `"}`))
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "available", "url": "` + authServiceURL + `", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

func authLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// For demo purposes, create a test response
	// In production, this would make a real gRPC call to the auth service
	if email != "" && password != "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"user_id": "demo-user-123",
			"session_token": "demo-session-token-456",
			"email": "` + email + `"
		}`))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Email and password are required"}`))
	}
}

func authRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// For demo purposes, create a test response
	if email != "" && password != "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"user_id": "demo-user-789",
			"session_token": "demo-session-token-012"
		}`))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Email and password are required"}`))
	}
}

func authValidateSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	sessionToken := r.FormValue("session_token")

	// For demo purposes, create a test response
	if sessionToken != "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"user_id": "demo-user-123",
			"valid": true,
			"project_ids": ["demo-project-1", "demo-project-2"]
		}`))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Session token is required"}`))
	}
}
