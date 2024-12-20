package opene

// ErrorCode represents a standardized error code
type ErrorCode string

const (
	// Basic error types that map directly to URLs
	CodeValidation ErrorCode = "validation"
	CodeNotFound   ErrorCode = "not_found"
	CodeConflict   ErrorCode = "conflict"
	CodeInternal   ErrorCode = "internal"
)
