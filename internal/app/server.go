package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/neox5/openk/internal/buildinfo"
	"github.com/neox5/openk/internal/logging"
	"github.com/neox5/openk/internal/server"
)

// StartServer initializes and starts both HTTP and gRPC servers
func StartServer() error {
	// Create fresh application context
	ctx := NewAppContext()

	info := buildinfo.Get()

	// Initialize logger
	cfg := logging.DefaultConfig()
	logger := logging.InitLogger(cfg)

	// Log startup information
	logger.LogAttrs(ctx, slog.LevelInfo, "starting openK server",
		slog.String("version", info.Version),
		slog.String("git_commit", info.GitCommit),
		slog.Time("build_time", info.BuildTime),
		slog.String("go_version", info.GoVersion),
	)

	// Create server configuration
	serverCfg := server.DefaultConfig()

	// Create gRPC server
	grpcServer, err := server.NewGRPCServer(ctx, serverCfg, logger)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "failed to create gRPC server",
			slog.String("error", err.Error()),
		)
		return err
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start gRPC server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := grpcServer.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		logger.LogAttrs(ctx, slog.LevelInfo, "received shutdown signal",
			slog.String("signal", sig.String()),
		)
	case err := <-errChan:
		logger.LogAttrs(ctx, slog.LevelError, "server error",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("server error: %w", err)
	}

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := grpcServer.Stop(shutdownCtx); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "shutdown error",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("shutdown error: %w", err)
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "server shutdown complete")
	return nil
}
