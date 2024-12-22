package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
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
		logger.LogAttrs(ctx, slog.LevelError, "invalid server configuration",
			slog.String("error", err.Error()),
		)
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
		return errors.New("TLS not yet implemented")
	}

	return s.server.ListenAndServe()
}
