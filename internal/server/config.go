package server

import (
	"time"

	"github.com/neox5/openk/internal/opene"
)

type Config struct {
	Host            string
	Port            int
	GRPCPort        int            // Added for gRPC server
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	EnableTLS       bool
}

func DefaultConfig() *Config {
	return &Config{
		Host:            "0.0.0.0",
		Port:            8080,
		GRPCPort:        9090,         // Default gRPC port
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		EnableTLS:       false,
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
		return opene.NewValidationError("server", "validate_config", "grpc port must be between 1 and 65535").
			WithMetadata(opene.Metadata{
				"value": c.GRPCPort,
			})
	}

	if c.GRPCPort == c.Port {
		return opene.NewValidationError("server", "validate_config", "grpc port must be different from http port").
			WithMetadata(opene.Metadata{
				"http_port": c.Port,
				"grpc_port": c.GRPCPort,
			})
	}

	return nil
}
