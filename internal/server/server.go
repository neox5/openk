package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/neox5/openk/internal/logging"
	"github.com/neox5/openk/internal/opene"
)

type Server struct {
	config *Config
	logger *slog.Logger
	server *http.Server
	mux    *http.ServeMux
}

func NewServer(ctx context.Context, cfg *Config, logger *slog.Logger) (*Server, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if err := cfg.Validate(); err != nil {
		logging.LogError(ctx, logger, "invalid server configuration", err)
		return nil, err // Error is already openE type from Validate()
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
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	s.registerRoutes()
	return s, nil
}

func (s *Server) Start() error {
	s.logger.LogAttrs(s.server.BaseContext(nil), slog.LevelInfo, "starting server",
		slog.String("address", s.server.Addr),
		slog.Duration("read_timeout", s.config.ReadTimeout),
		slog.Duration("write_timeout", s.config.WriteTimeout),
	)

	if s.config.EnableTLS {
		return opene.NewInternalError("server", "start", "TLS support not implemented")
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return opene.NewInternalError("server", "start", "server startup failed").
			WithMetadata(opene.Metadata{
				"address": s.server.Addr,
			}).
			Wrap(opene.AsError(err, "http", opene.CodeInternal))
	}

	return nil
}
