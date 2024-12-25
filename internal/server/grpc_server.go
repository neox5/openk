package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/neox5/openk/internal/logging"
	"github.com/neox5/openk/internal/opene"
	"google.golang.org/grpc"
)

// GRPCServer handles the gRPC server implementation
type GRPCServer struct {
	config   *Config
	logger   *slog.Logger
	server   *grpc.Server
	listener net.Listener
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(ctx context.Context, cfg *Config, logger *slog.Logger) (*GRPCServer, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if err := cfg.Validate(); err != nil {
		logging.LogError(ctx, logger, "invalid server configuration", err)
		return nil, err
	}

	if logger == nil {
		logger = slog.Default()
	}

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPCPort))
	if err != nil {
		return nil, opene.NewInternalError("server", "create_grpc_server", "failed to create listener").
			WithMetadata(opene.Metadata{
				"host": cfg.Host,
				"port": cfg.GRPCPort,
			}).
			Wrap(opene.AsError(err, "net", opene.CodeInternal))
	}

	// Create server options
	opts := []grpc.ServerOption{
		// TODO: Add interceptors
		// TODO: Add TLS configuration
	}

	// Create server
	server := grpc.NewServer(opts...)

	return &GRPCServer{
		config:   cfg,
		logger:   logger,
		server:   server,
		listener: listener,
	}, nil
}

// Start begins serving gRPC requests
func (s *GRPCServer) Start() error {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "starting gRPC server",
		slog.String("address", s.listener.Addr().String()),
	)

	// Start serving
	if err := s.server.Serve(s.listener); err != nil {
		return opene.NewInternalError("server", "start_grpc", "server startup failed").
			WithMetadata(opene.Metadata{
				"address": s.listener.Addr().String(),
			}).
			Wrap(opene.AsError(err, "grpc", opene.CodeInternal))
	}

	return nil
}

// Stop gracefully shuts down the server
func (s *GRPCServer) Stop() error {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "stopping gRPC server",
		slog.String("address", s.listener.Addr().String()),
	)

	s.server.GracefulStop()
	return nil
}
