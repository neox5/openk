package opene

import "net/http"

// NewValidationError creates an error for invalid input
func NewValidationError(msg string) *BaseError {
	return &BaseError{
		msg:       msg,
		errType:   "validation",
		status:    http.StatusBadRequest,
		sensitive: false,
		metadata:  make(map[string]interface{}),
	}
}

// NewNotFoundError creates an error for missing resources
func NewNotFoundError(msg string) *BaseError {
	return &BaseError{
		msg:       msg,
		errType:   "not_found",
		status:    http.StatusNotFound,
		sensitive: false,
		metadata:  make(map[string]interface{}),
	}
}

// NewUnauthorizedError creates an error for authentication failures
func NewUnauthorizedError(msg string) *BaseError {
	return &BaseError{
		msg:       msg,
		errType:   "unauthorized",
		status:    http.StatusUnauthorized,
		sensitive: false,
		metadata:  make(map[string]interface{}),
	}
}

// NewForbiddenError creates an error for authorization failures
func NewForbiddenError(msg string) *BaseError {
	return &BaseError{
		msg:       msg,
		errType:   "forbidden",
		status:    http.StatusForbidden,
		sensitive: false,
		metadata:  make(map[string]interface{}),
	}
}

// NewConflictError creates an error for resource conflicts
func NewConflictError(msg string) *BaseError {
	return &BaseError{
		msg:       msg,
		errType:   "conflict",
		status:    http.StatusConflict,
		sensitive: false,
		metadata:  make(map[string]interface{}),
	}
}

// NewInternalError creates an error for system failures
func NewInternalError(msg string) *BaseError {
	return &BaseError{
		msg:       msg,
		errType:   "internal",
		status:    http.StatusInternalServerError,
		sensitive: true,
		metadata:  make(map[string]interface{}),
	}
}
