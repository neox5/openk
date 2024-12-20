package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/neox5/openk/internal/logging"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

func main() {
	// Example 1: Default slog setup
	cfg := logging.DefaultConfig()
	logger := logging.InitLogger(cfg)

	logger.Info("hello from slog",
		"key1", "value1",
		"key2", 42,
	)

	// Example 2A: Development setup with console output
	devLogger := zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.Kitchen,
			NoColor:    false,
		},
	).With().Timestamp().Logger()

	devOpts := slogzerolog.Option{
		Level:  slog.LevelDebug,
		Logger: &devLogger,
	}
	loggerDev := slog.New(devOpts.NewZerologHandler())

	// Example 2B: Production setup with JSON output
	prodLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	prodOpts := slogzerolog.Option{
		Level:  slog.LevelInfo,
		Logger: &prodLogger,
	}
	loggerProd := slog.New(prodOpts.NewZerologHandler())

	// Example log entries using both loggers
	baseAttrs := []any{
		"environment", "dev",
		"release", "v1.0.0",
	}

	// Development output (human-readable)
	loggerDev = loggerDev.With(baseAttrs...)
	loggerDev.Error("database error",
		"category", "sql",
		"query.duration", 1*time.Second,
		"error", fmt.Errorf("connection timeout"),
	)

	// Production output (JSON)
	loggerProd = loggerProd.With(baseAttrs...)
	loggerProd.Error("database error",
		"category", "sql",
		"query.duration", 1*time.Second,
		"error", fmt.Errorf("connection timeout"),
	)
}
