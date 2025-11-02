# Makefile for Microservice Test Harness

# Variables
BINARY_NAME=microservice-test
BUILD_DIR=bin
SOURCE_DIR=.
TEMPLATES_DIR=.

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build the project
build: clean
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Built $(BINARY_NAME) successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Cleaned successfully!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) -u github.com/a-h/templ/cmd/templ@latest
	$(GOMOD) download

# Generate templ components
generate:
	@echo "Generating templ components..."
	templ generate
	@echo "Templ components generated!"

# Watch for changes and rebuild automatically
dev: generate
	@echo "Starting development server with hot reload..."
	while true; do \
		inotifywait -e modify -r $(SOURCE_DIR) --include '\.go$$|\.templ$$' 2>/dev/null || break; \
		echo "Changes detected, rebuilding..."; \
		$(MAKE) build; \
	done

# Run the application
run: build
	@echo "Starting $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with custom port
run-port: build
	@echo "Starting $(BINARY_NAME) on custom port..."
	@read -p "Enter port (default 8080): " PORT; \
	[ -z "$$PORT" ] && PORT=8080; \
	./$(BUILD_DIR)/$(BINARY_NAME)

# Test the application
test:
	@echo "Running tests..."
	$(GOTEST) ./...

# Install the application
install: deps generate build
	@echo "Installing $(BINARY_NAME)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Format Go code
fmt:
	@echo "Formatting Go code..."
	$(GOFMT) ./...

# Lint the code (requires golangci-lint to be installed)
lint:
	@echo "Linting code..."
	golangci-lint run

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t microservice-test-harness .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 microservice-test-harness

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  generate     - Generate templ components"
	@echo "  dev          - Development mode with hot reload"
	@echo "  run          - Run the application"
	@echo "  run-port     - Run with custom port"
	@echo "  test         - Run tests"
	@echo "  install      - Install to system"
	@echo "  fmt          - Format Go code"
	@echo "  lint         - Lint code"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help"

# Default target
.PHONY: all
all: deps generate build
	@echo "Setup complete!"

# .PHONY targets that don't create files
.PHONY: build clean deps generate dev run run-port test install fmt lint docker-build docker-run help all
