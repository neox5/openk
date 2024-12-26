package server

import (
	"log/slog"
	"time"

	"github.com/neox5/openk/internal/server/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// buildServerOptions combines all gRPC server options
func buildServerOptions(cfg *Config, logger *slog.Logger) ([]grpc.ServerOption, error) {
	opts := []grpc.ServerOption{}

	// Add connection options
	connectionOpts, err := buildConnectionOptions(cfg)
	if err != nil {
		return nil, err
	}
	opts = append(opts, connectionOpts...)

	// Add interceptor options
	interceptorOpts, err := buildInterceptorOptions(logger)
	if err != nil {
		return nil, err
	}
	opts = append(opts, interceptorOpts...)

	// Add transport options
	transportOpts, err := buildTransportOptions(cfg)
	if err != nil {
		return nil, err
	}
	opts = append(opts, transportOpts...)

	return opts, nil
}

// buildConnectionOptions configures how connections are managed
func buildConnectionOptions(cfg *Config) ([]grpc.ServerOption, error) {
	kaParams := keepalive.ServerParameters{
		MaxConnectionIdle:     cfg.MaxConnectionIdle,
		MaxConnectionAge:      cfg.MaxConnectionAge,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  5 * time.Second,
		Timeout:               1 * time.Second,
	}

	kaPolicy := keepalive.EnforcementPolicy{
		MinTime:             cfg.MinConnectionTime,
		PermitWithoutStream: cfg.PermitWithoutStream,
	}

	return []grpc.ServerOption{
		grpc.KeepaliveParams(kaParams),
		grpc.KeepaliveEnforcementPolicy(kaPolicy),
	}, nil
}

// buildInterceptorOptions configures server middleware
func buildInterceptorOptions(logger *slog.Logger) ([]grpc.ServerOption, error) {
	// Unary (request/response) interceptors
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLogging(logger),
		// Future interceptors:
		// - Authentication
		// - Request validation
		// - Metrics collection
		// - Panic recovery
		// - Rate limiting
	}

	// Streaming interceptors
	streamInterceptors := []grpc.StreamServerInterceptor{
		interceptors.StreamLogging(logger),
		// Future interceptors:
		// - Authentication
		// - Metrics collection
		// - Panic recovery
		// - Rate limiting
	}

	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}, nil
}

// buildTransportOptions configures transport-level settings
func buildTransportOptions(cfg *Config) ([]grpc.ServerOption, error) {
	opts := []grpc.ServerOption{}

	// Future transport options:
	// - TLS configuration
	// - Maximum message sizes
	// - Compression settings

	return opts, nil
}
