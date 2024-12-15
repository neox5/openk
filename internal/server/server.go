package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/neox5/openk/internal/server/middleware"
	"github.com/neox5/openk/internal/storage"
)

// Config holds server configuration
type Config struct {
	// Server address in format "host:port"
	Addr string

	// Read timeout for entire request
	ReadTimeout time.Duration

	// Write timeout for response
	WriteTimeout time.Duration

	// Idle timeout for keepalive connections
	IdleTimeout time.Duration

	// Maximum header size
	MaxHeaderBytes int

	// Maximum body size
	MaxBodySize int64
}

// Server represents the HTTP server instance
type Server struct {
	config  Config
	storage storage.MiniStorageBackend
	server  *http.Server
	mux     *http.ServeMux
}

// DefaultConfig returns server configuration with sensible defaults
func DefaultConfig() Config {
	return Config{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
		MaxBodySize:    1 << 20, // 1MB
	}
}

// New creates a new server instance
func New(config Config, storage storage.MiniStorageBackend) (*Server, error) {
	if storage == nil {
		return nil, fmt.Errorf("storage backend cannot be nil")
	}

	s := &Server{
		config:  config,
		storage: storage,
		mux:     http.NewServeMux(),
	}

	// Initialize HTTP server
	s.server = &http.Server{
		Addr:           config.Addr,
		Handler:        s.mux,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		IdleTimeout:    config.IdleTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	// Register routes
	s.registerRoutes()

	return s, nil
}

// Start begins listening for requests
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// handleHealth returns the health check handler
func (s *Server) handleHealth() http.HandlerFunc {
	type healthResponse struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		resp := healthResponse{
			Status:  "ok",
			Version: "0.1.0", // TODO: Get from version package
		}

		WriteJSON(w, http.StatusOK, resp)
	}
}
