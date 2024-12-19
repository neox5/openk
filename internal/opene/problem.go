package opene

// Problem represents RFC 7807 convertible errors
type Problem interface {
	error
	ToProblem() *ProblemDetails
}

// ProblemDetails implements RFC 7807
type ProblemDetails struct {
	Type     string      `json:"type"`               // URI reference
	Title    string      `json:"title"`              // Short summary
	Status   int         `json:"status"`             // HTTP status code
	Detail   string      `json:"detail,omitempty"`   // Detailed explanation
	Instance string      `json:"instance,omitempty"` // URI reference
	Extra    interface{} `json:"extra,omitempty"`    // Additional context
}
