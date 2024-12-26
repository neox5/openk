package server

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// ServerOption allows for customizing the gRPC server
type ServerOption func(*serverOptions)

type serverOptions struct {
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	keepaliveParams    keepalive.ServerParameters
	keepalivePolicy    keepalive.EnforcementPolicy
}

// defaultServerOptions returns the default server options
func defaultServerOptions() *serverOptions {
	return &serverOptions{
		keepaliveParams: keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Minute,
			MaxConnectionAge:      30 * time.Minute,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  5 * time.Second,
			Timeout:              1 * time.Second,
		},
		keepalivePolicy: keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		},
	}
}

// WithUnaryInterceptors adds unary interceptors to the server
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *serverOptions) {
		o.unaryInterceptors = append(o.unaryInterceptors, interceptors...)
	}
}

// WithStreamInterceptors adds stream interceptors to the server
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(o *serverOptions) {
		o.streamInterceptors = append(o.streamInterceptors, interceptors...)
	}
}

// WithKeepaliveParams customizes the keepalive parameters
func WithKeepaliveParams(params keepalive.ServerParameters) ServerOption {
	return func(o *serverOptions) {
		o.keepaliveParams = params
	}
}

// WithKeepalivePolicy customizes the keepalive enforcement policy
func WithKeepalivePolicy(policy keepalive.EnforcementPolicy) ServerOption {
	return func(o *serverOptions) {
		o.keepalivePolicy = policy
	}
}
