# Error Handling System Implementation Plan

## 1. Core Error Types (2-3 days)

### 1.1 Base Error Interface
```go
// OpenError defines the base interface for all errors
type OpenError interface {
    error
    Code() string              // Returns error code (e.g., "E1001")
    Domain() string            // Returns error domain (e.g., "crypto", "auth")
    HTTPStatus() int           // Maps to HTTP status if exposed via API
    Unwrap() error            // Standard unwrap for error chains
    WithMetadata(md Metadata)  // Adds context to error
    IsSensitive() bool        // Indicates if error contains sensitive data
}

// Metadata holds structured error context
type Metadata map[string]interface{}
```

### 1.2 Implementation Tasks
- [ ] Create base error struct
- [ ] Implement OpenError interface
- [ ] Add error wrapping support
- [ ] Add metadata handling
- [ ] Add sensitivity marking
- [ ] Create error builder pattern
- [ ] Add comprehensive tests

### 1.3 Error Categories
```go
// Define standard error categories
const (
    // Core categories
    CategoryValidation = "validation"
    CategoryAuth      = "auth"
    CategoryCrypto    = "crypto"
    CategoryStorage   = "storage"
    CategoryInternal  = "internal"
    
    // HTTP status mapping
    StatusValidation  = http.StatusBadRequest
    StatusAuth       = http.StatusUnauthorized
    StatusNotFound   = http.StatusNotFound
    StatusInternal   = http.StatusInternalServerError
)
```

## 2. Error Creation System (2-3 days)

### 2.1 Error Factory
```go
// ErrorFactory provides methods to create domain-specific errors
type ErrorFactory interface {
    // Core creation methods
    New(code string, msg string) OpenError
    Wrap(code string, err error, msg string) OpenError
    
    // Specialized creators
    NotFound(resource, id string) OpenError
    Validation(field, reason string) OpenError
    Unauthorized(reason string) OpenError
    Internal(err error) OpenError
}
```

### 2.2 Implementation Tasks
- [ ] Create error factory interface
- [ ] Implement default factory
- [ ] Add domain-specific factories
- [ ] Create helper functions
- [ ] Add factory tests
- [ ] Document factory usage

## 3. HTTP Integration (2-3 days)

### 3.1 RFC 7807 Problem Details
```go
// ProblemDetails implements RFC 7807
type ProblemDetails struct {
    Type     string      `json:"type"`
    Title    string      `json:"title"`
    Status   int         `json:"status"`
    Detail   string      `json:"detail,omitempty"`
    Instance string      `json:"instance,omitempty"`
    Meta     Metadata    `json:"meta,omitempty"`
}

// ToProblem converts OpenError to ProblemDetails
func ToProblem(err OpenError, req *http.Request) *ProblemDetails {
    return &ProblemDetails{
        Type:     fmt.Sprintf("https://openk.dev/errors/%s", err.Domain()),
        Title:    err.Error(),
        Status:   err.HTTPStatus(),
        Instance: req.URL.Path,
        Meta:     sanitizeMetadata(err),
    }
}
```

### 3.2 Implementation Tasks
- [ ] Create problem details type
- [ ] Implement error conversion
- [ ] Add metadata sanitization
- [ ] Create middleware
- [ ] Add response writers
- [ ] Create integration tests

## 4. Logging Integration (1-2 days)

### 4.1 Error Context for Logging
```go
// LogContext extracts structured data for logging
type LogContext struct {
    Code     string                 `json:"code"`
    Domain   string                 `json:"domain"`
    Message  string                 `json:"message"`
    Stack    string                 `json:"stack,omitempty"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func ToLogContext(err OpenError) *LogContext {
    return &LogContext{
        Code:     err.Code(),
        Domain:   err.Domain(),
        Message:  err.Error(),
        Stack:    getStackTrace(err),
        Metadata: sanitizeForLogging(err.Metadata()),
    }
}
```

### 4.2 Implementation Tasks
- [ ] Create log context type
- [ ] Implement context extraction
- [ ] Add stack trace support
- [ ] Create logging helpers
- [ ] Add logging middleware
- [ ] Create logging tests

## 5. Migration Strategy (2-3 days)

### 5.1 Package Updates
- [ ] Crypto Package
  - [ ] Update error definitions
  - [ ] Implement error factory
  - [ ] Update error handling
  - [ ] Update tests

- [ ] KMS Package
  - [ ] Update error definitions
  - [ ] Implement error factory
  - [ ] Update error handling
  - [ ] Update tests

- [ ] Storage Package
  - [ ] Update error definitions
  - [ ] Implement error factory
  - [ ] Update error handling
  - [ ] Update tests

### 5.2 Implementation Tasks
- [ ] Create migration guide
- [ ] Update existing errors
- [ ] Add new error handling
- [ ] Update tests
- [ ] Validate changes

## 6. Documentation (1-2 days)

### 6.1 Documentation Tasks
- [ ] Create error handling guide
- [ ] Document error codes
- [ ] Add usage examples
- [ ] Create troubleshooting guide
- [ ] Add API error documentation

### 6.2 Example Documentation
```markdown
## Error Handling Guide

### Creating Errors
```go
// Using factory
err := factory.New("E1001", "invalid input")
err := factory.Wrap("E1002", originalErr, "operation failed")

// Using helpers
err := factory.NotFound("user", "123")
err := factory.Validation("email", "invalid format")
```

### Error Handling
```go
func HandleOperation() error {
    err := doSomething()
    if err != nil {
        return factory.Wrap("E1003", err, "operation failed")
    }
    return nil
}
```
```

## 7. Testing Framework (1-2 days)

### 7.1 Test Utilities
```go
// ErrorAssertions provides testing helpers
type ErrorAssertions interface {
    AssertErrorCode(t *testing.T, err error, code string)
    AssertErrorDomain(t *testing.T, err error, domain string)
    AssertErrorMetadata(t *testing.T, err error, key string, value interface{})
    AssertProblemResponse(t *testing.T, resp *http.Response, code string)
}
```

### 7.2 Implementation Tasks
- [ ] Create test utilities
- [ ] Add assertion helpers
- [ ] Create mock factory
- [ ] Add example tests
- [ ] Create testing guide

## Timeline

Week 1:
- Days 1-3: Core Error Types
- Days 4-5: Error Creation System

Week 2:
- Days 1-2: HTTP Integration
- Days 3-4: Logging Integration
- Day 5: Testing Framework

Week 3:
- Days 1-3: Migration Strategy
- Days 4-5: Documentation and Review
