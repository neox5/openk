package opene

import (
	"errors"
	"fmt"
)

// Metadata holds structured error context data
type Metadata map[string]interface{}

// Error is the core interface for OpenK errors defining behavior
type Error interface {
	error                           // Standard error interface
	Unwrap() error                 // Standard unwrap for error chains
	WithMetadata(md Metadata) Error // Adds structured context to error
	Wrap(err Error) Error          // Wraps another opene.Error
}

// BaseError implements the Error interface and provides the standard
// error implementation with additional context fields
type BaseError struct {
	Message     string   // Error message
	Code        string   // Error code (e.g., "E1001")
	Domain      string   // Error domain (e.g., "crypto", "auth")
	StatusCode  int      // HTTP status code
	IsSensitive bool     // Indicates if error contains sensitive data
	WrappedErr  error    // The wrapped error
	Meta        Metadata // Additional error context
}

func (e *BaseError) Error() string { return e.Message }

func (e *BaseError) Unwrap() error { return e.WrappedErr }

func (e *BaseError) WithMetadata(md Metadata) Error {
	e.Meta = md
	return e
}

func (e *BaseError) Wrap(err Error) Error {
	return &BaseError{
		Message:     fmt.Sprintf("%s: %v", e.Message, err),
		Code:        e.Code,
		Domain:      e.Domain,
		StatusCode:  e.StatusCode,
		IsSensitive: e.IsSensitive || err.(*BaseError).IsSensitive,
		WrappedErr:  err,
		Meta:        e.Meta,
	}
}

// AsError converts a standard error to an opene.Error.
// Returns nil if err is nil, the original error if it's already an Error,
// or creates a new Error if it's a standard error.
func AsError(domain string, code string, err error) Error {
	if err == nil {
		return nil
	}
	
	var openErr Error
	if errors.As(err, &openErr) {
		return openErr
	}

	return &BaseError{
		Message:     err.Error(),
		Code:        code,
		Domain:      domain,
		StatusCode:  500, // Default to internal server error
		IsSensitive: false,
		WrappedErr:  err,
		Meta:        make(Metadata),
	}
}
