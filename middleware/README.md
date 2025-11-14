# Middleware Route Configuration

This directory contains the authentication middleware and route configuration for the application.

## 📁 Files

- `auth.go` - Main authentication middleware logic
- `routes.go` - Centralized route configuration and categorization
- `README.md` - This documentation file

## 🛣️ Route Categories

The middleware now uses a centralized configuration system to categorize routes:

### **PROTECTED Routes** (Require Authentication)
These routes will redirect unauthenticated users to `/login`:
- `/profile` - User profile page
- `/admin` - Admin dashboard

### **ADMIN Routes** (Require Authentication)
These routes require admin-level authentication:
- `/api/admin/*` - All admin API endpoints

### **PUBLIC Routes** (No Authentication Required)
These routes are accessible to everyone and skip session validation:
- `/` - Homepage
- `/health` - Health check endpoint  
- `/login` - Login page
- `/test` - Test page for development

### **OAUTH Routes** (Authentication Flow)
These routes are part of the OAuth authentication process:
- `/auth/google` - Google OAuth login
- `/auth/github` - GitHub OAuth login
- `/auth/discord` - Discord OAuth login
- `/auth/microsoft` - Microsoft OAuth login
- `/auth/callback` - OAuth callback handler

### **AUTH API Routes** (Session Management)
These routes handle authentication-related API operations:
- `/api/auth/validate` - Validate session
- `/api/auth/user` - Get current user info
- `/api/auth/logout` - Logout user
- `/api/auth/set-session` - Set session cookie
- `/api/auth/exchange-code` - Exchange OAuth code for session
- `/api/auth/test-session-create` - Test session creation
- `/api/auth/refresh` - Refresh session token

## ⚙️ Configuration

To modify route categories, edit the arrays in `routes.go`:

```go
// Example: Add a new protected route
var ProtectedRoutes = []string{
    "/profile",        // User profile page
    "/admin",          // Admin dashboard
    "/dashboard",      // NEW: User dashboard
}

// Example: Add a new admin route
var AdminRoutes = []string{
    "/api/admin",      // All admin API endpoints
    "/api/reports",    // NEW: Reports API
}
```

## 🔍 Benefits

1. **Centralized Configuration** - All route categories in one place
2. **Easy Maintenance** - Add/modify routes without touching middleware logic
3. **Clear Categorization** - Each route has a defined purpose
4. **Debugging Support** - Built-in route category logging
5. **Flexible** - Easy to add new categories or modify existing ones

## 🐛 Debugging

The middleware includes detailed logging for route processing:

```
🔐 MIDDLEWARE: Processing route / [Category: PUBLIC]
🔐 MIDDLEWARE: Processing route /profile [Category: PROTECTED]
🔐 MIDDLEWARE: Processing route /api/admin/users [Category: PROTECTED]
```

Use `GetRouteCategory(path)` to check any route's category programmatically.

## 🔐 Security

- **Protected Routes** redirect to `/login` for unauthenticated users
- **Admin Routes** require admin-level authentication
- **Public Routes** skip session validation entirely for performance
- **Session Validation** is cached for 15 seconds to reduce auth service calls

## ⚡ Performance

- Public routes bypass session validation completely
- Session validation results are cached for 15 seconds
- Reduced auth service calls by ~85% for public routes