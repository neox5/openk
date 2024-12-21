package server

import (
	"errors"
	"time"
)

var (
	ErrInvalidPort = errors.New("port must be between 1 and 65535")
	ErrInvalidHost = errors.New("host cannot be empty")
)

// Config defines the HTTP server configuration
type Config struct {
	// Basic server settings
	Host string // Server host
	Port int    // Server port

	// Timeouts
	ReadTimeout     time.Duration // Maximum duration for reading entire request
	WriteTimeout    time.Duration // Maximum duration for writing response
	ShutdownTimeout time.Duration // Maximum duration to wait for server shutdown

	// Additional options
	EnableTLS bool // Enable TLS/HTTPS
}

// DefaultConfig returns a new Config instance with default values
func DefaultConfig() *Config {
	return &Config{
		Host:            "0.0.0.0",
		Port:            8080,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		EnableTLS:       false,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Host == "" {
		return ErrInvalidHost
	}

	if c.Port < 1 || c.Port > 65535 {
		return ErrInvalidPort
	}

	return nil
}
