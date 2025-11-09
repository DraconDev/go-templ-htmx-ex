package main

import (
	"fmt"
	"github.com/DraconDev/go-templ-htmx-ex/auth"
)

func main() {
	authService := auth.NewService("http://localhost:8080")
	
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMTgxNTgzNjE0ODQ2OTE5NDE0MDUiLCJuYW1lIjoiRHJhY29uIiwiZW1haWwiOiJkcmFjc2hhcnBAZ21haWwuY29tIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FDZzhvY0lmTEVYUFVnbHZTbU1ZQllNelQtaUlnYWQ3X1hNTXJqUEVVZ1FPazVQOGlFbnVyS3ByPXM5Ni1jIiwiaXNzIjoiYXV0aC1tcyIsImlhdCI6MTc2MjY0NTM1NywiZXhwIjoxNzYyNzMxNzU3LCJ0b2tlbl90eXBlIjoiYWNjZXNzIn0.N2cqxeJB0Xz5rGF4MeJM38SmXptPSTwaTSAYbaQIq4I"
	
	fmt.Println("Testing auth service...")
	resp, err := authService.ValidateUser(token)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Response: %+v\n", resp)
	fmt.Printf("Success: %v\n", resp.Success)
	if resp.Success {
		fmt.Printf("Name: %s, Email: %s\n", resp.Name, resp.Email)
	}
}