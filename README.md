# Startup Platform

A modern Go-based startup application demonstrating Google OAuth authentication integration with a dedicated auth microservice. Built with **Templ** and **HTMX** for dynamic frontend interactions.

## Features

- 🔐 **Google OAuth 2.0 Authentication** via dedicated auth microservice
- 🏗️ **JWT Token Management** with secure session handling
- 📱 **Dynamic UI** with HTMX for seamless interactions
- 🏗️ **Server-side rendering** with Templ templates
- 🚀 **Modern Architecture** with microservice patterns
- 👤 **User Profile Management** with Google account information
- 🔒 **Secure Sessions** with HttpOnly cookies
- 🚀 **Fast Development** with hot reload capabilities

## Technology Stack

- **Go 1.21+** - Main programming language
- **Templ** - Type-safe HTML templating
- **HTMX** - Dynamic frontend interactions
- **Gorilla Mux** - HTTP routing
- **Tailwind CSS** - Styling (via CDN)
- **Godotenv** - Environment configuration

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Templ CLI tool (`go install github.com/a-h/templ/cmd/templ@latest`)

### Installation

1. **Clone and setup the project:**
   ```bash
   git clone <your-repo-url>
   cd go-templ-htmx-ex
   ```

2. **Install dependencies:**
   ```bash
   make deps
   ```

3. **Generate templ components:**
   ```bash
   make generate
   ```

4. **Build the application:**
   ```bash
   make build
   ```

5. **Run the application:**
   ```bash
   ./bin/microservice-test
   ```

6. **Open your browser:**
   Navigate to `http://localhost:8080`

## Development

### Hot Reload Development

For development with automatic rebuilding:

```bash
make dev
```

This will watch for changes and automatically rebuild the application.

### Available Make Commands

```bash
make help        # Show all available commands
make deps        # Install dependencies
make generate    # Generate templ components
make build       # Build the application
make run         # Build and run the application
make dev         # Development mode with hot reload
make clean       # Clean build artifacts
make test        # Run tests
make fmt         # Format Go code
```

## Configuration

### Environment Variables

Copy `.env.example` to `.env` and customize:

```bash
cp .env.example .env
```

Key configuration options:

- `PORT` - Server port (default: 8080)
- `SERVICE_TIMEOUT` - Timeout for service calls
- `LOG_LEVEL` - Logging level (debug, info, warn, error)
- Service URLs - Configure your microservice endpoints

## API Endpoints

### Main Pages

- `/` - Home page with overview
- `/test` - Main testing dashboard
- `/test/{service}` - Service-specific testing page
- `/health` - Application health check

### API Endpoints

- `GET /api/services` - Discover available services
- `POST /api/test` - Run a test on a service

## Usage Examples

### Basic Service Testing

1. Navigate to `/test`
2. Click "Discover Services" to auto-detect services
3. Use "Manual Service Test" to test specific endpoints
4. View results in real-time

### HTMX Features

The application uses HTMX for dynamic interactions:

- **Service Discovery**: Click to load services dynamically
- **Form Submissions**: Tests run without page reloads
- **Real-time Updates**: Results appear instantly

### Service URL Examples

For testing your microservices, use URLs like:

- Health Check: `http://your-service:8001/health`
- API Endpoint: `http://your-service:8001/api/users`
- Custom Endpoint: `http://your-service:8001/api/orders`

## Architecture

### Project Structure

```
├── main.go                    # Main application
├── components/                # Templ templates
│   ├── layout.templ          # Base layout
│   ├── home.templ            # Home page component
│   ├── microservice_test.templ # Testing dashboard
│   └── service_test.templ    # Service-specific testing
├── Makefile                  # Build automation
├── .env.example              # Environment configuration
└── README.md                 # This file
```

### Key Components

- **HTTP Server**: Gorilla Mux for routing
- **Templ Components**: Type-safe HTML templates
- **HTMX Integration**: Dynamic frontend without JavaScript frameworks
- **Service Testing**: HTTP client for testing microservice endpoints

## Testing Microservices

### Health Checks

Test service health endpoints:

```
GET /health
```

### API Testing

Test various HTTP methods:

- **GET** - Retrieve data
- **POST** - Create resources
- **PUT** - Update resources
- **DELETE** - Remove resources

### Stress Testing

Run multiple concurrent tests to check service performance and resilience.

## Production Deployment

### Building for Production

```bash
make build
```

### Docker Support

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

### System Installation

```bash
make install  # Installs to /usr/local/bin
```

## Troubleshooting

### Common Issues

1. **Port already in use**: Change `PORT` in `.env`
2. **Dependencies missing**: Run `make deps`
3. **Templ components not found**: Run `make generate`
4. **Services not responding**: Check service URLs and network connectivity

### Debug Mode

Set `LOG_LEVEL=debug` in your `.env` file for detailed logging.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make fmt` to format code
6. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- [Templ](https://templ.guide/) for type-safe templating
- [HTMX](https://htmx.org/) for dynamic frontend interactions
- [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
