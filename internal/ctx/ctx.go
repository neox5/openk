package ctx

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// WithService sets service identification information in the context.
func WithService(ctx context.Context, name, version, instance string) context.Context {
	ctx = context.WithValue(ctx, KeyServiceName, name)
	ctx = context.WithValue(ctx, KeyServiceVersion, version)
	return context.WithValue(ctx, KeyServiceInstance, instance)
}

// Request Context

// WithTraceID sets the provided trace ID in the context for request correlation.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, KeyTraceID, traceID)
}

// WithRequestID sets the request identifier in the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, KeyRequestID, id)
}

// WithUserID sets the user identifier in the context.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, KeyUserID, id)
}

// WithTenantID sets the tenant identifier in the context.
func WithTenantID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, KeyTenantID, id)
}

// StartSpan begins a new span and adds it to the context. If a span already exists,
// it becomes the parent of the new span. Returns context with span information.
func StartSpan(ctx context.Context, name string, kind SpanKind) context.Context {
	// Check if we're in an existing span
	if currentSpanID := ctx.Value(KeySpanID); currentSpanID != nil {
		ctx = context.WithValue(ctx, KeySpanParentID, currentSpanID)
	}

	// Generate new span ID
	spanID := uuid.New().String()
	ctx = context.WithValue(ctx, KeySpanID, spanID)

	// Ensure we have a trace ID (create if not exists)
	if ctx.Value(KeyTraceID) == nil {
		ctx = context.WithValue(ctx, KeyTraceID, uuid.New().String())
	}

	ctx = context.WithValue(ctx, KeySpanName, name)
	ctx = context.WithValue(ctx, KeySpanKind, kind)
	ctx = context.WithValue(ctx, KeySpanStartTime, time.Now().UTC())
	return context.WithValue(ctx, KeySpanStatus, SpanStatusUnset)
}

// EndSpan marks the current span as complete and calculates its duration.
func EndSpan(ctx context.Context) {
	startTime, ok := ctx.Value(KeySpanStartTime).(time.Time)
	if !ok {
		return // No span active
	}

	duration := time.Since(startTime)
	// TODO: Hook into telemetry system to report span completion with duration
	_ = duration
}

// SetSpanStatus updates the status of the current span.
func SetSpanStatus(ctx context.Context, status SpanStatus) context.Context {
	return context.WithValue(ctx, KeySpanStatus, status)
}

// SetSpanError adds error information to the current span.
func SetSpanError(ctx context.Context, err error) context.Context {
	if err == nil {
		return ctx
	}

	ctx = context.WithValue(ctx, KeyErrorType, err.Error())
	return context.WithValue(ctx, KeyErrorMessage, err.Error())
}

// Attribute represents a key-value pair for context attributes.
type Attribute struct {
	Key   AttributeKey
	Value interface{}
}

// NewAttribute creates a new attribute with the given key and value.
func NewAttribute(key AttributeKey, value interface{}) Attribute {
	return Attribute{
		Key:   key,
		Value: value,
	}
}

// WithAttribute adds a single attribute to the context.
func WithAttribute(ctx context.Context, attr Attribute) context.Context {
	return context.WithValue(ctx, attr.Key, attr.Value)
}

// WithAttributes adds multiple attributes to the context.
func WithAttributes(ctx context.Context, attrs ...Attribute) context.Context {
	for _, attr := range attrs {
		ctx = WithAttribute(ctx, attr)
	}
	return ctx
}
