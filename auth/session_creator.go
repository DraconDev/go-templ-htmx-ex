package auth

import (
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// SessionCreator handles session creation operations
type SessionCreator struct {
	http    *HTTPClient
	builder *RequestBuilder
	parser  *ResponseParser
	baseURL string
}

// NewSessionCreator creates a new session creator
func NewSessionCreator(httpClient *HTTPClient, builder *RequestBuilder, parser *ResponseParser, baseURL string) *SessionCreator {
	return &SessionCreator{
		http:    httpClient,
		builder: builder,
		parser:  parser,
		baseURL: baseURL,
	}
}

// CreateSession creates a session from authorization code
func (sc *SessionCreator) CreateSession(code string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/auth/session/create", sc.baseURL)
	params := map[string]string{"code": code}

	req, err := sc.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	_, bodyBytes, err := sc.http.Do(req)
	if err != nil {
		return nil, err
	}

	return sc.parser.ParseGenericResponse(bodyBytes)
}

// ExchangeCodeForTokens exchanges OAuth authorization code for server session
func (sc *SessionCreator) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/create", sc.baseURL)
	params := map[string]string{
		"auth_code": code,
	}

	req, err := sc.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to build request",
		}, err
	}

	_, bodyBytes, err := sc.http.Do(req)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to call auth service: " + err.Error(),
		}, err
	}

	return sc.parser.ParseTokenExchangeResponse(bodyBytes)
}