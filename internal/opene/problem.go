package opene

// Problem represents errors that can be converted to RFC 7807 format
type Problem interface {
	Error
	AsProblem() *ProblemDetails
}

// ProblemDetails implements RFC 7807 for HTTP API errors
type ProblemDetails struct {
	Type     string      `json:"type"`               // URI reference
	Title    string      `json:"title"`              // Short, human-readable title
	Status   int         `json:"status"`             // HTTP status code
	Detail   string      `json:"detail,omitempty"`   // Human-readable explanation
	Instance string      `json:"instance,omitempty"` // URI reference to specific occurrence
	Extra    interface{} `json:"extra,omitempty"`    // Additional context
}

// Make BaseError implement Problem interface
func (e *BaseError) AsProblem() *ProblemDetails {
	return &ProblemDetails{
		Type:   uriForDomain(e.Domain),
		Title:  e.Message,
		Status: e.StatusCode,
		Detail: e.Error(),
		Extra:  e.Meta,
	}
}

// uriForDomain creates a type URI for the RFC 7807 response
func uriForDomain(domain string) string {
	return "https://openk.dev/errors/" + domain
}
