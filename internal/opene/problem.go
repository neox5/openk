package opene

import (
	"errors"
	"strings"
)

const defaultBaseURI = "https://openk.dev/errors"

var errorBaseURI = defaultBaseURI

// Problem implements RFC 7807 for HTTP API errors
type Problem struct {
	Type     string      `json:"type"`               // URI reference that identifies the error type
	Title    string      `json:"title"`              // Short, human-readable summary
	Status   int         `json:"status"`             // HTTP status code
	Detail   string      `json:"detail,omitempty"`   // Human-readable explanation
	Instance string      `json:"instance,omitempty"` // URI reference to specific occurrence
	Extra    interface{} `json:"extra,omitempty"`    // Additional context
}

// SetErrorBaseURI configures the base URI for error types.
// Must include protocol (e.g., "https://") and should not end with a slash.
func SetErrorBaseURI(uri string) *Error {
	if uri == "" {
		return NewValidationError("base URI cannot be empty").
			WithDomain("opene").
			WithOperation("set_base_uri").
			WithMetadata(Metadata{
				"provided_uri": uri,
			})
	}
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		return NewValidationError("base URI must include protocol (http:// or https://)").
			WithDomain("opene").
			WithOperation("set_base_uri").
			WithMetadata(Metadata{
				"provided_uri": uri,
			})
	}

	// Remove trailing slash if present
	errorBaseURI = strings.TrimSuffix(uri, "/")
	return nil
}

// ResetErrorBaseURI resets the base URI to the default value
func ResetErrorBaseURI() {
	errorBaseURI = defaultBaseURI
}

// AsProblem converts an Error into a Problem (RFC 7807).
// If err is not an Error type or is sensitive, returns a generic problem.
func AsProblem(err error) *Problem {
	// Set default problem for non-Error types and sensitive errors
	var e *Error
	if !errors.As(err, &e) || e.IsSensitive {
		return &Problem{
			Type:   errorBaseURI + "/internal",
			Title:  "Internal Server Error",
			Status: 500,
		}
	}

	// Convert to Problem using error details
	return &Problem{
		Type:   errorBaseURI + "/" + string(e.Code),
		Title:  e.Message,
		Status: e.StatusCode,
		Detail: e.Error(),
		Extra:  e.Meta,
	}
}
