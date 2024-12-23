package app

import (
	"context"

	"github.com/neox5/openk/internal/buildinfo"
	"github.com/neox5/openk/internal/ctx"
)

const serviceName = "openk"

// NewContext creates the base application context with service information
func NewContext(parent context.Context) context.Context {
	info := buildinfo.Get()
	return ctx.WithService(parent,
		serviceName,
		info.ShortVersion(),
		generateInstanceID(),
	)
}

// generateInstanceID creates a unique identifier for this service instance
func generateInstanceID() string {
	// TODO: Implement proper instance ID generation
	return "dev-instance"
}
