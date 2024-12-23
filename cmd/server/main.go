package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/neox5/openk/internal/buildinfo"
	"github.com/neox5/openk/internal/ctx"
	"github.com/neox5/openk/internal/logging"
	"github.com/neox5/openk/internal/server"
)

const (
	serviceName = "openk"
)

func main() {
	// Get build information
	info := buildinfo.Get()

	// 1. Create base context with service info (FIRST!)
	serviceCtx := ctx.WithService(context.Background(),
		serviceName,
		info.ShortVersion(),
		generateInstanceID(),
	)

	// 2. Initialize logger with service context
	cfg := logging.DefaultConfig()
	logger := logging.InitLogger(cfg)

	// Log startup information
	logger.LogAttrs(serviceCtx, slog.LevelInfo, "starting OpenK server",
		slog.String("version", info.Version),
		slog.String("git_commit", info.GitCommit),
		slog.Time("build_time", info.BuildTime),
		slog.String("go_version", info.GoVersion),
	)

	// 3. Create server config
	serverCfg := server.DefaultConfig()

	// 4. Create server with service context
	srv, err := server.NewServer(serviceCtx, serverCfg, logger)
	if err != nil {
		logger.LogAttrs(serviceCtx, slog.LevelError, "failed to create server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// 5. Start server (using BaseContext configured in NewServer)
	if err := srv.Start(); err != nil {
		logger.LogAttrs(serviceCtx, slog.LevelError, "failed to start server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}

// generateInstanceID creates a unique identifier for this service instance
func generateInstanceID() string {
	// TODO: Implement proper instance ID generation
	return "dev-instance"
}
