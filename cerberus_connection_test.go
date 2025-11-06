package main

import (
	"fmt"
	"testing"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
)

// TestCerberusConnection tests the actual gRPC connection to Cerberus
func TestCerberusConnection(t *testing.T) {
	fmt.Println("🧪 Testing Cerberus gRPC Connection")
	fmt.Println("=====================================")

	// Test that the client can be created
	fmt.Println("1. Creating auth client...")
	authClient := auth.NewAuthClient("https://cerberus-auth-ms-548010171143.europe-west1.run.app/")
	
	if authClient == nil {
		t.Error("Failed to create auth client")
		return
	}
	fmt.Println("✅ Auth client created successfully")
	
	// Test health check
	fmt.Println("\n2. Testing health check...")
	resp, err := authClient.HealthCheck()
	
	if err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		t.Errorf("Health check failed: %v", err)
	} else {
		fmt.Printf("✅ Health check successful!\n")
		fmt.Printf("   Status: %s\n", resp.Status)
		fmt.Printf("   Message: %s\n", resp.Message)
	}
	
	// Test register
	fmt.Println("\n3. Testing register...")
	registerResp, err := authClient.Register("test@example.com", "password123", "test-project-123")
	if err != nil {
		fmt.Printf("❌ Register failed: %v\n", err)
		t.Errorf("Register failed: %v", err)
	} else {
		fmt.Printf("✅ Register successful!\n")
		fmt.Printf("   User ID: %s\n", registerResp.UserID)
		fmt.Printf("   Session Token: %s\n", registerResp.SessionToken)
		
		// Test login with the same credentials
		fmt.Println("\n4. Testing login...")
		loginResp, err := authClient.Login("test@example.com", "password123")
		if err != nil {
			fmt.Printf("❌ Login failed: %v\n", err)
			t.Errorf("Login failed: %v", err)
		} else {
			fmt.Printf("✅ Login successful!\n")
			fmt.Printf("   User ID: %s\n", loginResp.UserID)
			fmt.Printf("   Session Token: %s\n", loginResp.SessionToken)
			
			// Test session validation
			fmt.Println("\n5. Testing session validation...")
			validatResp, err := authClient.ValidateSession(loginResp.SessionToken)
			if err != nil {
				fmt.Printf("❌ Session validation failed: %v\n", err)
				t.Errorf("Session validation failed: %v", err)
			} else {
				fmt.Printf("✅ Session validation successful!\n")
				fmt.Printf("   Valid: %t\n", validatResp.Valid)
				fmt.Printf("   Project IDs: %v\n", validatResp.ProjectIDs)
				
				// Test get user details
				fmt.Println("\n6. Testing get user details...")
				userDetailsResp, err := authClient.GetUserDetails(loginResp.UserID)
				if err != nil {
					fmt.Printf("❌ Get user details failed: %v\n", err)
					t.Errorf("Get user details failed: %v", err)
				} else {
					fmt.Printf("✅ Get user details successful!\n")
					fmt.Printf("   User ID: %s\n", userDetailsResp.UserID)
					fmt.Printf("   Email: %s\n", userDetailsResp.Email)
				}
			}
		}
	}
	
	fmt.Println("\n🎉 Cerberus gRPC Test Completed!")
	fmt.Println("=================================")
}

// TestConnectionFormats tests different connection URL formats
func TestConnectionFormats(t *testing.T) {
	fmt.Println("🧪 Testing Different Connection Formats")
	fmt.Println("========================================")
	
	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "Cloud Run with port 443",
			url:  "cerberus-auth-ms-548010171143.europe-west1.run.app:443",
		},
		{
			name: "Cloud Run with port 80",
			url:  "cerberus-auth-ms-548010171143.europe-west1.run.app:80",
		},
		{
			name: "Just domain (default port)",
			url:  "cerberus-auth-ms-548010171143.europe-west1.run.app",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n🔗 Testing: %s\n", tc.url)
			
			authClient := auth.NewAuthClient(tc.url)
			if authClient == nil {
				t.Errorf("Failed to create client for %s", tc.url)
				return
			}
			
			// Try health check
			resp, err := authClient.HealthCheck()
			if err != nil {
				fmt.Printf("❌ Health check failed: %v\n", err)
			} else {
				fmt.Printf("✅ Health check successful: %s - %s\n", resp.Status, resp.Message)
			}
		})
	}
}