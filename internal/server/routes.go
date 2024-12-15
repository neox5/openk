package server

import (
	"net/http"

	"github.com/neox5/openk/internal/server/middleware"
)

// Registers all routes for the server
func (s *Server) registerRoutes() {
	// Apply base middleware stack
	withMiddleware := func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.ErrorHandler(
			middleware.MaxBodySize(s.config.MaxBodySize)(handler),
		)
	}

	// Health check - simple GET endpoint
	s.mux.HandleFunc("GET /health", withMiddleware(s.handleHealth()))

	// Derivation endpoints
	derivationHandler := newDerivationV1Handler(s.storage)

	// POST /api/v1/derivation/params
	s.mux.HandleFunc("POST /api/v1/derivation/params",
		withMiddleware(derivationHandler.storeParams))

	// GET /api/v1/derivation/params/{username}
	s.mux.HandleFunc("GET /api/v1/derivation/params/{username}",
		withMiddleware(derivationHandler.getParams))
}
