package app

import (
	"log/slog"

	"github.com/neox5/openk/internal/buildinfo"
	"github.com/neox5/openk/internal/logging"
	"github.com/neox5/openk/internal/server"
)

func StartServer() error {
	// Create fresh application context
	ctx := NewAppContext()

	info := buildinfo.Get()

	// Initialize logger
	cfg := logging.DefaultConfig()
	logger := logging.InitLogger(cfg)

	// Log startup information
	logger.LogAttrs(ctx, slog.LevelInfo, "starting OpenK server",
		slog.String("version", info.Version),
		slog.String("git_commit", info.GitCommit),
		slog.Time("build_time", info.BuildTime),
		slog.String("go_version", info.GoVersion),
	)

	// Create and start server
	serverCfg := server.DefaultConfig()
	srv, err := server.NewServer(ctx, serverCfg, logger)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "failed to create server",
			slog.String("error", err.Error()),
		)
		return err
	}

	return srv.Start()
}
