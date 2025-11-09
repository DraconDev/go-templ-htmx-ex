package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

// JWKSResponse represents the response from auth/jwks endpoint
type JWKSResponse struct {
	Keys []JWKSKey `json:"keys"`
}

type JWKSKey struct {
	Alg  string `json:"alg"`
	E    string `json:"e"`
	Kid  string `json:"kid"`
	Kty  string `json:"kty"`
	N    string `json:"n"`
	Use  string `json:"use"`
}

// UserInfo represents authenticated user information
type UserInfo struct {
	LoggedIn bool   `json:"logged_in"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Picture  string `json:"picture,omitempty"`
}

// PublicKeyCache stores verification keys
type PublicKeyCache struct {
	keys  map[string]*rsa.PublicKey
	token string
}

// NewPublicKeyCache creates a new public key cache
func NewPublicKeyCache(authServiceURL string) (*PublicKeyCache, error) {
	cache := &PublicKeyCache{
		keys: make(map[string]*rsa.PublicKey),
	}
	
	// Load initial keys
	if err := cache.loadKeys(authServiceURL); err != nil {
		return nil, err
	}
	
	// Start background refresh
	go cache.startRefresh(authServiceURL)
	
	return cache, nil
}

// loadKeys fetches public keys from JWKS endpoint
func (c *PublicKeyCache) loadKeys(authServiceURL string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(authServiceURL + "/auth/jwks")
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS endpoint returned status %d", resp.StatusCode)
	}
	
	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %v", err)
	}
	
	// Parse each key
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" && key.Alg == "RS256" {
			publicKey, err := parseRSAPublicKey(key.N, key.E)
			if err != nil {
				log.Printf("Failed to parse RSA key %s: %v", key.Kid, err)
				continue
			}
			c.keys[key.Kid] = publicKey
			log.Printf("Loaded public key: %s", key.Kid)
		}
	}
	
	if len(c.keys) == 0 {
		return fmt.Errorf("no valid RSA keys found in JWKS")
	}
	
	return nil
}

// parseRSAPublicKey parses RSA public key from modulus and exponent
func parseRSAPublicKey(modulusB64, exponentB64 string) (*rsa.PublicKey, error) {
	// Decode base64url
	modulusBytes, err := base64URLDecode(modulusB64)
	if err != nil {
		return nil, err
	}
	
	exponentBytes, err := base64URLDecode(exponentB64)
	if err != nil {
		return nil, err
	}
	
	// Convert to big integers
	modulus := new(big.Int).SetBytes(modulusBytes)
	exponent := new(big.Int).SetBytes(exponentBytes)
	
	// Create RSA public key
	publicKey := &rsa.PublicKey{
		N: modulus,
		E: int(exponent.Int64()),
	}
	
	return publicKey, nil
}

// startRefresh runs background key refresh
func (c *PublicKeyCache) startRefresh(authServiceURL string) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		if err := c.loadKeys(authServiceURL); err != nil {
			log.Printf("Failed to refresh keys: %v", err)
		}
	}
}

// ValidateJWT validates a JWT token locally
func (c *PublicKeyCache) ValidateJWT(token string) (UserInfo, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return UserInfo{LoggedIn: false}, fmt.Errorf("invalid JWT format")
	}
	
	// Decode header
	headerBytes, err := base64URLDecode(parts[0])
	if err != nil {
		return UserInfo{LoggedIn: false}, fmt.Errorf("failed to decode header: %v", err)
	}
	
	var header struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return UserInfo{LoggedIn: false}, fmt.Errorf("failed to parse header: %v", err)
	}
	
	// Verify algorithm
	if header.Alg != "RS256" {
		return UserInfo{LoggedIn: false}, fmt.Errorf("unsupported algorithm: %s", header.Alg)
	}
	
	// Get public key
	publicKey, exists := c.keys[header.Kid]
	if !exists {
		return UserInfo{LoggedIn: false}, fmt.Errorf("unknown key ID: %s", header.Kid)
	}
	
	// Verify signature
	if !verifyJWTSignature(parts[0]+"."+parts[1], parts[2], publicKey) {
		return UserInfo{LoggedIn: false}, fmt.Errorf("invalid signature")
	}
	
	// Decode payload
	payloadBytes, err := base64URLDecode(parts[1])
	if err != nil {
		return UserInfo{LoggedIn: false}, fmt.Errorf("failed to decode payload: %v", err)
	}
	
	var claims struct {
		Sub     string `json:"sub"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
		Exp     int64  `json:"exp"`
		Iss     string `json:"iss"`
	}
	
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return UserInfo{LoggedIn: false}, fmt.Errorf("failed to parse payload: %v", err)
	}
	
	// Check expiration
	if claims.Exp < time.Now().Unix() {
		return UserInfo{LoggedIn: false}, fmt.Errorf("token expired")
	}
	
	// Check issuer
	if claims.Iss != "auth-ms" {
		return UserInfo{LoggedIn: false}, fmt.Errorf("invalid issuer: %s", claims.Iss)
	}
	
	return UserInfo{
		LoggedIn: true,
		Name:     claims.Name,
		Email:    claims.Email,
		Picture:  claims.Picture,
	}, nil
}

// verifyJWTSignature verifies JWT signature using RSA
func verifyJWTSignature(data, signature string, publicKey *rsa.PublicKey) bool {
	// Decode signature
	sigBytes, err := base64URLDecode(signature)
	if err != nil {
		return false
	}
	
	// Hash the data
	hash := sha256.Sum256([]byte(data))
	
	// Verify signature
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], sigBytes)
	return err == nil
}

// base64URLDecode decodes base64url encoding
func base64URLDecode(data string) ([]byte, error) {
	// Add padding if needed
	switch len(data) % 4 {
	case 2:
		data += "=="
	case 3:
		data += "="
	case 1:
		return nil, fmt.Errorf("invalid base64url length")
	}
	
	return base64.URLEncoding.DecodeString(data)
}