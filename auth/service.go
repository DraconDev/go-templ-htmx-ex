package auth

import (
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
)

// Service handles communication with the auth microservice
type Service struct {
	config   *config.Config
	http     *HTTPClient
	builder  *RequestBuilder
	parser   *ResponseParser
	timeout  time.Duration
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	return &Service{
		config:  cfg,
		http:    NewHTTPClient(10 * time.Second),
		builder: NewRequestBuilder(cfg.AuthSecret),
		parser:  NewResponseParser(),
		timeout: 10 * time.Second,
	}
}
