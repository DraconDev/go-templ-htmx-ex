package main

import (
	"fmt"
	"log"
	"time"

	"star/auth"
)

func main() {
	// Create a new gRPC auth client
	fmt.Println("Creating gRPC auth client...")
	authClient, err := auth.NewGRPCAuthClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create gRPC auth client: %v", err)
	}
	defer authClient.Close()

	fmt.Println("gRPC auth client created successfully!")

	// Test health check
	fmt.Println("Testing health check...")
	timeout := time.After(5 * time.Second)
	select {
	case <-timeout:
		fmt.Println("Health check timed out - this is expected if auth server isn't running")
	case <-time.After(1 * time.Second):
		resp, err := authClient.HealthCheck()
		if err != nil {
			fmt.Printf("Health check failed: %v (expected if auth server isn't running)\n", err)
		} else {
			fmt.Printf("Health check successful: %+v\n", resp)
		}
	}

	// Test login with dummy credentials (will fail if server isn't running)
	fmt.Println("Testing login...")
	loginResp, err := authClient.Login("test@example.com", "password")
	if err != nil {
		fmt.Printf("Login failed: %v (expected if auth server isn't running)\n", err)
	} else {
		fmt.Printf("Login successful: %+v\n", loginResp)
	}

	// Test register
	fmt.Println("Testing register...")
	registerResp, err := authClient.Register("test@example.com", "password", "project-123")
	if err != nil {
		fmt.Printf("Register failed: %v (expected if auth server isn't running)\n", err)
	} else {
		fmt.Printf("Register successful: %+v\n", registerResp)
	}

	// Test validate session
	fmt.Println("Testing validate session...")
	validateResp, err := authClient.ValidateSession("dummy-token")
	if err != nil {
		fmt.Printf("Validate session failed: %v (expected if auth server isn't running)\n", err)
	} else {
		fmt.Printf("Validate session successful: %+v\n", validateResp)
	}

	// Test get user details
	fmt.Println("Testing get user details...")
	userDetailsResp, err := authClient.GetUserDetails("user-123")
	if err != nil {
		fmt.Printf("Get user details failed: %v (expected if auth server isn't running)\n", err)
	} else {
		fmt.Printf("Get user details successful: %+v\n", userDetailsResp)
	}

	fmt.Println("gRPC client test completed!")
}
