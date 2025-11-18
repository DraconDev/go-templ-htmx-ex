# What Should Be in Middleware Directory

## ✅ **Proper Middleware Concerns** (Keep in middleware/)

### 1. `auth.go` - Authentication Middleware
- **Should have**: `AuthMiddleware`, `GetUserFromContext`, route categorization
- **Should NOT have**: HTTP client calls to external services

### 2. `cache.go` - Caching Middleware  
- **Should have**: SessionCache, caching logic, TTL management
- **Purpose**: Performance optimization across requests

### 3. `admin.go` - Authorization Middleware
- **Should have**: `RequireAdmin`, `RequireConfigAdmin` 
- **Purpose**: Role-based access control

## ❌ **NOT Middleware** (Move to services/)

### 4. `auth_http.go` - HTTP Client Logic
- **Problem**: Contains HTTP calls to auth service
- **Issue**: Middleware shouldn't make external HTTP requests
- **Solution**: Move to `services/auth_service.go`

### 5. `session.go` - Mixed Concerns  
- **Problem**: Contains both session validation AND HTTP calls
- **Issue**: Session validation logic mixed with external service calls
- **Solution**: Simplify to just session management, move HTTP calls to services

## **The Clean Middleware Pattern**

**Middleware should ONLY contain:**
1. **Request processing** (reading cookies, headers, etc.)
2. **Validation logic** (checking permissions, authentication status)  
3. **Response modification** (setting headers, redirects)
4. **Performance optimizations** (caching, rate limiting)

**Middleware should NOT contain:**
1. **Business logic** (OAuth flows, database queries)
2. **External HTTP calls** (API requests to other services)
3. **Complex orchestration** (multi-step processes)

## **Recommended Clean Structure**

```
internal/middleware/
├── auth.go      (request → validate session → add user context)
├── cache.go     (session caching, performance)
└── admin.go     (role-based authorization)

internal/services/  
├── auth_service.go (HTTP calls to auth microservice)
└── user_service.go (business logic)
```

## **Clear Separation Benefits**

✅ **Testability**: Middleware easy to test without external dependencies
✅ **Performance**: No external HTTP calls slowing down requests  
✅ **Maintainability**: Clear boundaries between layers
✅ **Reusability**: Middleware can be used by multiple handlers

## **Bottom Line**

**Middleware = Request/Response Processing**
**Services = Business Logic & External Calls**

Keep middleware thin, fast, and focused on cross-cutting concerns.