# OpenK Development Guide

## Table of Contents
1. [Project Structure](#project-structure)
2. [Code Style Guide](#code-style-guide)
3. [gRPC Development](#grpc-development)
4. [Testing Guide](#testing-guide)

## Project Structure

```
.
├── cmd/                     # Command line applications
│   └── openk/
│       └── main.go
├── docs/                    # Documentation
│   ├── adr/                # Architecture Decision Records
│   ├── misc/               # Miscellaneous documentation
│   ├── models/             # Data model documentation
│   └── specs/              # Technical specifications
├── internal/               # Internal packages
│   ├── api_gen/           # Generated code
│   │   └── openk/         # Generated API code
│   │       ├── common/    # Common proto types
│   │       │   └── v1/
│   │       └── health/    # Health service protos
│   │           ├── v1/    
│   │           └── v2/    
│   ├── app/               # Application core
│   │   ├── client/       # Client implementations
│   │   └── server/       # Server core
│   ├── buildinfo/         # Build information
│   ├── cli/              # CLI implementation
│   ├── crypto/           # Cryptographic operations
│   ├── ctx/              # Context utilities
│   ├── kms/              # Key management
│   ├── logging/          # Logging utilities
│   ├── opene/            # Error handling
│   ├── secret/           # Secret management
│   └── server/           # Server implementation
│       ├── interceptors/ # gRPC interceptors
│       └── services/     # Service implementations
│           └── health/   # Health service
│               ├── health_server_v1.go  # V1 implementation
│               └── health_register.go   # Version registration
├── proto/                 # Proto definitions
│   ├── openk/            # API namespace
│   │   ├── common/       # Common definitions
│   │   │   └── v1/
│   │   │       └── common_v1.proto    # Common types
│   │   └── health/       # Health service
│   │       └── v1/
│   │           ├── health_v1.proto          # V1 types
│   │           └── health_service_v1.proto  # V1 service
│   ├── vendor/           # Vendored proto dependencies
│   ├── buf.yaml          # Buf configuration
│   └── buf.gen.yaml      # Code generation config
└── scripts/              # Build and maintenance scripts
```

## Code Style Guide

### Code Documentation
- Write clear package documentation
- Document exported symbols
- Include examples for complex functionality
- Follow godoc conventions
- Include usage examples for complex functionality

### Security Considerations
- Clear sensitive data after use
- Use constant-time comparisons
- Validate all inputs
- Handle errors securely
- Follow crypto best practices

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

Functions should take context when they:
- Make network calls
- Access databases or external storage
- Perform long-running operations
- Need cancellation/timeout support
- Propagate request-scoped data

### Logging

#### Using LogAttrs Pattern
Always use LogAttrs for logging to ensure efficiency. Choose the appropriate logger based on your component:

```go
// In server components (HTTP handlers, middleware)
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    h.logger.LogAttrs(ctx, slog.LevelInfo, "handling request",
        slog.String("path", r.URL.Path),
        slog.String("method", r.Method),
    )
}

// In independent packages (crypto, kms)
func (k *KeyManager) Rotate(ctx context.Context, keyID string) error {
    slog.Default().LogAttrs(ctx, slog.LevelInfo, "rotating key",
        slog.String("key_id", keyID),
        slog.Time("rotation_time", time.Now()),
    )
}
```

### Error Handling

All errors should be created using OpenE. This provides consistent error handling and RFC 7807 compliance.

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

// HTTP Error Handling
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := h.store.FetchKey(r.Context(), r.URL.Query().Get("id")); err != nil {
        prob := opene.AsProblem(err)
        w.Header().Set("Content-Type", "application/problem+json")
        w.WriteHeader(prob.Status)
        json.NewEncoder(w).Encode(prob)
        return
    }
}
```

## gRPC Development

### Proto Organization

#### Directory Structure
```
proto/openk/
├── common/              # Shared types
│   └── v1/
│       └── common_v1.proto
└── <service>/          # Service-specific protos
    └── v1/
        ├── <service>_v1.proto           # Service types
        └── <service>_service_v1.proto   # Service definition
```

#### Naming Conventions
1. Proto Files:
   - Types: `<service>_v1.proto` (e.g., `health_v1.proto`)
   - Service: `<service>_service_v1.proto` (e.g., `health_service_v1.proto`)
   - Common: `common_v1.proto`

2. Service Implementation:
   - Server: `<service>_server_v1.go` (e.g., `health_server_v1.go`)
   - Registration: `<service>_register.go` (e.g., `health_register.go`)

### Server Implementation Pattern

#### File Structure
```go
// health_server_v1.go
type HealthServerV1 struct {
    healthv1.UnimplementedHealthServiceServer
    // Server implementation
}

// health_register.go
func RegisterHealthServers(srv *grpc.Server, logger *slog.Logger) (*HealthServerV1, error) {
    // Registration of all versions
}
```

#### Version Management
- Each version has its own server implementation (`*_server_v1.go`)
- Single registration point for all versions (`*_register.go`)
- Clear version suffixes in type names (e.g., `HealthServerV1`)
- Version-specific types include version in name (e.g., `ComponentCheckV1`)

### Proto Best Practices

1. Package Naming:
   ```protobuf
   package openk.<service>.v1;
   option go_package = "github.com/neox5/openk/internal/api_gen/openk/<service>/v1;<service>v1";
   ```

2. Service Definition:
   ```protobuf
   service HealthService {
     rpc Check(CheckRequest) returns (CheckResponse);
   }
   ```

3. Version Management:
   - Each version in separate directory (`v1/`, `v2/`)
   - No breaking changes within a version
   - New versions for breaking changes

### Implementation Guidelines

1. Error Handling:
   ```go
   return nil, opene.NewInternalError("health", "check", "database unavailable")
   ```

2. Logging:
   ```go
   s.logger.LogAttrs(ctx, slog.LevelInfo, "health check requested",
       slog.String("version", "v1"),
       slog.Int("components", len(req.Components)),
   )
   ```

3. Context Usage:
   - Always pass context through gRPC methods
   - Use for cancellation and timeouts
   - Include in logging calls

4. Registration:
   ```go
   server := grpc.NewServer()
   v1Server, err := health.RegisterHealthServers(server, logger)
   if err != nil {
       return err
   }
   ```

### Testing Approach

1. Server Tests:
   ```
   services/health/
   ├── health_server_v1_test.go
   └── health_register_test.go
   ```

2. Integration Tests:
   - Use real gRPC client/server
   - Test all supported versions
   - Verify version compatibility

3. Test Naming:
   ```go
   func TestHealthServerV1_Check(t *testing.T) {...}
   func TestHealthServerV1_WatchHealth(t *testing.T) {...}
   ```

## Testing Guide

### File Structure
- Name test files: `foo_test.go` for `foo.go`
- Use separate test package: `package foo_test`
- Group imports with stdlib first

### Test Function Naming
```go
// Standard pattern 
TestType_Method(t *testing.T)

// Examples
TestRSA_Generate(t *testing.T)
TestMemory_SecureWipe(t *testing.T)
```

### Subtests
Format: Descriptive phrase in lowercase, grouped by success/error cases
```go
t.Run("success cases", func(t *testing.T) {
    t.Run("valid key is exported", func(t *testing.T) {...})
    t.Run("empty message is handled", func(t *testing.T) {...})
})

t.Run("error cases", func(t *testing.T) {
    t.Run("rejects nil key", func(t *testing.T) {...})
    t.Run("handles oversized input", func(t *testing.T) {...})
})
```

### Complex Test Cases
For tests with complex setup or multiple steps, prefer direct testing over tables:

```go
func TestRSA_ImportPrivateKey(t *testing.T) {
    t.Run("success cases", func(t *testing.T) {
        key := generateTestKey(t, RSAKeySize2048)
        der, err := MarshalPKCS8PrivateKey(key)
        require.NoError(t, err)

        imported, err := ImportRSAPrivateKey(der)
        assert.NoError(t, err)
        assert.Equal(t, key.D.Bytes(), imported.D.Bytes())
    })

    t.Run("error cases", func(t *testing.T) {
        t.Run("rejects invalid DER", func(t *testing.T) {
            imported, err := ImportRSAPrivateKey([]byte("invalid"))
            assert.ErrorIs(t, err, ErrInvalidPrivateKey)
            assert.Contains(t, err.Error(), "invalid format")
            assert.Nil(t, imported)
        })
    })
}
```

### Test Helpers
- Put common setup in helper functions
- Mark helpers with t.Helper()
- Keep helpers focused and simple
- Place helpers close to where they're used

Example:
```go
func generateTestKey(t *testing.T, size int) *rsa.PrivateKey {
    t.Helper()
    key, err := rsa.GenerateKey(rand.Reader, size)
    require.NoError(t, err)
    return key
}
```

### Test Structure

#### Basic Pattern
```go
func TestType_Method(t *testing.T) {
    // Optional common setup
    key := generateTestKey(t)

    t.Run("success cases", func(t *testing.T) {
        // Test happy paths
    })

    t.Run("error cases", func(t *testing.T) {
        // Test error conditions
    })
}
```

#### Table-Driven Tests
Use only when testing multiple similar input/output variations:

```go
func TestRSA_Generate(t *testing.T) {
    tests := []struct {
        name    string
        bits    int
    }{
        {
            name: "generates 2048-bit key",
            bits: RSAKeySize2048,
        },
        {
            name: "generates 4096-bit key",
            bits: RSAKeySize4096,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            key, err := GenerateRSAKeyPair(tt.bits)
            assert.NoError(t, err)
            assert.Equal(t, tt.bits, key.Size()*8)
        })
    }
}
```

### Error Checking
- Use `require` for critical setup that must succeed
- Use `assert` for test validations
- Check both error type and message when relevant

```go
// Setup with require
key := generateTestKey(t)
der, err := MarshalPrivateKey(key)
require.NoError(t, err)

// Test validations with assert
result, err := ImportKey(der)
assert.ErrorIs(t, err, ErrInvalidKey)
assert.Contains(t, err.Error(), "invalid format")
```

### Best Practices (in order of priority)

1. Keep tests focused and small
   - Each test should validate one piece of functionality
   - Split complex tests into smaller, focused subtests
   - Keep setup close to where it's used

2. Group logically with t.Run()
   - Separate success and error cases
   - Use descriptive names for test groups
   - Keep nesting to maximum of 2 levels

3. Use table-driven tests appropriately
   - Only for similar input/output variations
   - Avoid for complex setup or multi-step tests
   - Keep test cases clear and readable

4. Use require vs assert correctly
   - require: For critical setup that must succeed
   - assert: For actual test validations
   - Fail fast when setup fails

### Running Tests

```bash
go test ./...                     # All tests
go test -run TestType_Method ./...   # Single test
go test -v ./...                  # Verbose output
```
