package health

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	healthv1 "github.com/neox5/openk/internal/api_gen/openk/health/v1"
	"github.com/neox5/openk/internal/buildinfo"
	"github.com/neox5/openk/internal/opene"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// HealthServerV1 implements the health check service.
type HealthServerV1 struct {
	healthv1.UnimplementedHealthServiceServer
	logger    *slog.Logger
	startTime time.Time
	info      *buildinfo.Info
}

// NewHealthServerV1 creates a new health service implementation.
func NewHealthServerV1(logger *slog.Logger) (*HealthServerV1, error) {
	if logger == nil {
		logger = slog.Default()
	}

	return &HealthServerV1{
		logger:    logger,
		startTime: time.Now(),
		info:      buildinfo.Get(),
	}, nil
}

// Check implements the health check endpoint.
func (s *HealthServerV1) Check(ctx context.Context, req *healthv1.CheckRequest) (*healthv1.CheckResponse, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "health check requested",
		slog.String("version", s.info.Version),
		slog.Int("components_requested", len(req.Components)),
	)

	// Get component status
	components := s.checkComponents(ctx, req.Components)

	// Determine overall status based on components
	overallStatus := s.determineOverallStatus(components)

	// Create health check response
	check := &healthv1.HealthCheck{
		Status:     overallStatus,
		Version:    s.info.Version,
		Uptime:     time.Since(s.startTime).String(),
		Timestamp:  timestamppb.Now(),
		Components: components,
	}

	return &healthv1.CheckResponse{
		Health: check,
	}, nil
}

// WatchHealth implements the health check streaming endpoint.
func (s *HealthServerV1) WatchHealth(req *healthv1.WatchHealthRequest, stream healthv1.HealthService_WatchHealthServer) error {
	if req.UpdateIntervalSeconds < 1 {
		return opene.NewValidationError("health", "watch_health", "update interval must be at least 1 second").
			WithMetadata(opene.Metadata{
				"provided_interval": req.UpdateIntervalSeconds,
			})
	}

	s.logger.LogAttrs(stream.Context(), slog.LevelInfo, "health watch started",
		slog.String("version", s.info.Version),
		slog.Int("components_requested", len(req.Components)),
		slog.Uint64("interval_seconds", uint64(req.UpdateIntervalSeconds)),
	)

	ticker := time.NewTicker(time.Duration(req.UpdateIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			s.logger.LogAttrs(context.Background(), slog.LevelInfo, "health watch ended",
				slog.String("reason", "context cancelled"),
			)
			return nil

		case <-ticker.C:
			components := s.checkComponents(stream.Context(), req.Components)
			overallStatus := s.determineOverallStatus(components)

			check := &healthv1.HealthCheck{
				Status:     overallStatus,
				Version:    s.info.Version,
				Uptime:     time.Since(s.startTime).String(),
				Timestamp:  timestamppb.Now(),
				Components: components,
			}

			if err := stream.Send(check); err != nil {
				s.logger.LogAttrs(stream.Context(), slog.LevelError, "failed to send health update",
					slog.String("error", err.Error()),
				)
				return opene.NewInternalError("health", "watch_health", "failed to send update").
					WithMetadata(opene.Metadata{
						"error": err.Error(),
					})
			}
		}
	}
}

// checkComponents checks the status of requested components.
func (s *HealthServerV1) checkComponents(ctx context.Context, requested []string) []*healthv1.ComponentStatus {
	// If no specific components requested, check all
	if len(requested) == 0 {
		requested = []string{"system", "memory"}
	}

	components := make([]*healthv1.ComponentStatus, 0, len(requested))
	now := timestamppb.Now()

	for _, name := range requested {
		status := &healthv1.ComponentStatus{
			Name:      name,
			LastCheck: now,
		}

		switch name {
		case "system":
			status.Status = healthv1.Status_STATUS_UP
			status.Message = fmt.Sprintf("Running on %s/%s", runtime.GOOS, runtime.GOARCH)

		case "memory":
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)

			if mem.Sys > 1024*1024*1024 { // 1GB
				status.Status = healthv1.Status_STATUS_DEGRADED
				status.Message = fmt.Sprintf("High memory usage: %d MB", mem.Sys/1024/1024)
			} else {
				status.Status = healthv1.Status_STATUS_UP
				status.Message = fmt.Sprintf("Memory usage: %d MB", mem.Sys/1024/1024)
			}

		default:
			status.Status = healthv1.Status_STATUS_UNSPECIFIED
			status.Message = "Unknown component"
		}

		components = append(components, status)
	}

	return components
}

// determineOverallStatus determines the overall system status based on component statuses.
func (s *HealthServerV1) determineOverallStatus(components []*healthv1.ComponentStatus) healthv1.Status {
	if len(components) == 0 {
		return healthv1.Status_STATUS_UNSPECIFIED
	}

	hasDown := false
	hasDegraded := false

	for _, comp := range components {
		switch comp.Status {
		case healthv1.Status_STATUS_DOWN:
			hasDown = true
		case healthv1.Status_STATUS_DEGRADED:
			hasDegraded = true
		}
	}

	if hasDown {
		return healthv1.Status_STATUS_DOWN
	}
	if hasDegraded {
		return healthv1.Status_STATUS_DEGRADED
	}
	return healthv1.Status_STATUS_UP
}
