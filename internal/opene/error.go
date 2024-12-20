package opene

// Metadata holds structured error context
type Metadata map[string]interface{}

// Error provides extended error information
type Error struct {
    Message     string
    Code        ErrorCode
    Domain      string
    StatusCode  int
    IsSensitive bool
    WrappedErr  error
    Meta        Metadata
}

// Error implements the error interface
func (e *Error) Error() string {
    return e.Message
}

// Unwrap implements error unwrapping
func (e *Error) Unwrap() error {
    return e.WrappedErr
}

// UnwrapAll returns the root cause error by unwrapping through Error types
func (e *Error) UnwrapAll() error {
    if e.WrappedErr == nil {
        return nil
    }
    
    if wrappedErr, ok := e.WrappedErr.(*Error); ok {
        return wrappedErr.UnwrapAll()
    }
    return e.WrappedErr
}

// WithMetadata adds metadata to the error
func (e *Error) WithMetadata(md Metadata) *Error {
    e.Meta = md
    return e
}

// Wrap wraps another Error with additional context
func (e *Error) Wrap(err *Error) *Error {
    if err == nil {
        return e
    }

    return &Error{
        Message:     e.Message + ": " + err.Message,
        Code:        e.Code,
        Domain:      e.Domain,
        StatusCode:  e.StatusCode,
        IsSensitive: e.IsSensitive || err.IsSensitive,
        WrappedErr:  err,
        Meta:        e.Meta,
    }
}

// Sensitive marks the error as containing sensitive information
func (e *Error) Sensitive() *Error {
    e.IsSensitive = true
    return e
}

// AsError converts a standard error to an Error
func AsError(err error, domain string, code ErrorCode) *Error {
    if err == nil {
        return nil
    }
    
    if e, ok := err.(*Error); ok {
        return e
    }

    return &Error{
        Message:     err.Error(),
        Code:        code,
        Domain:      domain,
        StatusCode:  500,
        IsSensitive: false,
        WrappedErr:  err,
        Meta:        make(Metadata),
    }
}
