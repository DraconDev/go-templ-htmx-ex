package auth

import (
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// ExchangeService handles token exchange operations
type ExchangeService struct {
	http    *HTTPClient
	builder *RequestBuilder
	parser  *ResponseParser
	baseURL string
}

// NewExchangeService creates a new exchange service
func NewExchangeService(httpClient *HTTPClient, builder *RequestBuilder, parser *ResponseParser, baseURL string) *ExchangeService {
	return &ExchangeService{
		http:    httpClient,
		builder: builder,
		parser:  parser,
		baseURL: baseURL,
	}
}

// ExchangeCodeForTokens exchanges OAuth authorization code for server session
func (es *ExchangeService) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/create", es.baseURL)
	params := map[string]string{
		"auth_code": code,
	}

	req, err := es.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to build request",
		}, err
	}

	_, bodyBytes, err := es.http.Do(req)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to call auth service: " + err.Error(),
		}, err
	}

	return es.parser.ParseTokenExchangeResponse(bodyBytes)
}