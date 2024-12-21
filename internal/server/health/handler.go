package health

import (
	"encoding/json"
	"net/http"
	"time"
)

// Handler implements the health check endpoint
type Handler struct {
	startTime time.Time
}

// NewHandler creates a new health check handler
func NewHandler() *Handler {
	return &Handler{
		startTime: time.Now(),
	}
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Uptime:    time.Since(h.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
