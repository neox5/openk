package health

import (
	"log/slog"

	healthv1 "github.com/neox5/openk/internal/api_gen/openk/health/v1"
	"google.golang.org/grpc"
)

// RegisterHealthServers registers all versions of the health service.
func RegisterHealthServers(srv *grpc.Server, logger *slog.Logger) (*HealthServerV1, error) {
	// Create v1 server
	v1Server, err := NewHealthServerV1(logger)
	if err != nil {
		return nil, err
	}

	// Register v1 server
	healthv1.RegisterHealthServiceServer(srv, v1Server)

	return v1Server, nil
}
