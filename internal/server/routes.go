package server

import (
	"github.com/neox5/openk/internal/server/health"
	"github.com/neox5/openk/internal/server/middleware"
)

// registerRoutes sets up all HTTP routes
func (s *Server) registerRoutes() {
	// Create handlers
	healthHandler := health.NewHandler()

	// Add middleware
	loggingMiddleware := middleware.NewLogging(s.logger)

	// Register routes
	s.mux.Handle("GET /health", loggingMiddleware(healthHandler))
}
