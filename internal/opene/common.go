package opene

import "net/http"

// NewValidationError creates a generic validation error
func NewValidationError(msg string) *BaseError {
	return &BaseError{
		Message:    msg,
		Code:       "VALIDATION_ERROR",
		Domain:     "validation",
		StatusCode: http.StatusBadRequest,
		Meta:       make(Metadata),
	}
}

// NewNotFoundError creates a generic not found error
func NewNotFoundError(msg string) *BaseError {
	return &BaseError{
		Message:    msg,
		Code:       "NOT_FOUND",
		Domain:     "resource",
		StatusCode: http.StatusNotFound,
		Meta:       make(Metadata),
	}
}

// NewConflictError creates a generic conflict error
func NewConflictError(msg string) *BaseError {
	return &BaseError{
		Message:    msg,
		Code:       "CONFLICT",
		Domain:     "resource",
		StatusCode: http.StatusConflict,
		Meta:       make(Metadata),
	}
}

// NewInternalError creates a generic internal error
func NewInternalError(msg string) *BaseError {
	return &BaseError{
		Message:     msg,
		Code:        "INTERNAL_ERROR",
		Domain:      "internal",
		StatusCode:  http.StatusInternalServerError,
		IsSensitive: true,
		Meta:        make(Metadata),
	}
}
