package httperror

import (
	"encoding/json"
	"net/http"
)

type HTTPStatus int

const (
	StatusInvalidRequest HTTPStatus = http.StatusBadRequest
	StatusUnauthorized   HTTPStatus = http.StatusUnauthorized
	StatusNotFound       HTTPStatus = http.StatusNotFound
	StatusConflict       HTTPStatus = http.StatusConflict
	StatusInternal       HTTPStatus = http.StatusInternalServerError
)

// HTTPError implements RFC 7807 for HTTP API error responses
type HTTPError struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   HTTPStatus  `json:"status"`
	Detail   string      `json:"detail,omitempty"`
	Instance string      `json:"instance,omitempty"`
	Extra    interface{} `json:"extra,omitempty"`
}

const (
	typePrefix     = "https://openk.dev/errors"
	typeValidation = typePrefix + "/validation"
	typeAuth       = typePrefix + "/auth"
	typeNotFound   = typePrefix + "/not-found"
	typeInternal   = typePrefix + "/internal"
	typeConflict   = typePrefix + "/conflict"
)

// ValidationError creates an error response for request validation failures
func ValidationError(field, reason string) *HTTPError {
	return &HTTPError{
		Type:   typeValidation,
		Title:  "Validation Failed",
		Status: StatusInvalidRequest,
		Extra: map[string]interface{}{
			"validation": []map[string]string{{
				"field":  field,
				"reason": reason,
			}},
		},
	}
}

// NotFoundError creates an error response when a requested resource doesn't exist
func NotFoundError(resource, id string) *HTTPError {
	return &HTTPError{
		Type:   typeNotFound,
		Title:  "Resource Not Found",
		Status: StatusNotFound,
		Detail: resource + " with id " + id + " not found",
	}
}

// InternalError creates an error response for unexpected server-side errors
func InternalError() *HTTPError {
	return &HTTPError{
		Type:   typeInternal,
		Title:  "Internal Error",
		Status: StatusInternal,
	}
}

// WriteError writes any error as an RFC 7807 response, converting non-HTTPErrors to internal errors
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var problemDetails *HTTPError
	if e, ok := err.(*HTTPError); ok {
		problemDetails = e
	} else {
		// Never expose internal error details to clients
		problemDetails = InternalError()
	}

	if problemDetails.Instance == "" {
		problemDetails.Instance = r.URL.Path
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(int(problemDetails.Status))
	json.NewEncoder(w).Encode(problemDetails)
}

// WithInstance adds the request URI to the error response
func (e *HTTPError) WithInstance(uri string) *HTTPError {
	e.Instance = uri
	return e
}

// WithDetail adds a detailed error message to the response
func (e *HTTPError) WithDetail(detail string) *HTTPError {
	e.Detail = detail
	return e
}

// Error implements the error interface for HTTPError
func (e *HTTPError) Error() string {
	if e.Detail != "" {
		return e.Detail
	}
	return e.Title
}

// WriteJSON writes a JSON response with the given status code and data
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
