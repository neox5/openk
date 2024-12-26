# OpenK Development Guide

## Table of Contents
1. [Project Structure](#project-structure)
2. [Code Style Guide](#code-style-guide)
3. [gRPC Development](#grpc-development)
4. [Testing Guide](#testing-guide)

## Project Structure

```
.
├── proto/                      # Proto definitions
│   ├── buf.yaml               # Buf configuration
│   ├── buf.gen.yaml           # Code generation config
│   ├── openk/                 # API namespace
│   │   └── <service>/        # Service directories
│   │       ├── v1/           # Version directories
│   │       │   ├── service.proto  # Service definitions
│   │       │   └── types.proto    # Type definitions
│   │       └── v2/
│   └── vendor/               # Vendored dependencies
│       └── google/
│           └── protobuf/
├── internal/
│   ├── api_gen/             # Generated code
│   │   └── openk/
│   │       └── <service>/
│   │           └── v{n}/
│   └── server/              # Server implementation
│       ├── grpc_server.go   # gRPC server setup
│       ├── gateway.go       # REST gateway
│       ├── interceptors/    # gRPC interceptors
│       ├── services/        # Service implementations
│       └── health/          # Health checks
├── cmd/                     # Command line applications
│   └── openk/
│       └── main.go
└── docs/                    # Documentation
    ├── adr/                # Architecture Decision Records
    └── specs/              # Technical specifications
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

### Buf Configuration

#### buf.yaml
```yaml
version: v2
modules:
  - path: .
    name: buf.build/neox5/openk
    excludes:
      - vendor/google/protobuf
lint:
  use:
    - DEFAULT
  except:
    - PACKAGE_VERSION_SUFFIX
breaking:
  use:
    - FILE
```

#### buf.gen.yaml
```yaml
version: v2
clean: true
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/neox5/openk/internal/api_gen
plugins:
  - remote: buf.build/protocolbuffers/go:v1.31.0
    out: ../internal/api_gen
    opt: 
      - paths=source_relative
  - remote: buf.build/grpc/go:v1.3.0
    out: ../internal/api_gen
    opt:
      - paths=source_relative
```

### Key Configuration Settings

#### buf.yaml
Key settings:
- `modules`: Local module configuration without external dependencies
- `lint`: Enforces API design standards (DEFAULT ruleset)
- `breaking`: Detects breaking API changes using FILE ruleset
- `excludes`: Prevents vendor directory from being processed

#### buf.gen.yaml
Key settings:
- `clean`: Clears output directory before generation
- `managed`: Enables consistent package naming
- `plugins`: Configures Go and gRPC code generation
- `out`: Places generated code in internal/api_gen

### Package Naming
- Format: `openk.<service>.<version>`
- Example: `openk.health.v1`
- Never reuse package names across different protofiles
- Always include version in package name

### Go Package Naming
- Format: `github.com/neox5/openk/internal/api_gen/openk/<service>/<version>;<service><version>`
- Example: `github.com/neox5/openk/internal/api_gen/openk/health/v1;healthv1`
- Import path matches generated code location
- Short package names for clean imports

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
