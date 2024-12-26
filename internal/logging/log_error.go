package logging

import (
	"context"
	"log/slog"

	"github.com/neox5/openk/internal/opene"
)

// LogError logs an error with openE-aware attributes if available
func LogError(ctx context.Context, logger *slog.Logger, msg string, err error) {
	attrs := []slog.Attr{
		slog.String("error", err.Error()),
	}

	// If it's an openE error, add structured information
	if e, ok := err.(*opene.Error); ok {
		if e.Domain != "" {
			attrs = append(attrs, slog.String("domain", e.Domain))
		}
		if e.Operation != "" {
			attrs = append(attrs, slog.String("operation", e.Operation))
		}
		if e.Code != "" {
			attrs = append(attrs, slog.String("code", string(e.Code)))
		}
		// Only include metadata that adds value beyond what's in the message
		for k, v := range e.Meta {
			attrs = append(attrs, slog.Any(k, v))
		}
	}

	logger.LogAttrs(ctx, slog.LevelError, msg, attrs...)
}
