package interceptors

import (
	"context"
	"log/slog"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// UnaryLogging returns a UnaryServerInterceptor that logs requests and responses
func UnaryLogging(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Extract method name for logging
		method := path.Base(info.FullMethod)

		// Get peer information if available
		var peerAddr string
		if p, ok := peer.FromContext(ctx); ok {
			peerAddr = p.Addr.String()
		}

		// Log request
		logger.LogAttrs(ctx, slog.LevelDebug, "handling gRPC request",
			slog.String("method", method),
			slog.String("peer_addr", peerAddr),
		)

		// Handle request
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Log response
		attrs := []slog.Attr{
			slog.String("method", method),
			slog.String("peer_addr", peerAddr),
			slog.Duration("duration", duration),
		}

		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
			logger.LogAttrs(ctx, slog.LevelError, "gRPC request failed", attrs...)
		} else {
			logger.LogAttrs(ctx, slog.LevelInfo, "gRPC request completed", attrs...)
		}

		return resp, err
	}
}

// StreamLogging returns a StreamServerInterceptor that logs streaming operations
func StreamLogging(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Extract method name
		method := path.Base(info.FullMethod)

		// Get peer information
		var peerAddr string
		if p, ok := peer.FromContext(ss.Context()); ok {
			peerAddr = p.Addr.String()
		}

		// Log stream start
		logger.LogAttrs(ss.Context(), slog.LevelDebug, "starting gRPC stream",
			slog.String("method", method),
			slog.String("peer_addr", peerAddr),
			slog.Bool("is_client_stream", info.IsClientStream),
			slog.Bool("is_server_stream", info.IsServerStream),
		)

		// Handle stream
		err := handler(srv, ss)

		// Calculate duration
		duration := time.Since(start)

		// Log stream end
		attrs := []slog.Attr{
			slog.String("method", method),
			slog.String("peer_addr", peerAddr),
			slog.Duration("duration", duration),
		}

		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
			logger.LogAttrs(ss.Context(), slog.LevelError, "gRPC stream failed", attrs...)
		} else {
			logger.LogAttrs(ss.Context(), slog.LevelInfo, "gRPC stream completed", attrs...)
		}

		return err
	}
}
