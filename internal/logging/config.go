package logging

import (
	"log/slog"
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
