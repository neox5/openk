package opene

import (
	"errors"
	"strings"
)

var errorBaseURI = "https://example.com/errors/"

// SetErrorBaseURI configures the base URI for error types.
// Must be called before creating any errors.
func SetErrorBaseURI(uri string) {
	// Ensure trailing slash
	if !strings.HasSuffix(uri, "/") {
		uri += "/"
	}
	errorBaseURI = uri
}

// Problem implements RFC 7807 for HTTP API errors
type Problem struct {
	Type     string      `json:"type"`               // URI reference
	Title    string      `json:"title"`              // Short, human-readable title
	Status   int         `json:"status"`             // HTTP status code
	Detail   string      `json:"detail,omitempty"`   // Human-readable explanation
	Instance string      `json:"instance,omitempty"` // URI reference to specific occurrence
	Extra    interface{} `json:"extra,omitempty"`    // Additional context
}

// AsProblem converts any error into a Problem.
// If the error is not an Error type, returns a generic internal error Problem.
func AsProblem(err error) *Problem {
	var e *Error
	if !errors.As(err, &e) {
		return &Problem{
			Type:   uriForDomain("internal"),
			Title:  "Internal Server Error",
			Status: 500,
		}
	}

	// For sensitive errors, return a generic 500 error
	if e.IsSensitive {
		return &Problem{
			Type:   uriForDomain("internal"),
			Title:  "Internal Server Error",
			Status: 500,
		}
	}

	return &Problem{
		Type:   uriForDomain(e.Domain),
		Title:  e.Message,
		Status: e.StatusCode,
		Detail: e.Error(),
		Extra:  e.Meta,
	}
}

func uriForDomain(domain string) string {
	return errorBaseURI + domain
}
