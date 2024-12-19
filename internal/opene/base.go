package opene

// BaseError provides the foundation for all errors in the system
type BaseError struct {
	msg       string
	cause     error
	sensitive bool
	errType   string
	status    int
	metadata  map[string]interface{}
}

// Error returns the error message
func (e *BaseError) Error() string { return e.msg }

// Unwrap returns the wrapped error
func (e *BaseError) Unwrap() error { return e.cause }

// ToProblem converts the error to RFC 7807 format
func (e *BaseError) ToProblem() *ProblemDetails {
	if e.sensitive {
		return &ProblemDetails{
			Type:   "https://openk.dev/errors/" + e.errType,
			Title:  "Operation Failed",
			Status: e.status,
		}
	}

	return &ProblemDetails{
		Type:   "https://openk.dev/errors/" + e.errType,
		Title:  e.msg,
		Status: e.status,
		Extra:  e.metadata,
	}
}
