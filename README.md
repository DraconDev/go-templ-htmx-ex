# Startup Platform

A modern Go-based startup application demonstrating **GitHub OAuth** and **Google OAuth** authentication integration with a dedicated auth microservice. Built with **Templ** and **HTMX** for dynamic frontend interactions.

## Features

- 🔐 **GitHub & Google OAuth 2.0 Authentication** via dedicated auth microservice
- 🏗️ **JWT Token Management** with secure session handling
- 📱 **Dynamic UI** with HTMX for seamless interactions
- 🏗️ **Server-side rendering** with Templ templates
- 🚀 **Modern Modular Architecture** with microservice patterns
- 👤 **User Profile Management** with OAuth provider account information
- 🔒 **Secure Sessions** with HttpOnly cookies
- 🎯 **Real-time Authentication Status** - UI dynamically updates to show login/logout
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
   ./bin/startup-platform
   ```

6. **Open your browser:**
   Navigate to `http://localhost:8081`

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

- `PORT` - Server port (default: 8081)
- `AUTH_SERVICE_URL` - Auth microservice URL (default: http://localhost:8080)
- `REDIRECT_URL` - Application redirect URL (default: http://localhost:8081)
- `LOG_LEVEL` - Logging level (debug, info, warn, error)

## Authentication Flow

This application demonstrates a complete OAuth authentication flow:

1. **Initiate Login**: User clicks "Login with Google/GitHub"
2. **OAuth Provider**: Redirected to Google/GitHub for authentication
3. **Token Exchange**: Auth microservice exchanges code for JWT
4. **Session Creation**: JWT stored in secure HttpOnly cookie
5. **Dynamic UI Update**: Navigation shows user info, avatar, logout button
6. **User Access**: User can access protected pages like `/profile`
7. **Logout**: Clear session via logout button

## API Endpoints

### Main Pages

- `/` - Home page with authentication features
- `/auth/google` - Initiate Google OAuth login
- `/auth/github` - Initiate GitHub OAuth login
- `/auth/callback` - Handle OAuth callback
- `/profile` - User profile page (requires authentication)
- `/health` - Application health check

### Authentication API

- `GET /api/auth/user` - Get current user info and login status
- `POST /api/auth/validate` - Validate JWT token
- `POST /api/auth/logout` - Clear session and logout
- `POST /api/auth/set-session` - Set session from OAuth callback
- `GET /api/auth/health` - Auth service health check

## Architecture

### Project Structure

```
├── main.go                    # Main application entry point
├── config/
│   └── config.go             # Configuration management
├── models/
│   └── user.go               # User and authentication models
├── auth/
│   └── service.go            # Auth service module
├── handlers/
│   └── auth.go               # Authentication HTTP handlers
└── templates/
    ├── layout.templ          # Base layout with dynamic navigation
    ├── home.templ            # Home page component
    ├── profile.templ         # User profile page
    └── auth_callback.templ   # OAuth callback processing page
```

### Key Components

- **HTTP Server**: Gorilla Mux for routing
- **Templ Components**: Type-safe HTML templates
- **HTMX Integration**: Dynamic frontend without JavaScript frameworks
- **Modular Architecture**: Separated concerns for maintainability
- **OAuth Integration**: GitHub and Google authentication providers

### Modular Design

The application follows clean architecture principles:

- **Configuration Management** (`config/`): Centralized environment configuration
- **Domain Models** (`models/`): User and authentication data structures
- **Business Logic** (`auth/`): Auth service communication and JWT handling
- **HTTP Layer** (`handlers/`): Request/response handling and routing
- **Presentation** (`templates/`): Type-safe HTML templating

## Authentication Features

### Dynamic Navigation

The application features a smart navigation system that automatically updates based on authentication status:

- **Unauthenticated Users**: See "Login with Google" and "Login with GitHub" buttons
- **Authenticated Users**: See profile picture, name, and "Logout" button
- **Real-time Updates**: Navigation changes without page reload using HTMX

### User Profile Management

After successful OAuth authentication:

- **Profile Picture**: Displays user's avatar from OAuth provider
- **User Information**: Shows name and email when available
- **Secure Session**: JWT token stored in HttpOnly cookie for security
- **Protected Routes**: Profile page requires authentication

## Testing

### Authentication Testing

1. **Start the server**: `make run`
2. **Visit**: `http://localhost:8081`
3. **Click**: "Login with Google" or "Login with GitHub"
4. **Authenticate**: Complete OAuth flow with your provider
5. **Verify**: Navigation updates to show logged-in state
6. **Test Profile**: Navigate to `/profile` to see user information
7. **Test Logout**: Click logout button to clear session

### API Testing

Test authentication endpoints:

```bash
# Check authentication status
curl http://localhost:8081/api/auth/user

# Validate session
curl -X POST http://localhost:8081/api/auth/validate

# Health check
curl http://localhost:8081/api/auth/health
```

## Production Deployment

### Building for Production

```bash
make build
```

This creates a binary in `bin/startup-platform` that can be deployed to any system with Go runtime.

### Environment Configuration

Ensure production `.env` file has:

```bash
PORT=8081
AUTH_SERVICE_URL=https://your-auth-service.com
REDIRECT_URL=https://your-app-domain.com
LOG_LEVEL=info
```

### Docker Support

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

## Troubleshooting

### Common Issues

1. **Port already in use**: Change `PORT` in `.env`
2. **Dependencies missing**: Run `make deps`
3. **Templ components not found**: Run `make generate`
4. **Authentication fails**: Check `AUTH_SERVICE_URL` includes protocol (`http://` or `https://`)
5. **OAuth callback issues**: Verify `REDIRECT_URL` matches your domain

### Debug Mode

Set `LOG_LEVEL=debug` in your `.env` file for detailed logging that shows:

- OAuth flow progression
- JWT token validation
- Auth service communication
- Session management details

### Common Authentication Issues

- **"unsupported protocol scheme"**: Ensure `AUTH_SERVICE_URL` starts with `http://` or `https://`
- **"415 Unsupported Media Type"**: Verify auth service expects JSON (not form data)
- **Session not persisting**: Check HttpOnly cookie settings and SameSite policy

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following the modular architecture
4. Add tests if applicable
5. Run `make fmt` to format code
6. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- [Templ](https://templ.guide/) for type-safe templating
- [HTMX](https://htmx.org/) for dynamic frontend interactions
- [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
- OAuth providers (Google, GitHub) for authentication services
