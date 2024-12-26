package server

import (
	"time"

	"github.com/neox5/openk/internal/opene"
)

type Config struct {
	// Common server settings
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	EnableTLS       bool

	// gRPC specific settings
	GRPCPort              int           // Port for gRPC server
	MaxConnectionAge      time.Duration // Maximum duration for a connection
	MaxConnectionIdle     time.Duration // Maximum idle time for a connection
	MinConnectionTime     time.Duration // Minimum time a connection must be alive
	PermitWithoutStream   bool          // Allow connections without active streams
}

func DefaultConfig() *Config {
	return &Config{
		Host:                 "0.0.0.0",
		Port:                8080,
		GRPCPort:            9090,
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        10 * time.Second,
		ShutdownTimeout:     30 * time.Second,
		EnableTLS:           false,
		MaxConnectionAge:    30 * time.Minute,
		MaxConnectionIdle:   15 * time.Minute,
		MinConnectionTime:   5 * time.Second,
		PermitWithoutStream: true,
	}
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return opene.NewValidationError("server", "validate_config", "host cannot be empty")
	}

	if c.Port < 1 || c.Port > 65535 {
		return opene.NewValidationError("server", "validate_config", "port must be between 1 and 65535").
			WithMetadata(opene.Metadata{
				"value": c.Port,
			})
	}

	if c.GRPCPort < 1 || c.GRPCPort > 65535 {
		return opene.NewValidationError("server", "validate_config", "gRPC port must be between 1 and 65535").
			WithMetadata(opene.Metadata{
				"value": c.GRPCPort,
			})
	}

	if c.GRPCPort == c.Port {
		return opene.NewValidationError("server", "validate_config", "gRPC port must be different from HTTP port").
			WithMetadata(opene.Metadata{
				"grpc_port": c.GRPCPort,
				"http_port": c.Port,
			})
	}

	return nil
}
