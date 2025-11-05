# gRPC Auth Integration with Cerberus Authentication Service

This document describes the integration of the `auth-cerberus` gRPC authentication service into the go-templ-htmx-ex project.

## Overview

The integration provides a complete gRPC client for the Cerberus authentication service, enabling the application to communicate with a remote auth server over gRPC instead of using HTTP-based authentication.

## Architecture

### Components

1. **gRPC Client (`auth/grpc_client.go`)**: Complete implementation of gRPC client for Cerberus auth service
2. **Generated gRPC Files (`auth/cerberus/`)**: Auto-generated protobuf bindings for the auth service
3. **Test Client (`cmd/grpc-auth-test/main.go`)**: Standalone test application to verify gRPC functionality
4. **Proto Definition (`proto/auth.proto`)**: Source protobuf definition for the auth service

### API Endpoints

The gRPC client provides the following authentication methods:

- **Login**: Authenticate user with email/password
- **Register**: Create new user account with project association
- **ValidateSession**: Verify session token validity
- **GetUserDetails**: Retrieve user information by user ID
- **HealthCheck**: Check if the auth service is available

## Setup and Configuration

### Dependencies

The gRPC client requires the following Go dependencies:

```go
google.golang.org/grpc v1.76.0
google.golang.org/protobuf v1.36.10
```

These are automatically resolved through the project's `go.mod` file.

### Connection Configuration

The gRPC client is configured to connect to the auth service:

```go
// Insecure connection for development
// TODO: Add proper TLS certificates for production
conn, err := grpc.NewClient(
    baseURL,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithBlock(),
)
```

### Environment Variables

The following environment variables should be configured:

- `AUTH_GRPC_URL`: URL of the Cerberus auth service (e.g., `localhost:50051`)
- `AUTH_GRPC_TIMEOUT`: Connection timeout (default: 10 seconds)

## Usage Examples

### Basic Usage

```go
package main

import (
    "log"
    "github.com/dracon/go-templ-htmx-ex/auth"
)

func main() {
    // Initialize gRPC client
    authClient, err := auth.NewGRPCAuthClient("localhost:50051")
    if err != nil {
        log.Fatalf("Failed to create auth client: %v", err)
    }
    defer authClient.Close()

    // Login user
    resp, err := authClient.Login("user@example.com", "password")
    if err != nil {
        log.Printf("Login failed: %v", err)
        return
    }

    log.Printf("Login successful: %+v", resp)
}
```

### Error Handling

The gRPC client provides comprehensive error handling:

```go
// Check for gRPC-specific errors
if err != nil {
    st, ok := status.FromError(err)
    if ok {
        log.Printf("gRPC error (%s): %s", st.Code(), st.Message())
    } else {
        log.Printf("General error: %v", err)
    }
    return
}
```

### Response Format

All auth methods return a `GRPCAuthResponse` structure:

```go
type GRPCAuthResponse struct {
    UserID       string   `json:"user_id"`
    SessionToken string   `json:"session_token"`
    Email        string   `json:"email,omitempty"`
    Valid        bool     `json:"valid,omitempty"`
    ProjectIDs   []string `json:"project_ids,omitempty"`
    Error        string   `json:"error,omitempty"`
    Status       string   `json:"status,omitempty"`
    Message      string   `json:"message,omitempty"`
}
```

## Testing

### Running the Test Client

```bash
go run cmd/grpc-auth-test/main.go
```

The test client will attempt to connect to `localhost:50051` and test all available authentication methods. Connection failures are expected if the auth server is not running.

### Expected Output

```
Creating gRPC auth client...
gRPC auth client created successfully!
Testing health check...
Health check failed: connection error (expected if auth server isn't running)
Testing login...
Login failed: connection error (expected if auth server isn't running)
Testing register...
Register failed: connection error (expected if auth server isn't running)
Testing validate session...
Validate session failed: connection error (expected if auth server isn't running)
Testing get user details...
Get user details failed: connection error (expected if auth server isn't running)
gRPC client test completed!
```

## Production Considerations

### Security

1. **TLS Configuration**: Replace `insecure.NewCredentials()` with proper TLS certificates
2. **Authentication**: Consider mutual TLS (mTLS) for enhanced security
3. **Connection Pooling**: Implement connection pooling for high-throughput scenarios

### Performance

1. **Connection Reuse**: The gRPC client maintains persistent connections
2. **Timeouts**: Configurable timeouts for all operations
3. **Retry Logic**: Implement exponential backoff for transient failures

### Monitoring

1. **Health Checks**: Implement regular health checks to detect service availability
2. **Metrics**: Add metrics collection for monitoring auth service performance
3. **Logging**: Structured logging for debugging and monitoring

## Troubleshooting

### Common Issues

1. **Connection Refused**: Ensure the auth service is running on the specified host:port
2. **TLS Errors**: Check certificate configuration and hostname verification
3. **Timeout Errors**: Verify network connectivity and increase timeout values
4. **Authentication Failures**: Verify user credentials and service configuration

### Debug Mode

Enable debug logging by setting:

```go
logging.SetLevel(logging.DEBUG)
```

## Future Enhancements

1. **Connection Pooling**: Implement connection pooling for better performance
2. **Load Balancing**: Add client-side load balancing for multiple auth service instances
3. **Caching**: Implement session token caching to reduce auth service calls
4. **Metrics**: Add Prometheus metrics for monitoring
5. **Circuit Breaker**: Implement circuit breaker pattern for resilience

## Migration from HTTP Auth

To migrate from HTTP-based authentication to gRPC:

1. Replace HTTP client calls with gRPC client methods
2. Update error handling to use gRPC status codes
3. Remove HTTP-specific dependencies
4. Update configuration management for gRPC endpoints

## Support

For issues related to the gRPC integration:

1. Check the auth service logs for server-side errors
2. Verify network connectivity and firewall rules
3. Review the proto definitions for API compatibility
4. Check the test client output for specific error messages
