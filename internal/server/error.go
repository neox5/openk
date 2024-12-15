package server

import (
	"encoding/json"
	"net/http"
)

type Status int

const (
	StatusInvalidRequest Status = 400
	StatusUnauthorized   Status = 401
	StatusNotFound       Status = 404
	StatusConflict       Status = 409
	StatusInternal       Status = 500
)

type Error struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   Status      `json:"status"`
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

func ValidationError(field, reason string) *Error {
	return &Error{
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

func NotFoundError(resource, id string) *Error {
	return &Error{
		Type:   typeNotFound,
		Title:  "Resource Not Found",
		Status: StatusNotFound,
		Detail: resource + " with id " + id + " not found",
	}
}

func InternalError() *Error {
	return &Error{
		Type:   typeInternal,
		Title:  "Internal Error",
		Status: StatusInternal,
	}
}

// Write sends an error as an RFC 7807 response
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var problemDetails *Error
	if e, ok := err.(*Error); ok {
		problemDetails = e
	} else {
		problemDetails = InternalError().WithDetail(err.Error())
	}

	if problemDetails.Instance == "" {
		problemDetails.Instance = r.URL.Path
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(int(problemDetails.Status))
	json.NewEncoder(w).Encode(problemDetails)
}

func (e *Error) WithInstance(uri string) *Error {
	e.Instance = uri
	return e
}

func (e *Error) WithDetail(detail string) *Error {
	e.Detail = detail
	return e
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return e.Detail
	}
	return e.Title
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
