package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// AuthTestRequest represents test request structure
type AuthTestRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	ProjectID   string `json:"project_id,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
}

// AuthTestResponse represents test response structure
type AuthTestResponse struct {
	UserID        string   `json:"user_id"`
	SessionToken  string   `json:"session_token"`
	Email         string   `json:"email,omitempty"`
	Valid         bool     `json:"valid,omitempty"`
	ProjectIDs    []string `json:"project_ids,omitempty"`
	Status        string   `json:"status,omitempty"`
	Error         string   `json:"error,omitempty"`
}

func main() {
	fmt.Println("Testing Authentication with Cerberus Auth Service")
	fmt.Println("=================================================")

	baseURL := "http://localhost:8080"
	authServiceURL := "https://cerberus-auth-ms-548010171143.europe-west1.run.app"

	// Test 1: Health Check
	fmt.Println("\n1. Testing Health Check...")
	testHealthCheck(baseURL, authServiceURL)

	// Test 2: Login (with dummy credentials for testing)
	fmt.Println("\n2. Testing Login...")
	testLogin(baseURL)

	// Test 3: Register
	fmt.Println("\n3. Testing Registration...")
	testRegister(baseURL)

	// Test 4: Session Validation
	fmt.Println("\n4. Testing Session Validation...")
	testSessionValidation(baseURL)
}

func testHealthCheck(baseURL, authServiceURL string) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	// Test our local health check
	resp, err := client.Get(baseURL + "/api/auth/health")
	if err != nil {
		log.Printf("Error connecting to local server: %v", err)
		fmt.Println("❌ Local server not running. Please start it first with 'go run main.go'")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Local auth health check successful")
		
		var healthResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&healthResp); err == nil {
			fmt.Printf("   Response: %+v\n", healthResp)
		}
	} else {
		fmt.Printf("❌ Health check failed with status: %d\n", resp.StatusCode)
	}
}

func testLogin(baseURL string) {
	// Create a test login request
	loginData := AuthTestRequest{
		Email:    "test@example.com",
		Password: "testpassword123",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		log.Printf("Error marshaling login request: %v", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", baseURL+"/api/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating login request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending login request: %v", err)
		fmt.Println("❌ Make sure the local server is running")
		return
	}
	defer resp.Body.Close()

	var authResp AuthTestResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		log.Printf("Error decoding login response: %v", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Login successful")
		fmt.Printf("   User ID: %s\n", authResp.UserID)
		fmt.Printf("   Session Token: %s\n", authResp.SessionToken)
		fmt.Printf("   Email: %s\n", authResp.Email)
		
		// Store token for later validation test
		// Note: In a real implementation, you'd store this securely
		fmt.Printf("   Save this JWT token for validation: %s\n", authResp.SessionToken)
	} else {
		fmt.Printf("❌ Login failed (Status: %d)\n", resp.StatusCode)
		fmt.Printf("   Error: %s\n", authResp.Error)
	}
}

func testRegister(baseURL string) {
	// Create a test registration request
	registerData := AuthTestRequest{
		Email:      "newuser@example.com",
		Password:   "newpassword123",
		ProjectID:  "test-project-123",
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		log.Printf("Error marshaling register request: %v", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", baseURL+"/api/auth/register", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating register request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending register request: %v", err)
		return
	}
	defer resp.Body.Close()

	var authResp AuthTestResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		log.Printf("Error decoding register response: %v", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Registration successful")
		fmt.Printf("   User ID: %s\n", authResp.UserID)
		fmt.Printf("   Session Token: %s\n", authResp.SessionToken)
	} else {
		fmt.Printf("❌ Registration failed (Status: %d)\n", resp.StatusCode)
		fmt.Printf("   Error: %s\n", authResp.Error)
	}
}

func testSessionValidation(baseURL string) {
	// Test with a dummy token - in real scenario, you'd use the token from login
	testToken := "dummy-jwt-token-for-testing"

	validateData := AuthTestRequest{
		SessionToken: testToken,
	}

	jsonData, err := json.Marshal(validateData)
	if err != nil {
		log.Printf("Error marshaling validation request: %v", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", baseURL+"/api/auth/validate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating validation request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending validation request: %v", err)
		return
	}
	defer resp.Body.Close()

	var authResp AuthTestResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		log.Printf("Error decoding validation response: %v", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Session validation successful")
		fmt.Printf("   Valid: %t\n", authResp.Valid)
		fmt.Printf("   User ID: %s\n", authResp.UserID)
		if len(authResp.ProjectIDs) > 0 {
			fmt.Printf("   Project IDs: %v\n", authResp.ProjectIDs)
		}
	} else {
		fmt.Printf("❌ Session validation failed (Status: %d)\n", resp.StatusCode)
		fmt.Printf("   Error: %s\n", authResp.Error)
	}
}
