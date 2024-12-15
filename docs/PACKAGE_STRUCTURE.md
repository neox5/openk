# Server Package Structure Guide

## Overview
This guide explains the organization of the HTTP server implementation.

## Structure
```
internal/server/
├── error.go                 # Error types and handling
├── middleware/             # HTTP middleware
├── routes.go              # All route definitions
├── server.go              # Core server implementation
└── {domain}_v1.go        # Domain handlers with DTOs
```

## File Organization

### Domain Handlers
Each `{domain}_v1.go` file contains:
```go
package server

// DTOs
type storeParamsRequest struct {
    // Request fields
}

type paramsResponse struct {
    // Response fields
}

// Handler
type domainV1Handler struct {
    // Dependencies
}

func (h *domainV1Handler) handleOperation(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### Routes
All routes defined in `routes.go`:
```go
func (s *Server) registerRoutes() {
    handler := newDomainV1Handler(...)
    s.mux.Handle("/api/v1/...", middleware.ErrorHandler(
        http.HandlerFunc(handler.handleOperation)))
}
```

## Guidelines
- One file per domain area
- Include version in filename
- Group related DTOs with handler
- Use middleware for cross-cutting concerns
- Follow RFC 7807 for errors