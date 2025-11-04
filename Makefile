# Makefile for Microservice Test Harness

BINARY_NAME=microservice-test
BUILD_DIR=bin

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

build: generate clean
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Built $(BINARY_NAME) successfully!"

clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Cleaned successfully!"

deps:
	@echo "Installing dependencies..."
	$(GOGET) -u github.com/a-h/templ/cmd/templ@latest
	$(GOMOD) download

generate:
	@echo "Generating templ components..."
	templ generate -path .
	@echo "Templ components generated!"

dev: generate
	@echo "Starting development server with hot reload..."
	while true; do \
		inotifywait -e modify -r . --include '\.templ$$' 2>/dev/null || break; \
		echo "Changes detected, rebuilding..."; \
		$(MAKE) generate; \
		$(MAKE) build; \
	done

air: generate
	@echo "Starting Air development server..."
	go run github.com/air-verse/air@latest

live: generate
	@echo "Starting Air live reload server..."
	go run github.com/air-verse/air@latest

run: build
	@echo "Starting $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	$(GOTEST) ./...

fmt:
	@echo "Formatting Go code..."
	$(GOFMT) ./...

all: deps generate build
	@echo "Setup complete!"

.PHONY: build clean deps generate dev air live run test fmt all
