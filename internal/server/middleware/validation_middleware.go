package middleware

import "net/http"

const (
	// DefaultMaxBodySize is 1MB
	DefaultMaxBodySize = 1 << 20
)

// MaxBodySize limits the size of request bodies
func MaxBodySize(maxBytes int64) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Only apply to requests that might have bodies
			if r.ContentLength > 0 || r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next(w, r)
		}
	}
}
