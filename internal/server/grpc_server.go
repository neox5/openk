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
func NewGRPCServer(ctx context.Context, cfg *Config, logger *slog.Logger, opts ...ServerOption) (*GRPCServer, error) {
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

	// Process options
	options := defaultServerOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Apply config-based keepalive settings
	options.keepaliveParams.MaxConnectionAge = cfg.MaxConnectionAge
	options.keepaliveParams.MaxConnectionIdle = cfg.MaxConnectionIdle
	options.keepalivePolicy.MinTime = cfg.MinConnectionTime
	options.keepalivePolicy.PermitWithoutStream = cfg.PermitWithoutStream

	// Create server options
	serverOpts := []grpc.ServerOption{
		grpc.KeepaliveParams(options.keepaliveParams),
		grpc.KeepaliveEnforcementPolicy(options.keepalivePolicy),
	}

	if len(options.unaryInterceptors) > 0 {
		serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(options.unaryInterceptors...))
	}
	if len(options.streamInterceptors) > 0 {
		serverOpts = append(serverOpts, grpc.ChainStreamInterceptor(options.streamInterceptors...))
	}

	// Create server
	server := grpc.NewServer(serverOpts...)

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
