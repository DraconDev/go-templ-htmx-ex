# Makefile for Go templ/HTMX authentication platform

# Build targets
build:
	go build -v

test:
	go test -v ./...

lint: ## Run golangci-lint
	golangci-lint run --timeout=2m

fmt: ## Format code
	gofmt -w .

check: fmt lint build ## Run all checks (fmt + lint + build)

# Development targets
dev:
	air -c .air.toml

# Clean targets
clean:
	rm -f go-templ-htmx-ex
	find . -name "*.go" -exec gofmt -l {} \;

# Database targets
db-init:
	@echo "Database initialization required - configure PostgreSQL and run migrations"

.PHONY: build test lint fmt check dev clean db-init
