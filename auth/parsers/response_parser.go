package parsers

import (
	"encoding/json"
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// ResponseParser handles parsing responses from auth service
type ResponseParser struct{}

// NewResponseParser creates a new response parser
func NewResponseParser() *ResponseParser {
	return &ResponseParser{}
}

// ParseAuthResponse parses a response as AuthResponse
func (p *ResponseParser) ParseAuthResponse(bodyBytes []byte) (*models.AuthResponse, error) {
	var authResp models.AuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err == nil && authResp.Success {
		return &authResp, nil
	}
	return nil, fmt.Errorf("invalid auth service response format")
}

// ParseGenericResponse parses a response as generic map
func (p *ResponseParser) ParseGenericResponse(bodyBytes []byte) (map[string]interface{}, error) {
	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// ParseTokenExchangeResponse parses a response as TokenExchangeResponse
func (p *ResponseParser) ParseTokenExchangeResponse(bodyBytes []byte) (*models.TokenExchangeResponse, error) {
	var respData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &respData); err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to parse response: " + err.Error(),
		}, fmt.Errorf("failed to parse auth service response: %v", err)
	}

	// Check for errors in response
	if errMsg, hasError := respData["error"]; hasError {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   fmt.Sprintf("%v", errMsg),
		}, fmt.Errorf("auth service error: %v", errMsg)
	}

	// Extract session_id for server session management
	var sessionID string
	if sessionInterface, exists := respData["session_id"]; exists {
		if sessionStr, ok := sessionInterface.(string); ok {
			sessionID = sessionStr
		}
	}

	if sessionID == "" {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Missing session_id in auth service response",
		}, fmt.Errorf("missing session_id in auth service response")
	}

	return &models.TokenExchangeResponse{
		Success: true,
		IdToken: sessionID,
	}, nil
}