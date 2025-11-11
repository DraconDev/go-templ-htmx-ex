# 🚀 Go + HTMX + Templ Startup Platform

A **production-ready startup platform** with **Google OAuth authentication**, **PostgreSQL database**, **admin dashboard**, and **modern UI**. Built with **Templ**, **HTMX**, and **SQLC** for high performance and scalability.

![Platform Preview](https://via.placeholder.com/800x400/1e293b/3b82f6?text=Modern+Startup+Platform+with+Database)

## 🎯 Perfect For

- **🚀 Modern startup projects** requiring user authentication and database
- **🏗️ Microservice architectures** with dedicated auth services
- **📊 Admin dashboards** with real-time user analytics
- **⚡ HTMX + Templ** applications without JavaScript frameworks
- **📈 Rapid MVP development** with enterprise-grade architecture
- **💰 E-commerce platforms** ready for payment integration

## ✨ Implemented Features

### 🔐 **Production Authentication System**
- **🔐 Google OAuth 2.0** via dedicated auth microservice
- **⚡ JWT Token Management** with local validation for performance
- **🎯 Real User Data** - shows actual names, emails, profile pictures
- **🔒 Secure Session Handling** with HttpOnly cookies
- **📱 Dynamic Navigation** - Login/Logout button changes based on auth state
- **🧪 Comprehensive Test Coverage** for authentication flows

### 💾 **Database Integration (SQLC)**
- **📊 PostgreSQL Database** with professional schema
- **🧬 SQLC Generated Queries** - type-safe database operations
- **📈 Real Data Display** - no mock data, shows actual user metrics
- **🔍 User Analytics** - total users, signups today, recent users
- **🛠️ Database Migrations** - versioned schema management
- **⚡ Connection Pooling** - optimized for production performance

### 📊 **Admin Dashboard**
- **📈 User Statistics** - total users, today's signups, weekly growth
- **👥 User Management** - view recent users with real data
- **🔒 Admin Protection** - only admin users can access
- **📱 Responsive Design** - works on all devices
- **⚡ Real-time Data** - displays actual database metrics

### 🎨 **Modern UI & UX**
- **🌙 Professional Dark Theme** with clean design
- **🎭 Interactive Elements** with hover effects and transitions
- **📱 Responsive Design** that works perfectly across all devices
- **⚡ HTMX Integration** - seamless user interactions without JavaScript frameworks
- **🚀 Glass Morphism Effects** - modern visual design

### 🏗️ **Enterprise Architecture**
- **📦 Docker Containerization** - ready for production deployment
- **🧪 Test Coverage** - comprehensive testing for core functionality
- **📊 Health Monitoring** - application and database health checks
- **🔧 Environment Configuration** - 12-factor app design
- **🛠️ Microservice Ready** - designed for distributed systems

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL database
- Templ CLI: `go install github.com/a-h/templ/cmd/templ@latest`

### One-Command Setup

```bash
# Clone and setup
git clone <your-repo>
cd go-templ-htmx-ex

# Install dependencies
make deps

# Generate templates
make generate

# Setup database
make db-migrate

# Run development server
make dev
```

**Visit:** `http://localhost:8081`

## 🔧 Configuration

### Environment Setup

```bash
# Copy example config
cp .env.example .env

# Edit for your setup
# PORT=8081
# AUTH_SERVICE_URL=http://localhost:8080  # Your auth service
# DATABASE_URL=postgresql://user:pass@localhost:5432/dbname
# ADMIN_EMAIL=admin@yourdomain.com        # Admin user email
```

### Database Setup

```bash
# Create PostgreSQL database
createdb startup_platform

# Run migrations
make db-migrate

# Optional: Seed with sample data
make db-seed
```

## 📊 Current Implementation Status

### ✅ **Completed Features (28/30 tasks)**

- **✅ PostgreSQL Database Integration** - Full database with users table
- **✅ SQLC Query Generation** - Type-safe database operations  
- **✅ Google OAuth Authentication** - Production-ready auth flow
- **✅ JWT Token Management** - Secure session handling
- **✅ Admin Dashboard** - Real user analytics and management
- **✅ User Profile System** - Complete user profile pages
- **✅ Docker Containerization** - Production deployment ready
- **✅ Test Coverage** - Comprehensive authentication testing
- **✅ Real Data Display** - No mock data, actual database metrics
- **✅ Modern UI** - Professional dark theme with HTMX
- **✅ Health Monitoring** - Application and database health checks

### 🔄 **Remaining Tasks**

- **📊 Health Checks & Monitoring** - Enhanced monitoring endpoints
- **🚀 Production Optimization** - Final production deployment optimizations

## 🏗️ Architecture

### Project Structure

```
├── main.go                    # Application entry point
├── Dockerfile                 # Production containerization
├── sqlc.yaml                 # SQLC configuration
├── config/
│   └── config.go            # Environment configuration
├── db/
│   ├── migrations/          # Database schema migrations
│   ├── sqlc/               # Generated type-safe queries
│   └── init.go             # Database connection setup
├── handlers/
│   ├── admin.go            # Admin dashboard handlers
│   ├── auth.go             # Authentication handlers
│   ├── handlers.go         # User and business logic
│   └── admin_test.go       # Authentication test suite
├── middleware/
│   └── auth.go             # JWT validation middleware
├── auth/
│   └── service.go          # Auth service integration
└── templates/
    ├── layout.templ        # Base layout with navigation
    ├── home.templ          # Landing page
    ├── profile.templ       # User profile page
    ├── admin_dashboard.templ # Admin analytics dashboard
    └── auth_callback.templ # OAuth callback processing
```

### Database Schema

```sql
-- Users table with admin support
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    auth_id TEXT UNIQUE NOT NULL,
    picture TEXT,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### API Endpoints

#### Main Pages
- `GET /` - Home page with authentication features
- `GET /profile` - User profile page (auth required)
- `GET /admin` - Admin dashboard (admin required)
- `GET /health` - Application health check

#### Authentication API
- `GET /auth/google` - Google OAuth login initiation
- `GET /auth/callback` - OAuth callback processing
- `GET /api/auth/user` - Get current user info
- `POST /api/auth/validate` - Validate JWT token
- `POST /api/auth/logout` - Logout user
- `GET /api/auth/health` - Auth system health

#### Admin API
- Admin dashboard loads real-time user statistics from database
- Shows total users, signups today, this week's growth
- Displays recent users with actual profile data

## 🧪 Testing

### Test Coverage

```bash
# Run all tests
make test

# Run specific test suite
go test ./handlers/ -v

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Results

```
=== RUN   TestAdminDashboardAccess
--- PASS: TestAdminDashboardAccess (0.00s)
=== RUN   TestAdminDashboardUnauthorized  
--- PASS: TestAdminDashboardUnauthorized (0.00s)
=== RUN   TestAdminDashboardNoAuth
--- PASS: TestAdminDashboardNoAuth (0.00s)
=== RUN   TestAdminDashboardMiddlewareIntegration
--- PASS: TestAdminDashboardMiddlewareIntegration (0.00s)
PASS
ok      github.com/DraconDev/go-templ-htmx-ex/handlers    0.002s
```

### Manual Testing

1. **Start services:**
   ```bash
   make dev  # Your app
   # Ensure auth service is running on port 8080
   # Ensure PostgreSQL is running
   ```

2. **Test authentication flow:**
   - Visit `http://localhost:8081`
   - Click "Login with Google"
   - Complete authentication
   - Verify user profile displays real data

3. **Test admin features:**
   - Login with admin email (configured in .env)
   - Navigate to `/admin`
   - Verify dashboard shows real user statistics

4. **Test logout:**
   - Click logout button
   - Verify navigation changes to "Login with Google"

## 📈 Performance Features

### Optimized Database Operations
- **SQLC Generated Queries** - Type-safe and optimized
- **Connection Pooling** - Efficient database connections
- **Real-time Analytics** - Instant dashboard data loading

### Authentication Performance
- **Local JWT Validation** - Fast authentication checks
- **Cookie-based Sessions** - Efficient session management
- **Smart Caching** - Reduced database queries

### UI Performance
- **HTMX Integration** - Seamless updates without full page reloads
- **Optimized Templates** - Fast rendering with Templ
- **Responsive Design** - Consistent performance across devices

## 🐳 Docker Deployment

### Development

```bash
# Build development image
make docker-build-dev

# Run with docker-compose
docker-compose up -d
```

### Production

```bash
# Build production image
make docker-build

# Run production container
docker run -p 8081:8081 \
  -e AUTH_SERVICE_URL=https://auth.your-domain.com \
  -e DATABASE_URL=postgresql://user:pass@db:5432/dbname \
  your-app:latest
```

## 📚 Development Commands

```bash
# Development
make dev              # Hot reload development server
make deps             # Install all dependencies
make generate         # Generate Templ components
make generate-db      # Regenerate SQLC queries

# Database
make db-migrate       # Run database migrations
make db-seed         # Seed with sample data
make db-reset        # Reset database (development only)

# Building & Testing
make build           # Production build
make test            # Run all tests
make lint            # Run code linter
make fmt             # Format code

# Docker
make docker-build    # Build production image
make docker-run      # Run container
make docker-logs     # View container logs
```

## 🚀 For Startup Projects

### What You Get Out of the Box

1. **✅ Production Authentication** - Google OAuth with JWT
2. **✅ Database Integration** - PostgreSQL with type-safe queries
3. **✅ Admin Dashboard** - Real user analytics and management
4. **✅ Test Coverage** - Comprehensive testing for core features
5. **✅ Docker Ready** - Production deployment with containers
6. **✅ Modern UI** - Professional design with HTMX + Templ

### Ready for Business Logic

```bash
# Add your specific business features
mkdir handlers/business
vim handlers/business/your_feature.go

# Add database tables for your needs
vim db/migrations/002_your_feature.sql

# Create templates for your pages
vim templates/your_feature.templ
```

### Monetization Ready

The platform is designed for immediate revenue generation:
- **💰 Payment Integration** - Connect to payment microservice
- **📊 User Analytics** - Track growth and engagement
- **👥 Admin Tools** - Manage users and content
- **🔒 Security** - Enterprise-grade authentication

## 📈 Analytics & Monitoring

### Database Metrics
- **Total Users** - Complete user count from database
- **Signups Today** - Daily growth tracking
- **Weekly Growth** - Week-over-week user acquisition
- **Recent Users** - Latest user registrations with profiles

### Application Health
- **Database Connectivity** - Real-time connection status
- **Authentication Status** - Auth service health monitoring
- **Response Times** - Performance metrics for optimization

## 🔒 Security Features

### Production-Ready Security
- **🔐 JWT Signature Verification** - Local validation with public keys
- **⏰ Token Expiration** - Automatic invalidation of expired tokens
- **🏷️ Issuer Validation** - Ensures tokens from correct auth service
- **🍪 Secure Cookies** - HttpOnly, SameSite, Secure flags
- **🛡️ Admin Protection** - Role-based access control

### Database Security
- **🔒 Parameterized Queries** - SQL injection protection via SQLC
- **🛡️ Connection Security** - Encrypted database connections
- **👤 Admin Role Management** - Database-level admin controls

## 🎯 Why This Platform

### For Developers
- ⚡ **Fast Development** - 5-minute setup with modern tooling
- 🧠 **Easy to Understand** - Clean, well-documented Go code
- 🔧 **Easy to Customize** - Modular architecture with SQLC
- 📈 **Easy to Scale** - Proven patterns and enterprise architecture
- 🧪 **Well Tested** - Comprehensive test coverage for reliability

### For Startups
- 🚀 **Quick to Market** - Production-ready auth and database
- 💰 **Cost Effective** - Optimized database queries and caching
- 🛡️ **Enterprise Security** - JWT validation, secure cookies, signature verification
- 📊 **Business Intelligence** - Real user analytics and admin dashboard
- 🔒 **Compliance Ready** - Professional security practices

### For Scale
- 🔄 **Microservice Ready** - Works with any auth service
- 📈 **High Performance** - Optimized database operations and caching
- 🏗️ **Container Ready** - Docker and Kubernetes support
- 🔧 **Configurable** - Environment-based setup for any deployment
- 📊 **Monitorable** - Health checks and performance metrics

## 🏆 Current Status

**✅ 28/30 Tasks Completed (93% Progress)**

- **✅ Database Integration** - Full PostgreSQL with SQLC
- **✅ Authentication System** - Google OAuth with JWT
- **✅ Admin Dashboard** - Real user analytics and management
- **✅ Test Coverage** - Comprehensive authentication testing
- **✅ Docker Containerization** - Production deployment ready
- **✅ Real Data Integration** - No mock data, actual database metrics

**🔄 Remaining: Health monitoring and production optimization**

## 📄 License

MIT License - Use this platform freely for your startup projects.

## 🚀 Start Building

Ready to build your next startup with **modern architecture** and **production-ready features**?

```bash
git clone <your-fork>
cd go-templ-htmx-ex
make deps
make generate
make db-migrate
make dev
# Start building your amazing product! 🚀
```

### 🎯 **What You'll Get**

- **✅ Production Authentication** - Google OAuth with JWT security
- **✅ Database Integration** - PostgreSQL with type-safe SQLC queries
- **✅ Admin Dashboard** - Real user analytics and management tools
- **✅ Test Coverage** - Comprehensive testing for reliability
- **✅ Docker Ready** - Production deployment with containers
- **✅ Modern UI** - Professional design with HTMX + Templ
- **✅ Performance Optimized** - Fast queries and efficient caching

**Built with ❤️ for the startup community**
**Crafted with 🚀 modern architecture for 2025+**
