package server

import (
	"errors"
	"time"
)

var (
	ErrInvalidPort = errors.New("port must be between 1 and 65535")
	ErrInvalidHost = errors.New("host cannot be empty")
)

type Config struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	EnableTLS       bool
}

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

func (c *Config) Validate() error {
	if c.Host == "" {
		return ErrInvalidHost
	}

	if c.Port < 1 || c.Port > 65535 {
		return ErrInvalidPort
	}

	return nil
}
