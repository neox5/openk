# ADR-011: Server Package Structure

## Status
Proposed

## Context
We need a clear structure for the HTTP server implementation that is simple, maintainable, and scalable.

## Decision
We will organize all HTTP server code under the internal/server package following a flat structure:

### Core Structure
```
internal/server/
├── error.go                 # Error types and handling
├── middleware/             # HTTP middleware
├── routes.go              # All route definitions
├── server.go              # Core server implementation
└── {domain}_v1.go        # Domain handlers with DTOs
```

### File Organization
Each domain handler file ({domain}_v1.go) contains:
- DTOs for that domain
- Handler implementation
- Version as part of filename

## Consequences

### Positive
- Clear location for all HTTP-related code
- Simple flat structure
- Easy to find files
- Clear versioning
- Minimizes package dependencies

### Negative
- All handlers in one package
- Need clear naming conventions

## Implementation Notes
1. Keep handlers focused on HTTP concerns
2. Version in filename (e.g., derivation_v1.go)
3. Follow RFC 7807 for error responses
4. Group related functionality in single files