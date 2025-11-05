package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AuthRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ProjectID    string `json:"project_id,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
}

type AuthResponse struct {
	UserID       string   `json:"user_id"`
	SessionToken string   `json:"session_token"`
	Email        string   `json:"email"`
	Valid        bool     `json:"valid"`
	ProjectIDs   []string `json:"project_ids"`
	Status       string   `json:"status"`
	Error        string   `json:"error"`
}

func main() {
	fmt.Println("🔐 Cerberus Auth Service Test Client")
	fmt.Println("====================================")
	
	baseURL := "http://localhost:8080"
	
	// Test 1: Health check
	fmt.Println("\n1. Testing Auth Health Check...")
	testHealthCheck(baseURL)
	
	// Test 2: Login
	fmt.Println("\n2. Testing Login...")
	testLogin(baseURL)
	
	// Test 3: Register  
	fmt.Println("\n3. Testing Registration...")
	testRegister(baseURL)
}

func testHealthCheck(baseURL string) {
	resp, err := http.Get(baseURL + "/api/auth/health")
	if err != nil {
		fmt.Printf("❌ Cannot connect to server: %v\n", err)
		fmt.Println("   Make sure to run: go run main.go")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Auth service is healthy")
		
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			fmt.Printf("   Response: %v\n", result)
		}
	} else {
		fmt.Printf("❌ Health check failed with status %d\n", resp.StatusCode)
	}
}

func testLogin(baseURL string) {
	loginData := AuthRequest{
		Email:    "test@example.com",
		Password: "testpassword",
	}
	
	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Login request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		return
	}

	if resp.StatusCode == 200 {
		fmt.Println("✅ Login successful!")
		fmt.Printf("   User ID: %s\n", result.UserID)
		fmt.Printf("   JWT Token: %s\n", result.SessionToken)
		fmt.Printf("   Email: %s\n", result.Email)
		
		// Test session validation with the token
		fmt.Println("\n   Testing session validation...")
		testValidateSession(baseURL, result.SessionToken)
	} else {
		fmt.Printf("❌ Login failed: %s\n", result.Error)
	}
}

func testRegister(baseURL string) {
	registerData := AuthRequest{
		Email:     "newuser@example.com",
		Password:  "newpassword123",
		ProjectID: "test-project",
	}
	
	jsonData, _ := json.Marshal(registerData)
	resp, err := http.Post(baseURL+"/api/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Registration request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		return
	}

	if resp.StatusCode == 200 {
		fmt.Println("✅ Registration successful!")
		fmt.Printf("   User ID: %s\n", result.UserID)
		fmt.Printf("   JWT Token: %s\n", result.SessionToken)
	} else {
		fmt.Printf("❌ Registration failed: %s\n", result.Error)
	}
}

func testValidateSession(baseURL, token string) {
	validateData := AuthRequest{
		SessionToken: token,
	}
	
	jsonData, _ := json.Marshal(validateData)
	resp, err := http.Post(baseURL+"/api/auth/validate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Validation request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		return
	}

	if resp.StatusCode == 200 && result.Valid {
		fmt.Println("✅ Session validation successful!")
		fmt.Printf("   Token is valid for user: %s\n", result.UserID)
		if len(result.ProjectIDs) > 0 {
			fmt.Printf("   Project IDs: %v\n", result.ProjectIDs)
		}
	} else {
		fmt.Printf("❌ Session validation failed: %s\n", result.Error)
	}
}
