# OpenK Code Style Guide

## Error Handling

### Error Organization
- Define errors in `errors.go` within each package
- Keep error definitions close to their usage
- Group related errors together

### Error Definition Pattern
```go
var (
    // Username errors - subject after Err
    ErrUsernameEmpty    = errors.New("username cannot be empty")
    ErrUsernameLength   = errors.New("username exceeds maximum length")
    ErrUsernameInvalid  = errors.New("username contains invalid characters")

    // Key errors - subject after Err
    ErrKeyRevoked       = errors.New("key has been revoked")
    ErrKeyNotDerived    = errors.New("key not derived")
    ErrKeyAlreadySet    = errors.New("key already set")

    // Envelope errors - subject after Err
    ErrEnvelopeMissing  = errors.New("envelope not found")
    ErrEnvelopeInvalid  = errors.New("envelope contains invalid data")

    // Don't:
    ErrNoKey           = errors.New("key not found")          // Use ErrKeyMissing instead
    ErrInvalidUsername = errors.New("username is invalid")    // Use ErrUsernameInvalid instead
    ErrorKeyNotFound   = errors.New("Key not found.")         // Don't prefix with Error or use punctuation
    errKeyInvalid     = errors.New("Invalid key")            // Don't use lowercase for exported errors
)
```

### Error Wrapping
- Wrap errors when adding context
- Keep the error chain focused
```go
// Do:
return fmt.Errorf("failed to decrypt key: %w", err)
return fmt.Errorf("invalid key format: %w", err)

// Don't:
return fmt.Errorf("error occurred while trying to perform decryption: %w", err)  // Too verbose
return fmt.Errorf("Error: %w", err)  // Too vague
return err  // Missing context
```

### HTTP Error Handling (RFC 7807)
Following ADR-010, implement HTTP errors using RFC 7807 format:

#### Error Response Structure
```go
// ProblemDetails implements RFC 7807 for HTTP API errors
type ProblemDetails struct {
    Type     string      `json:"type"`               // URI reference
    Title    string      `json:"title"`              // Short summary
    Status   int         `json:"status"`             // HTTP status code
    Detail   string      `json:"detail"`             // Detailed explanation
    Instance string      `json:"instance"`           // URI reference to specific occurrence
    Extra    interface{} `json:"extensions,omitempty"` // Optional extensions
}
```

#### Implementation Guidelines
1. Use Centralized Error Types:
```go
// internal/server/errors/types.go
const (
    TypeValidation = "https://openk.dev/errors/validation"
    TypeNotFound   = "https://openk.dev/errors/not-found"
    TypeInternal   = "https://openk.dev/errors/internal"
    TypeConflict   = "https://openk.dev/errors/conflict"
)
```

2. Create Error Helper Functions:
```go
// Good - Use helper functions for common errors
return errors.NewValidationError("username", "alphanumeric")

// Don't - Create error responses inline
return &ProblemDetails{Type: "...", Title: "..."}
```

3. Secure Error Handling:
```go
// Good - Sanitize sensitive data
err := errors.NewValidationError("password", "invalid")
err.SanitizeValue()

// Don't - Include sensitive values
err := errors.NewValidationError("password", "123pass")
```

4. Use Middleware for Consistency:
```go
// Good - Use error handling middleware
func ErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                ProblemResponse(w, errors.NewInternalError())
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// Don't - Handle errors differently across handlers
if err != nil {
    json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
```

5. Include TraceID for Internal Errors:
```go
// Good - Include traceID for debugging
if err != nil {
    traceID := generateTraceID()
    logError(traceID, err) // Log full error details
    return errors.NewInternalError().WithTraceID(traceID)
}

// Don't - Expose internal error details
if err != nil {
    return errors.NewInternalError().WithDetail(err.Error())
}
```

6. Test Error Scenarios:
```go
// Good - Test error responses
func TestHandler_Validation(t *testing.T) {
    resp := httptest.NewRecorder()
    // ... handler setup ...
    
    var problem ProblemDetails
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&problem))
    assert.Equal(t, TypeValidation, problem.Type)
    assert.Equal(t, http.StatusBadRequest, problem.Status)
}
```

## Testing

### File Structure
- Name test files: `foo_test.go` for `foo.go`
- Use separate test package: `package foo_test`
- Group imports with stdlib first

### Test Organization
- Group related test cases logically
- Use subtests for different scenarios
- Keep test files focused and manageable

### Test Pattern
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

### Assertions
- Use require for critical setup
- Use assert for test validations
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
assert.Nil(t, result)
```

## Package Organization
- Keep packages focused and cohesive
- Use internal for non-exported code
- Group related functionality
- Follow standard Go project layout

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

## References
- ADR-010: HTTP Error Handling Standard
- Go Code Review Comments: https://github.com/golang/go/wiki/CodeReviewComments
- Effective Go: https://golang.org/doc/effective_go.html
