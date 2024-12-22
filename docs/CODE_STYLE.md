# OpenK Code Style Guide

## Table of Contents
1. [Logging](#logging)
2. [Error Handling](#error-handling)
3. [Code Documentation](#code-documentation)
4. [Security Considerations](#security-considerations)

## Logging

### Using LogAttrs Pattern
Always use LogAttrs for logging to ensure efficiency. Choose the appropriate logger based on your component:

```go
// In server components (HTTP handlers, middleware)
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // Use instance logger when available
    h.logger.LogAttrs(ctx, slog.LevelInfo, "handling request",
        slog.String("path", r.URL.Path),
        slog.String("method", r.Method),
        slog.String("remote_addr", r.RemoteAddr),
    )

    // Process request...

    h.logger.LogAttrs(ctx, slog.LevelInfo, "request complete",
        slog.Int("status", status),
        slog.Duration("duration", time.Since(start)),
        slog.Int64("bytes_written", written),
    )
}

// In independent packages (crypto, kms)
func (k *KeyManager) Rotate(ctx context.Context, keyID string) error {
    // Use slog directly in packages
    slog.Default().LogAttrs(ctx, slog.LevelInfo, "rotating key",
        slog.String("key_id", keyID),
        slog.Time("rotation_time", time.Now()),
    )

    result, err := k.doRotation(ctx, keyID)
    if err != nil {
        slog.Default().LogAttrs(ctx, slog.LevelError, "key rotation failed",
            slog.String("key_id", keyID),
            slog.String("error", err.Error()),
        )
        return err
    }

    slog.Default().LogAttrs(ctx, slog.LevelInfo, "key rotation complete",
        slog.String("key_id", keyID),
        slog.Time("completion_time", result.CompletedAt),
        slog.Duration("duration", result.Duration),
    )
    return nil
}
```

### Context Usage
Always pass context when the function receives it:

```go
// Example showing when to include context
func (s *Server) Start(ctx context.Context) error {
    // Operation with context - pass it
    s.logger.LogAttrs(ctx, slog.LevelInfo, "starting server",
        slog.String("address", s.addr),
        slog.Duration("shutdown_timeout", s.shutdownTimeout),
    )
    
    if err := s.startServer(); err != nil {
        s.logger.LogAttrs(ctx, slog.LevelError, "server start failed",
            slog.String("error", err.Error()),
        )
        return err
    }

    return nil
}

// Example showing when not to include context
func (s *Server) Version() string {
    // No context available - use nil
    s.logger.LogAttrs(nil, slog.LevelDebug, "version requested",
        slog.String("version", s.version),
    )
    return s.version
}
```

## Error Handling

All errors should be created using OpenE. This provides consistent error handling, proper context, and RFC 7807 compliance.

### Creating Errors
```go
// Validation errors
func ValidateUsername(username string) error {
    if username == "" {
        return opene.NewValidationError("username cannot be empty").
            WithDomain("auth").
            WithOperation("validate_username").
            WithMetadata(opene.Metadata{
                "field": "username",
            })
    }
    return nil
}

// Not found errors
func (s *Store) FetchKey(ctx context.Context, id string) (*Key, error) {
    key, exists := s.keys[id]
    if !exists {
        return nil, opene.NewNotFoundError("key not found").
            WithDomain("kms").
            WithOperation("fetch_key").
            WithMetadata(opene.Metadata{
                "key_id": id,
            })
    }
    return key, nil
}

// Internal errors
func (k *KeyManager) Rotate(ctx context.Context, keyID string) error {
    if err := k.store.Rotate(ctx, keyID); err != nil {
        return opene.NewInternalError("key rotation failed").
            WithDomain("kms").
            WithOperation("rotate_key").
            WithMetadata(opene.Metadata{
                "key_id": keyID,
            }).
            Wrap(opene.AsError(err, "storage", opene.CodeInternal))
    }
    return nil
}
```

### HTTP Error Handling
Always use AsProblem for HTTP responses:

```go
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    key, err := h.store.FetchKey(r.Context(), r.URL.Query().Get("id"))
    if err != nil {
        prob := opene.AsProblem(err)
        w.Header().Set("Content-Type", "application/problem+json")
        w.WriteHeader(prob.Status)
        json.NewEncoder(w).Encode(prob)
        return
    }
    // Handle success...
}
```

## Code Documentation
- Write clear package documentation
- Document exported symbols
- Include examples for complex functionality
- Follow godoc conventions

## Security Considerations
- Clear sensitive data after use
- Use constant-time comparisons
- Validate all inputs
- Handle errors securely
- Follow crypto best practices
