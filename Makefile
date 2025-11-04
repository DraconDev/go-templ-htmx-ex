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
	@pkill -f "$(BUILD_DIR)/$(BINARY_NAME)" 2>/dev/null || true; \
	echo "Killed any existing server processes"; \
	while true; do \
		inotifywait -e modify -r . --include '\.(go|templ)$$' 2>/dev/null || break; \
		echo "Changes detected, rebuilding..."; \
		pkill -f "$(BUILD_DIR)/$(BINARY_NAME)" 2>/dev/null || true; \
		$(MAKE) generate; \
		$(MAKE) build; \
		./$(BUILD_DIR)/$(BINARY_NAME); \
	done

watch: generate
	@echo "Starting comprehensive watch mode (Go + Templ files)..."
	@if command -v inotifywait >/dev/null 2>&1; then \
		while true; do \
			inotifywait -e modify -r . --include '\.(go|templ)$$' 2>/dev/null || break; \
			echo "File changes detected, rebuilding and restarting..."; \
			pkill -f "$(BUILD_DIR)/$(BINARY_NAME)" 2>/dev/null || true; \
			$(MAKE) generate && \
			$(MAKE) build && \
			./$(BUILD_DIR)/$(BINARY_NAME) & \
			echo "Server restarted with PID $$!"; \
		done; \
	else \
		echo "inotifywait not found. Install with: sudo apt-get install inotify-tools (Debian/Ubuntu)"; \
		echo "Falling back to simple file watching with find..."; \
		while true; do \
			pkill -f "$(BUILD_DIR)/$(BINARY_NAME)" 2>/dev/null || true; \
			find . -name "*.go" -o -name "*.templ" | while read file; do \
				if [ "$$file" -nt /tmp/last_build ]; then \
					touch /tmp/last_build; \
					echo "Changes detected, rebuilding and restarting..."; \
					$(MAKE) generate && \
					$(MAKE) build && \
					./$(BUILD_DIR)/$(BINARY_NAME) & \
					break; \
				fi; \
			done; \
			sleep 5; \
		done; \
	fi

air: generate
	@echo "Starting development server with Air live reload..."
	air

air-watch: generate
	@echo "Starting development server with Air live reload (alternative name)..."
	air

dev-watch: generate
	@echo "Starting development server with comprehensive file watching..."
	@if command -v inotifywait >/dev/null 2>&1; then \
		while true; do \
			inotifywait -e modify -r . --include '\.(go|templ)$$' 2>/dev/null || break; \
			echo "File changes detected, rebuilding and restarting..."; \
			pkill -f "$(BUILD_DIR)/$(BINARY_NAME)" 2>/dev/null || true; \
			$(MAKE) generate && \
			$(MAKE) build && \
			./$(BUILD_DIR)/$(BINARY_NAME) & \
			echo "Server restarted with PID $$!"; \
		done; \
	else \
		echo "inotifywait not found. Please install inotify-tools for optimal file watching."; \
		echo "Install with: sudo apt-get install inotify-tools (Debian/Ubuntu) or brew install inotify-tools (macOS)"; \
		echo "Running once without live reload..."; \
		./$(BUILD_DIR)/$(BINARY_NAME); \
	fi

run: generate build
	@echo "Starting $(BINARY_NAME)..."
	@pkill -f "$(BUILD_DIR)/$(BINARY_NAME)" 2>/dev/null || true; \
	echo "Killed any existing server processes"; \
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	$(GOTEST) ./...

fmt:
	@echo "Formatting Go code..."
	$(GOFMT) ./...

all: deps generate build
	@echo "Setup complete!"

.PHONY: build clean deps generate dev watch air air-watch dev-watch run test fmt all
