package opene

import "net/http"

func NewValidationError(domain, operation, msg string) *Error {
	return &Error{
		Message:    msg,
		Code:       CodeValidation,
		Domain:     domain,
		Operation:  operation,
		StatusCode: http.StatusBadRequest,
		Meta:       make(Metadata),
	}
}

func NewNotFoundError(domain, operation, msg string) *Error {
	return &Error{
		Message:    msg,
		Code:       CodeNotFound,
		Domain:     domain,
		Operation:  operation,
		StatusCode: http.StatusNotFound,
		Meta:       make(Metadata),
	}
}

func NewConflictError(domain, operation, msg string) *Error {
	return &Error{
		Message:    msg,
		Code:       CodeConflict,
		Domain:     domain,
		Operation:  operation,
		StatusCode: http.StatusConflict,
		Meta:       make(Metadata),
	}
}

func NewInternalError(domain, operation, msg string) *Error {
	return &Error{
		Message:     msg,
		Code:        CodeInternal,
		Domain:      domain,
		Operation:   operation,
		StatusCode:  http.StatusInternalServerError,
		IsSensitive: true,
		Meta:        make(Metadata),
	}
}
