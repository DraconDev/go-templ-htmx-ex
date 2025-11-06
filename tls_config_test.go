package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	auth_cerberus "github.com/DraconDev/go-templ-htmx-ex/auth/cerberus"
)

// TestDifferentTLSConfigs tests various TLS configurations to find the working one
func TestDifferentTLSConfigs() {
	fmt.Println("🧪 Testing Different TLS Configurations")
	fmt.Println("=========================================")
	
	baseURL := "cerberus-auth-ms-548010171143.europe-west1.run.app:443"
	
	testCases := []struct {
		name        string
		dialOptions []grpc.DialOption
	}{
		{
			name: "Insecure (no TLS)",
			dialOptions: []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			},
		},
		{
			name: "Default TLS",
			dialOptions: []grpc.DialOption{
				grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
				grpc.WithBlock(),
			},
		},
		{
			name: "InsecureSkipVerify",
			dialOptions: []grpc.DialOption{
				grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true,
				})),
				grpc.WithBlock(),
			},
		},
	}
	
	for _, tc := range testCases {
		fmt.Printf("\n🔗 Testing: %s\n", tc.name)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		conn, err := grpc.NewClient(baseURL, tc.dialOptions...)
		if err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
			continue
		}
		defer conn.Close()
		
		client := cerberus.NewAuthServiceClient(conn)
		
		// Test health check
		req := &cerberus.HealthCheckRequest{
			Service: "test-client",
		}
		
		resp, err := client.HealthCheck(ctx, req)
		if err != nil {
			fmt.Printf("❌ Health check failed: %v\n", err)
		} else {
			fmt.Printf("✅ Health check successful!\n")
			fmt.Printf("   Status: %s\n", resp.GetStatus())
			fmt.Printf("   Message: %s\n", resp.GetMessage())
			
			// If this works, we found the right configuration
			if tc.name != "Insecure (no TLS)" {
				fmt.Printf("🎉 Working configuration found: %s\n", tc.name)
				return
			}
		}
	}
}

// TestWithCustomResolver tests if we need a custom name resolver
func TestWithCustomResolver() {
	fmt.Println("\n🧪 Testing with Custom DNS Resolution")
	fmt.Println("=====================================")
	
	// Try different address formats
	addresses := []string{
		"cerberus-auth-ms-548010171143.europe-west1.run.app:443",
		"cerberus-auth-ms-548010171143.europe-west1.run.app",
	}
	
	for _, addr := range addresses {
		fmt.Printf("\n🔗 Testing address: %s\n", addr)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})),
			grpc.WithBlock(),
		)
		if err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
			continue
		}
		defer conn.Close()
		
		client := cerberus.NewAuthServiceClient(conn)
		
		req := &cerberus.HealthCheckRequest{
			Service: "test-client",
		}
		
		resp, err := client.HealthCheck(ctx, req)
		if err != nil {
			fmt.Printf("❌ Health check failed: %v\n", err)
		} else {
			fmt.Printf("✅ Health check successful!\n")
			fmt.Printf("   Status: %s\n", resp.GetStatus())
			fmt.Printf("   Message: %s\n", resp.GetMessage())
		}
	}
}

func main() {
	TestDifferentTLSConfigs()
	TestWithCustomResolver()
}