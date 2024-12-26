package health_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/neox5/openk/internal/api_gen/openk/health/v1"
	"github.com/neox5/openk/internal/server/services/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestHealthServerV1_Check(t *testing.T) {
	logger := slog.Default()
	server, err := health.NewHealthServerV1(logger)
	require.NoError(t, err)
	require.NotNil(t, server)

	t.Run("success cases", func(t *testing.T) {
		t.Run("check all components", func(t *testing.T) {
			ctx := context.Background()
			req := &healthv1.CheckRequest{}

			resp, err := server.Check(ctx, req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Health)

			// Verify basic health check fields
			assert.NotEmpty(t, resp.Health.Version)
			assert.NotEmpty(t, resp.Health.Uptime)
			assert.NotNil(t, resp.Health.Timestamp)

			// Should have at least system and memory components
			assert.GreaterOrEqual(t, len(resp.Health.Components), 2)

			// Verify component details
			for _, comp := range resp.Health.Components {
				assert.NotEmpty(t, comp.Name)
				assert.NotEmpty(t, comp.Message)
				assert.NotNil(t, comp.LastCheck)
				assert.NotEqual(t, healthv1.Status_STATUS_UNSPECIFIED, comp.Status)
			}
		})

		t.Run("check specific component", func(t *testing.T) {
			ctx := context.Background()
			req := &healthv1.CheckRequest{
				Components: []string{"system"},
			}

			resp, err := server.Check(ctx, req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Health)

			assert.Len(t, resp.Health.Components, 1)
			comp := resp.Health.Components[0]
			assert.Equal(t, "system", comp.Name)
			assert.Equal(t, healthv1.Status_STATUS_UP, comp.Status)
		})
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Run("unknown component", func(t *testing.T) {
			ctx := context.Background()
			req := &healthv1.CheckRequest{
				Components: []string{"unknown"},
			}

			resp, err := server.Check(ctx, req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Health)

			assert.Len(t, resp.Health.Components, 1)
			comp := resp.Health.Components[0]
			assert.Equal(t, "unknown", comp.Name)
			assert.Equal(t, healthv1.Status_STATUS_UNSPECIFIED, comp.Status)
			assert.Equal(t, "Unknown component", comp.Message)
		})
	})
}

func TestHealthServerV1_WatchHealth(t *testing.T) {
	logger := slog.Default()
	server, err := health.NewHealthServerV1(logger)
	require.NoError(t, err)
	require.NotNil(t, server)

	t.Run("success cases", func(t *testing.T) {
		t.Run("watch all components", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			req := &healthv1.WatchHealthRequest{
				UpdateIntervalSeconds: 1,
			}

			stream := &mockHealthStream{
				ctx: ctx,
				t:   t,
			}

			err := server.WatchHealth(req, stream)
			require.NoError(t, err)

			// Should have received at least 2 updates
			assert.GreaterOrEqual(t, len(stream.updates), 2)

			// Verify each update
			for _, update := range stream.updates {
				assert.NotEmpty(t, update.Version)
				assert.NotEmpty(t, update.Uptime)
				assert.NotNil(t, update.Timestamp)
				assert.GreaterOrEqual(t, len(update.Components), 2)
			}
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("invalid interval", func(t *testing.T) {
			ctx := context.Background()
			req := &healthv1.WatchHealthRequest{
				UpdateIntervalSeconds: 0,
			}

			stream := &mockHealthStream{
				ctx: ctx,
				t:   t,
			}

			err := server.WatchHealth(req, stream)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "update interval must be at least 1 second")
		})
	})
}

// Mock stream for testing WatchHealth
type mockHealthStream struct {
	ctx     context.Context
	t       *testing.T
	updates []*healthv1.HealthCheck
	grpc.ServerStream
}

func (s *mockHealthStream) Context() context.Context {
	return s.ctx
}

func (s *mockHealthStream) Send(update *healthv1.HealthCheck) error {
	s.updates = append(s.updates, update)
	return nil
}

func (s *mockHealthStream) SetHeader(metadata.MD) error {
	return nil
}

func (s *mockHealthStream) SendHeader(metadata.MD) error {
	return nil
}

func (s *mockHealthStream) SetTrailer(metadata.MD) {
}
