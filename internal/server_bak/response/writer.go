package response

import (
	"encoding/json"
	"net/http"

	"github.com/neox5/openk/internal/errors"
)

// WriteError writes an RFC 7807 error response
func WriteError(w http.ResponseWriter, r *http.Request, err *errors.HTTPError) {
	if err.Instance == "" {
		err.Instance = r.URL.Path
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
