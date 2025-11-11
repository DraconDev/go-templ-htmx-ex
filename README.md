# 🚀 Go + HTMX + Templ Startup Platform

A **minimal, production-ready startup platform** with **Google OAuth**, **PostgreSQL database**, and **admin dashboard**. Built with **Templ**, **HTMX**, and **SQLC** for high performance.

## 🎯 What This Is

- **🚀 Fast startup foundation** with real authentication & database
- **📊 Simple admin dashboard** with user analytics  
- **🔐 Google OAuth ready** with JWT sessions
- **🧪 Basic test coverage** for core functionality
- **🐳 Docker ready** for production deployment

## ✨ What You Get

### 🔐 **Authentication System**
- Google OAuth 2.0 login
- JWT session management
- User profile pages
- Logout functionality

### 💾 **Database Integration**
- PostgreSQL with users table
- SQLC generated type-safe queries
- Real user data (no mock data)
- Basic user analytics

### 📊 **Admin Dashboard** 
- Total users count
- Signups today/this week
- Recent users list
- Admin-only access

### 🎨 **Clean UI**
- Simple dark theme
- Responsive design
- HTMX interactions
- Modern but minimal

## 🚀 Quick Start

```bash
# Clone and setup
git clone <your-repo>
cd go-templ-htmx-ex

# Install dependencies
make deps

# Generate templates
make generate

# Setup database
createdb startup_platform
make db-migrate

# Run development
make dev
```

**Visit:** `http://localhost:8081`

## 🔧 Configuration

```bash
# Copy environment config
cp .env.example .env

# Edit these values:
# PORT=8081
# AUTH_SERVICE_URL=http://localhost:8080  # Your auth service
# DATABASE_URL=postgresql://user:pass@localhost:5432/dbname
# ADMIN_EMAIL=admin@yourdomain.com
```

## 📁 Project Structure

```
├── main.go              # Application entry point
├── Dockerfile           # Production container
├── sqlc.yaml           # Database query generation
├── handlers/
│   ├── admin.go        # Admin dashboard
│   ├── auth.go         # Authentication
│   └── handlers.go     # User pages
├── db/
│   ├── migrations/     # Database schema
│   └── sqlc/          # Generated queries
├── middleware/
│   └── auth.go        # JWT validation
└── templates/
    ├── layout.templ   # Base layout
    ├── home.templ     # Landing page
    ├── profile.templ  # User profile
    └── admin_dashboard.templ # Admin page
```

## 🧪 Testing

```bash
# Run tests
make test

# Output shows authentication flow tests passing
go test ./handlers/ -v
```

## 🐳 Docker

```bash
# Build and run
make docker-build
docker run -p 8081:8081 your-app
```

## 📊 Current Features

### ✅ **What Works**
- Google OAuth login/logout
- User profile pages with real data
- Admin dashboard with database statistics
- PostgreSQL database integration
- Basic test coverage
- Docker containerization

### 🔄 **What's Missing**
- Mobile app API endpoints
- Payment integration
- Email notifications
- Advanced analytics

## 📈 Performance

- **Navigation:** ~5-10ms with local JWT validation
- **Admin Dashboard:** Real-time database queries
- **Database:** SQLC generated optimized queries
- **UI:** HTMX for seamless updates

## 🚀 For Your Startup

This gives you the foundation to build on:

```bash
# Add your business features
mkdir handlers/business
vim handlers/business/your_feature.go

# Add database tables
vim db/migrations/002_your_feature.sql

# Create templates
vim templates/your_feature.templ
```

## 📄 License

MIT License

---

**Simple. Fast. Ready to build on.**
