# Development Commands

This document describes the available make commands for development with automatic file watching.

## Available Commands

### `make dev`
Basic development mode that watches for `.go` and `.templ` file changes and automatically rebuilds and runs the server.

**Features:**
- Watches both Go source files and templ templates
- Automatically regenerates templ components
- Rebuilds and restarts the server on changes
- Uses inotifywait for efficient file monitoring

### `make watch`
Comprehensive watch mode with enhanced process management.

**Features:**
- Watches both `.go` and `.templ` files
- Process management to avoid multiple instances
- Graceful handling of server lifecycle
- Fallback support for systems without inotifywait
- Better error handling and user feedback

### `make dev-watch`
Advanced development mode with automatic process cleanup.

**Features:**
- Kills any existing server processes before starting
- Ensures clean restart on every file change
- Best for development when you want to ensure no stale processes

## Requirements

### Linux/Debian/Ubuntu
```bash
sudo apt-get install inotify-tools
```

### macOS
```bash
brew install inotify-tools
```

### Fallback
If `inotifywait` is not available, the commands will fall back to simple file polling using `find`, though this is less efficient.

## Usage Examples

Start development mode:
```bash
make dev
```

Start comprehensive watch mode:
```bash
make watch
```

Start advanced development mode:
```bash
make dev-watch
```

## How It Works

1. **File Watching**: Uses `inotifywait` to monitor file changes in real-time
2. **Change Detection**: Watches for modifications to `.go` and `.templ` files
3. **Auto-build**: Automatically runs `make generate` and `make build`
4. **Auto-restart**: Starts the server after successful build
5. **Process Management**: Handles existing processes to avoid conflicts

## Manual Development Workflow

If you prefer manual control:

1. Edit your `.go` or `.templ` files
2. Run `make generate` to regenerate templ components
3. Run `make build` to compile the Go code
4. Run `./bin/microservice-test` to start the server

The make commands automate steps 1-4 for a seamless development experience.
