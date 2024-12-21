# ADR-012: Logging Architecture

## Status
Proposed

## Context
OpenK needs a consistent logging approach that allows for:
- Standard logging interface across all packages
- Flexibility to change logger implementations at the server level
- Consistent structured logging patterns
- Proper context propagation
- Performance monitoring capabilities

## Decision

### 1. Core Architecture
- Use Go's slog as the standard logging interface across all packages
- Allow logger implementation swapping via slog.Handler interface
- Configure logging at server initialization
- Maintain consistent structured logging patterns

### 2. Package Requirements

#### 2.1 Internal Packages (crypto, kms, etc.)
- Must use slog.Logger directly
- Must not depend on specific logger implementations
- Must follow structured logging patterns
```go
slog.Info("rotating key",
    "key_id", id,
    "algorithm", alg,
)
```

#### 2.2 Server Layer
- Controls logger implementation selection
- Configures global logger via slog.SetDefault
- Provides logging middleware
- Handles context propagation

### 3. Logging Configuration
```go
type LogConfig struct {
    // Minimum level to log
    Level slog.Level

    // Logger implementation to use
    Implementation string // "slog" | "zerolog" | etc.

    // Output format
    Format string // "json" | "text"

    // Additional handler options
    HandlerOptions map[string]any
}
```

### 4. Context Values
Standard context keys for logging:
- TraceID
- SpanID
- RequestID
- UserID
- TenantID

### 5. Structured Logging Standards

#### 5.1 Field Naming
- Use snake_case for field names
- Use descriptive but concise names
- Include type indicators for special values

#### 5.2 Required Fields
Each log entry must include:
- Timestamp
- Log level
- Message
- Relevant context IDs

#### 5.3 Error Logging
Error logs must include:
- Error message
- Error type
- Stack trace (when available)
- Operation context

### 6. Implementation Guidelines

#### 6.1 Logger Setup
```go
func initLogger(cfg LogConfig) *slog.Logger {
    var handler slog.Handler
    switch cfg.Implementation {
    case "zerolog":
        zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
        handler = slogzerolog.Option{
            Level: cfg.Level,
            Logger: &zl,
        }.NewZerologHandler()
    default:
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: cfg.Level,
        })
    }
    return slog.New(handler)
}
```

#### 6.2 HTTP Middleware Pattern
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        logger := slog.With(
            "trace_id", ctx.Value(KeyTraceID),
            "request_id", ctx.Value(KeyRequestID),
            "path", r.URL.Path,
            "method", r.Method,
        )
        
        // ... request handling
    })
}
```

## Consequences

### Positive
- Consistent logging interface across codebase
- Flexibility to change implementations
- Standard library compatibility
- Clear structured logging patterns
- Performance optimized by default

### Negative
- Need to maintain consistent logging patterns
- Must ensure proper context propagation
- Some implementation-specific features may not be available
- Additional configuration complexity at server level

## Implementation Notes
1. Document standard logging patterns
2. Create logging middleware
3. Define context propagation approach
4. Setup proper log levels and formats
5. Monitor logging performance

## References
- slog documentation
- zerolog integration
- OpenTelemetry logging standards