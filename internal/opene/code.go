package opene

// ErrorCode represents a standardized error code
type ErrorCode string

const (
	CodeValidationError ErrorCode = "VALIDATION_ERROR"
	CodeNotFound        ErrorCode = "NOT_FOUND"
	CodeConflict        ErrorCode = "CONFLICT"
	CodeInternalError   ErrorCode = "INTERNAL_ERROR"
)
