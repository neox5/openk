package interceptors

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
)

// logBuffer captures log output for testing
type logBuffer struct {
	lines []string
}

func (b *logBuffer) Write(p []byte) (n int, err error) {
	b.lines = append(b.lines, string(p))
	return len(p), nil
}

func (b *logBuffer) String() string {
	return strings.Join(b.lines, "\n")
}

func setupTest() (*logBuffer, *slog.Logger) {
	buffer := &logBuffer{}
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(buffer, opts)
	logger := slog.New(handler)
	return buffer, logger
}

// TestUnaryLogging tests the unary interceptor
func TestUnaryLogging(t *testing.T) {
	buffer, logger := setupTest()
	interceptor := UnaryLogging(logger)

	t.Run("successful request", func(t *testing.T) {
		info := &grpc.UnaryServerInfo{
			FullMethod: "/test.Service/TestMethod",
		}

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			time.Sleep(time.Millisecond) // Ensure measurable duration
			return "response", nil
		}

		_, err := interceptor(context.Background(), "request", info, handler)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		logs := buffer.String()

		// Verify request log
		if !strings.Contains(logs, "handling gRPC request") {
			t.Error("missing request log")
		}
		if !strings.Contains(logs, "TestMethod") {
			t.Error("missing method name in logs")
		}

		// Verify response log
		if !strings.Contains(logs, "gRPC request completed") {
			t.Error("missing completion log")
		}
		if !strings.Contains(logs, "duration") {
			t.Error("missing duration in logs")
		}
	})

	t.Run("request with error", func(t *testing.T) {
		buffer, logger := setupTest() // Fresh buffer
		interceptor := UnaryLogging(logger)

		expectedErr := errors.New("test error")
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, expectedErr
		}

		info := &grpc.UnaryServerInfo{
			FullMethod: "/test.Service/ErrorMethod",
		}

		_, err := interceptor(context.Background(), "request", info, handler)
		if err != expectedErr {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}

		logs := buffer.String()

		// Verify error logging
		if !strings.Contains(logs, "gRPC request failed") {
			t.Error("missing error log")
		}
		if !strings.Contains(logs, expectedErr.Error()) {
			t.Error("missing error message in logs")
		}
	})
}

// mockServerStream implements grpc.ServerStream for testing
type mockServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *mockServerStream) Context() context.Context {
	return s.ctx
}

// TestStreamLogging tests the streaming interceptor
func TestStreamLogging(t *testing.T) {
	buffer, logger := setupTest()
	interceptor := StreamLogging(logger)

	t.Run("successful stream", func(t *testing.T) {
		info := &grpc.StreamServerInfo{
			FullMethod:     "/test.Service/StreamMethod",
			IsClientStream: true,
			IsServerStream: true,
		}

		stream := &mockServerStream{ctx: context.Background()}

		handler := func(srv interface{}, stream grpc.ServerStream) error {
			time.Sleep(time.Millisecond) // Ensure measurable duration
			return nil
		}

		err := interceptor(nil, stream, info, handler)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		logs := buffer.String()

		// Verify stream start log
		if !strings.Contains(logs, "starting gRPC stream") {
			t.Error("missing stream start log")
		}
		if !strings.Contains(logs, "StreamMethod") {
			t.Error("missing method name in logs")
		}

		// Verify stream completion log
		if !strings.Contains(logs, "gRPC stream completed") {
			t.Error("missing stream completion log")
		}
		if !strings.Contains(logs, "duration") {
			t.Error("missing duration in logs")
		}
	})

	t.Run("stream with error", func(t *testing.T) {
		buffer, logger := setupTest() // Fresh buffer
		interceptor := StreamLogging(logger)

		expectedErr := errors.New("stream error")
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return expectedErr
		}

		info := &grpc.StreamServerInfo{
			FullMethod:     "/test.Service/ErrorStream",
			IsClientStream: true,
			IsServerStream: true,
		}

		stream := &mockServerStream{ctx: context.Background()}

		err := interceptor(nil, stream, info, handler)
		if err != expectedErr {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}

		logs := buffer.String()

		// Verify error logging
		if !strings.Contains(logs, "gRPC stream failed") {
			t.Error("missing stream error log")
		}
		if !strings.Contains(logs, expectedErr.Error()) {
			t.Error("missing error message in logs")
		}
	})
}
