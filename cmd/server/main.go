package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/neox5/openk/internal/server"
)

func main() {
	// Create logger
	logger := slog.Default()

	// Create server config
	cfg := server.DefaultConfig()

	// Create server instance
	srv, err := server.NewServer(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx := context.Background()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
