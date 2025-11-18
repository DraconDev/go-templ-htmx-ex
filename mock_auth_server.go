package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/auth/session/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		var req struct {
			AuthCode string `json:"auth_code"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Invalid request",
			})
			return
		}
		
		fmt.Printf("Mock auth service received auth_code: %s\n", req.AuthCode)
		
		// Return the expected format with session_id and user_context
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"session_id": "web_session_123e4567-e89b-12d3-a456-426614174000",
			"user_context": map[string]interface{}{
				"user_id": "189289790288429057",
				"name":    "Dracon", 
				"email":   "dracsharp@gmail.com",
				"picture": "https://cdn.discordapp.com/avatars/189289790288429057/697bc005655b32b8bc543f9eca5899a7.png",
				"projects": map[string]interface{}{},
			},
		})
	})
	
	http.HandleFunc("/auth/session/refresh", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"user_context": map[string]interface{}{
				"user_id": "189289790288429057",
				"name":    "Dracon", 
				"email":   "dracsharp@gmail.com",
				"picture": "https://cdn.discordapp.com/avatars/189289790288429057/697bc005655b32b8bc543f9eca5899a7.png",
			},
		})
	})
	
	fmt.Println("Mock auth service starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}