package logging

import (
	"log/slog"
	"os"
)

// Config holds logger configuration
type Config struct {
	// Minimum log level to output
	Level slog.Level

	// Add source code location to log entries
	AddSource bool

	// Output format (text or json)
	JSONOutput bool

	// Additional options for specific handlers
	HandlerOptions map[string]any
}

// DefaultConfig returns a default logging configuration
func DefaultConfig() *Config {
	return &Config{
		Level:          slog.LevelInfo,
		AddSource:      false,
		JSONOutput:     true,
		HandlerOptions: make(map[string]any),
	}
}

// InitLogger creates and configures the default logger
func InitLogger(cfg *Config) *slog.Logger {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	var handler slog.Handler

	// Use default handlers for initial setup/testing
	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	if cfg.JSONOutput {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
