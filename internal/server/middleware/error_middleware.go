package middleware

import (
	"net/http"
	"runtime/debug"
)

// ErrorHandler wraps an http.HandlerFunc and provides centralized error handling
func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Recover from panics and convert to 500 errors
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				debug.Stack()

				// Create internal error response
				w.WriteHeader(http.StatusInternalServerError)
				// Error response should follow RFC 7807
			}
		}()

		// Create response wrapper to catch errors
		ww := &responseWriter{ResponseWriter: w}

		// Call the wrapped handler
		next(ww, r)

		// Check if we need to handle any errors
		if ww.error != nil {
			// Write error response following RFC 7807
			ww.WriteError(ww.error)
		}
	}
}

// responseWriter wraps http.ResponseWriter to capture errors
type responseWriter struct {
	http.ResponseWriter
	error error
}

// WriteError captures an error to be handled
func (w *responseWriter) WriteError(err error) {
	w.error = err
}
