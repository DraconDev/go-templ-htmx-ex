package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
)

func TestCerberus() {
	fmt.Println("🧪 Starting Cerberus Auth Service Test")
	fmt.Println("=====================================")

	// Create auth client
	fmt.Println("Creating auth client...")
	authClient := auth.NewAuthClient("https://cerberus-auth-ms-548010171143.europe-west1.run.app")
	
	// Test 1: Health Check
	fmt.Println("\n1. Testing Health Check...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		fmt.Println("❌ Health check timed out")
	case <-time.After(1 * time.Second):
		resp, err := authClient.HealthCheck()
		if err != nil {
			fmt.Printf("❌ Health check failed: %v\n", err)
		} else {
			fmt.Printf("✅ Health check successful!\n")
			fmt.Printf("   User ID: %s\n", resp.UserID)
		}
	}

	// Test 2: Register
	fmt.Println("\n2. Testing Register...")
	testEmail := "test@example.com"
	testPassword := "testpassword123"
	testProjectID := "test-project-123"

	registerResp, err := authClient.Register(testEmail, testPassword, testProjectID)
	if err != nil {
		fmt.Printf("❌ Register failed: %v\n", err)
	} else {
		fmt.Printf("✅ Register successful!\n")
		fmt.Printf("   User ID: %s\n", registerResp.UserID)
		fmt.Printf("   Session Token: %s\n", registerResp.SessionToken)

		// Test 3: Login with the same credentials
		fmt.Println("\n3. Testing Login...")
		loginResp, err := authClient.Login(testEmail, testPassword)
		if err != nil {
			fmt.Printf("❌ Login failed: %v\n", err)
		} else {
			fmt.Printf("✅ Login successful!\n")
			fmt.Printf("   User ID: %s\n", loginResp.UserID)
			fmt.Printf("   Session Token: %s\n", loginResp.SessionToken)

			// Test 4: Validate Session
			fmt.Println("\n4. Testing Session Validation...")
			validatResp, err := authClient.ValidateSession(loginResp.SessionToken)
			if err != nil {
				fmt.Printf("❌ Session validation failed: %v\n", err)
			} else {
				fmt.Printf("✅ Session validation successful!\n")
				fmt.Printf("   Valid: %t\n", validatResp.Valid)
				fmt.Printf("   Project IDs: %v\n", validatResp.ProjectIDs)

				// Test 5: Get User Details
				fmt.Println("\n5. Testing Get User Details...")
				userDetailsResp, err := authClient.GetUserDetails(loginResp.UserID)
				if err != nil {
					fmt.Printf("❌ Get user details failed: %v\n", err)
				} else {
					fmt.Printf("✅ Get user details successful!\n")
					fmt.Printf("   User ID: %s\n", userDetailsResp.UserID)
					fmt.Printf("   Email: %s\n", userDetailsResp.Email)
				}
			}
		}
	}

	fmt.Println("\n🎉 Cerberus Auth Service Test Completed!")
	fmt.Println("=====================================")
}

func main() {
	TestCerberus()
}