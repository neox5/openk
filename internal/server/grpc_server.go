package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/neox5/openk/internal/logging"
	"github.com/neox5/openk/internal/opene"
	"github.com/neox5/openk/internal/server/services/health"
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

	// Build all server options
	opts, err := buildServerOptions(cfg, logger)
	if err != nil {
		return nil, err
	}

	// Create server with options
	grpcServer := grpc.NewServer(opts...)

	s := &GRPCServer{
		config:   cfg,
		logger:   logger,
		server:   grpcServer,
		listener: listener,
	}

	// Register services
	if err := registerServices(ctx, s.server, logger); err != nil {
		listener.Close()
		return nil, err
	}

	return s, nil
}

// registerServices configures and registers all gRPC services
func registerServices(ctx context.Context, server *grpc.Server, logger *slog.Logger) error {
	// Register health service
	_, err := health.RegisterHealthServers(server, logger)
	if err != nil {
		return opene.NewInternalError("server", "register_health", "failed to register health service").
			Wrap(opene.AsError(err, "grpc", opene.CodeInternal))
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "registered gRPC services",
		slog.Bool("health_service", true),
	)

	return nil
}

// Start begins serving gRPC requests
func (s *GRPCServer) Start() error {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "starting gRPC server",
		slog.String("address", s.listener.Addr().String()),
	)

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
func (s *GRPCServer) Stop(ctx context.Context) error {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "stopping gRPC server",
		slog.String("address", s.listener.Addr().String()),
	)

	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return opene.NewInternalError("server", "stop_grpc", "graceful shutdown timed out").
			WithMetadata(opene.Metadata{
				"address": s.listener.Addr().String(),
			})
	case <-stopped:
		return nil
	}
}
