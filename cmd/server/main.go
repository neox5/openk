package main

import (
	"context"
	"os"

	"github.com/neox5/openk/internal/ctx"
	"github.com/neox5/openk/internal/logging"
	"github.com/neox5/openk/internal/server"
)

const (
	serviceName    = "openk"
	serviceVersion = "0.1.0" // This should come from build info
)

func main() {
	// 1. Create base context with service info (FIRST!)
	serviceCtx := ctx.WithService(context.Background(),
		serviceName,
		serviceVersion,
		generateInstanceID(),
	)

	// 2. Initialize logger with service context
	cfg := logging.DefaultConfig()
	logger := logging.InitLogger(cfg)

	// 3. Create server config
	serverCfg := server.DefaultConfig()

	// 4. Create server with service context
	srv, err := server.NewServer(serviceCtx, serverCfg, logger)
	if err != nil {
		logging.LogError(serviceCtx, logger, "server creation failed", err)
		os.Exit(1)
	}

	// 5. Start server (using BaseContext configured in NewServer)
	if err := srv.Start(); err != nil {
		logging.LogError(serviceCtx, logger, "server start failed", err)
		os.Exit(1)
	}
}

// generateInstanceID creates a unique identifier for this service instance
func generateInstanceID() string {
	// TODO: Implement proper instance ID generation
	return "dev-instance"
}
