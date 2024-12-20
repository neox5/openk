package opene

import "net/http"

func NewValidationError(msg string) *Error {
	return &Error{
		Message:    msg,
		Code:       CodeValidationError,
		StatusCode: http.StatusBadRequest,
		Meta:       make(Metadata),
	}
}

func NewNotFoundError(msg string) *Error {
	return &Error{
		Message:    msg,
		Code:       CodeNotFound,
		StatusCode: http.StatusNotFound,
		Meta:       make(Metadata),
	}
}

func NewConflictError(msg string) *Error {
	return &Error{
		Message:    msg,
		Code:       CodeConflict,
		StatusCode: http.StatusConflict,
		Meta:       make(Metadata),
	}
}

func NewInternalError(msg string) *Error {
	return &Error{
		Message:     msg,
		Code:        CodeInternalError,
		StatusCode:  http.StatusInternalServerError,
		IsSensitive: true,
		Meta:        make(Metadata),
	}
}
