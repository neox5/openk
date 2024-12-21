package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// Server represents the HTTP server
type Server struct {
	config *Config
	logger *slog.Logger
	server *http.Server
	mux    *http.ServeMux
}

// NewServer creates a new server instance
func NewServer(cfg *Config, logger *slog.Logger) (*Server, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	if logger == nil {
		logger = slog.Default()
	}

	s := &Server{
		config: cfg,
		logger: logger,
		mux:    http.NewServeMux(),
	}

	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      s.mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Register routes
	s.registerRoutes()

	return s, nil
}

// Start begins listening for requests
func (s *Server) Start() error {
	s.logger.Info("starting server",
		"address", s.server.Addr,
		"read_timeout", s.config.ReadTimeout,
		"write_timeout", s.config.WriteTimeout,
	)

	if s.config.EnableTLS {
		// TODO: Implement TLS configuration
		return errors.New("TLS not yet implemented")
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")

	// Create a timeout context for shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, s.config.ShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
