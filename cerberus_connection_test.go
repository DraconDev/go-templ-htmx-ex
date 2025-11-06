package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
)

func main() {
	// Try different connection URLs based on environment
	connectionURLs := []string{
		os.Getenv("AUTH_GRPC_URL"), // Environment variable first
		"localhost:50051",          // Local development
		"cerberus-auth-ms-548010171143.europe-west1.run.app:443", // Cloud Run
	}

	fmt.Println("🔍 Testing Cerberus Connection Options:")
	fmt.Println("Environment Variables:")
	fmt.Printf("  AUTH_GRPC_URL: %s\n", os.Getenv("AUTH_GRPC_URL"))
	fmt.Printf("  AUTH_GRPC_TIMEOUT: %s\n", os.Getenv("AUTH_GRPC_TIMEOUT"))
	fmt.Println()

	for i, url := range connectionURLs {
		if url == "" {
			fmt.Printf("Skipping empty URL %d\n", i+1)
			continue
		}

		fmt.Printf("🔗 Testing connection to: %s\n", url)

		// Set a timeout for the connection attempt
		timeout := time.After(5 * time.Second)
		done := make(chan bool)

		go func(connectionURL string) {
			authClient, err := auth.NewGRPCAuthClient(connectionURL)
			if err != nil {
				fmt.Printf("❌ Failed to create auth client: %v\n", err)
				done <- false
				return
			}
			defer authClient.Close()

			fmt.Printf("✅ gRPC client created successfully for %s!\n", connectionURL)

			// Test health check with a short timeout
			healthTimeout := time.After(3 * time.Second)
			select {
			case <-healthTimeout:
				fmt.Printf("⏰ Health check timed out for %s\n", connectionURL)
			default:
				resp, err := authClient.HealthCheck()
				if err != nil {
					fmt.Printf("⚠️  Health check failed: %v\n", err)
					// This is often expected if the service is running but not responding
				} else {
					fmt.Printf("🎉 Health check successful: %+v\n", resp)
				}
			}

			done <- true
		}(url)

		select {
		case success := <-done:
			if success {
				fmt.Printf("✅ Connection test completed for %s\n", url)
				break // Found a working connection
			}
		case <-timeout:
			fmt.Printf("⏰ Connection test timed out for %s\n", url)
		}

		fmt.Println()
	}

	fmt.Println("🏁 Connection test completed!")

	// Try to use the working connection if we found one
	if workingURL := os.Getenv("AUTH_GRPC_URL"); workingURL != "" {
		fmt.Printf("\n🔄 Testing full functionality with %s...\n", workingURL)
		testFullFunctionality(workingURL)
	}
}

func testFullFunctionality(baseURL string) {
	authClient, err := auth.NewGRPCAuthClient(baseURL)
	if err != nil {
		log.Printf("Failed to create auth client: %v", err)
		return
	}
	defer authClient.Close()

	// Test all methods with timeout
	timeout := time.After(10 * time.Second)

	// Health Check
	fmt.Println("🏥 Testing Health Check...")
	select {
	case <-timeout:
		fmt.Println("⏰ Health check timed out")
	default:
		resp, err := authClient.HealthCheck()
		if err != nil {
			fmt.Printf("⚠️  Health check: %v\n", err)
		} else {
			fmt.Printf("✅ Health check: %+v\n", resp)
		}
	}

	// Test Login (with dummy data)
	fmt.Println("🔐 Testing Login...")
	select {
	case <-timeout:
		fmt.Println("⏰ Login test timed out")
	default:
		resp, err := authClient.Login("test@example.com", "password")
		if err != nil {
			fmt.Printf("⚠️  Login: %v\n", err)
		} else {
			fmt.Printf("✅ Login: %+v\n", resp)
		}
	}

	// Test Register
	fmt.Println("📝 Testing Register...")
	select {
	case <-timeout:
		fmt.Println("⏰ Register test timed out")
	default:
		resp, err := authClient.Register("test@example.com", "password", "test-project")
		if err != nil {
			fmt.Printf("⚠️  Register: %v\n", err)
		} else {
			fmt.Printf("✅ Register: %+v\n", resp)
		}
	}

	fmt.Println("🎯 Full functionality test completed!")
}