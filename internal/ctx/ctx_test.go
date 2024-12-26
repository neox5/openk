package ctx_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/neox5/openk/internal/ctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithService(t *testing.T) {
	baseCtx := context.Background()
	svcCtx := ctx.WithService(baseCtx, "test-service", "1.0.0", "instance-1")

	assert.Equal(t, "test-service", svcCtx.Value(ctx.KeyServiceName))
	assert.Equal(t, "1.0.0", svcCtx.Value(ctx.KeyServiceVersion))
	assert.Equal(t, "instance-1", svcCtx.Value(ctx.KeyServiceInstance))
}

func TestRequestContext(t *testing.T) {
	baseCtx := context.Background()

	t.Run("trace ID", func(t *testing.T) {
		traceID := "trace-123"
		traced := ctx.WithTraceID(baseCtx, traceID)
		assert.Equal(t, traceID, traced.Value(ctx.KeyTraceID))
	})

	t.Run("request ID", func(t *testing.T) {
		reqID := "req-123"
		requested := ctx.WithRequestID(baseCtx, reqID)
		assert.Equal(t, reqID, requested.Value(ctx.KeyRequestID))
	})

	t.Run("user ID", func(t *testing.T) {
		userID := "user-123"
		identified := ctx.WithUserID(baseCtx, userID)
		assert.Equal(t, userID, identified.Value(ctx.KeyUserID))
	})

	t.Run("tenant ID", func(t *testing.T) {
		tenantID := "tenant-123"
		tenanted := ctx.WithTenantID(baseCtx, tenantID)
		assert.Equal(t, tenantID, tenanted.Value(ctx.KeyTenantID))
	})
}

func TestSpanOperations(t *testing.T) {
	baseCtx := context.Background()

	t.Run("start span without parent", func(t *testing.T) {
		spanCtx := ctx.StartSpan(baseCtx, "test-operation", ctx.SpanKindInternal)

		assert.NotNil(t, spanCtx.Value(ctx.KeySpanID))
		assert.NotNil(t, spanCtx.Value(ctx.KeyTraceID))
		assert.Equal(t, "test-operation", spanCtx.Value(ctx.KeySpanName))
		assert.Equal(t, ctx.SpanKindInternal, spanCtx.Value(ctx.KeySpanKind))
		assert.NotNil(t, spanCtx.Value(ctx.KeySpanStartTime))
		assert.Equal(t, ctx.SpanStatusUnset, spanCtx.Value(ctx.KeySpanStatus))
		assert.Nil(t, spanCtx.Value(ctx.KeySpanParentID))
	})

	t.Run("start span with parent", func(t *testing.T) {
		// Create parent span
		parentCtx := ctx.StartSpan(baseCtx, "parent-operation", ctx.SpanKindServer)
		parentSpanID := parentCtx.Value(ctx.KeySpanID)
		require.NotNil(t, parentSpanID)

		// Create child span
		childCtx := ctx.StartSpan(parentCtx, "child-operation", ctx.SpanKindInternal)
		
		assert.NotNil(t, childCtx.Value(ctx.KeySpanID))
		assert.NotEqual(t, parentSpanID, childCtx.Value(ctx.KeySpanID))
		assert.Equal(t, parentSpanID, childCtx.Value(ctx.KeySpanParentID))
		assert.Equal(t, "child-operation", childCtx.Value(ctx.KeySpanName))
	})

	t.Run("end span with timing", func(t *testing.T) {
		spanCtx := ctx.StartSpan(baseCtx, "timed-operation", ctx.SpanKindInternal)
		startTime := spanCtx.Value(ctx.KeySpanStartTime).(time.Time)
		
		// Small delay to ensure measurable duration
		time.Sleep(10 * time.Millisecond)
		
		ctx.EndSpan(spanCtx)
		assert.True(t, time.Since(startTime) >= 10*time.Millisecond)
	})

	t.Run("set span status", func(t *testing.T) {
		spanCtx := ctx.StartSpan(baseCtx, "status-operation", ctx.SpanKindInternal)
		
		// Update status
		updatedCtx := ctx.SetSpanStatus(spanCtx, ctx.SpanStatusOK)
		assert.Equal(t, ctx.SpanStatusOK, updatedCtx.Value(ctx.KeySpanStatus))
	})

	t.Run("set span error", func(t *testing.T) {
		spanCtx := ctx.StartSpan(baseCtx, "error-operation", ctx.SpanKindInternal)
		testErr := errors.New("test error")
		
		errorCtx := ctx.SetSpanError(spanCtx, testErr)
		assert.Equal(t, testErr.Error(), errorCtx.Value(ctx.KeyErrorType))
		assert.Equal(t, testErr.Error(), errorCtx.Value(ctx.KeyErrorMessage))
	})

	t.Run("end span without start", func(t *testing.T) {
		// Should not panic
		ctx.EndSpan(baseCtx)
	})
}

func TestAttributes(t *testing.T) {
	baseCtx := context.Background()

	t.Run("single attribute", func(t *testing.T) {
		key := ctx.NewAttributeKey("test-key")
		attr := ctx.NewAttribute(key, "test-value")
		
		attrCtx := ctx.WithAttribute(baseCtx, attr)
		assert.Equal(t, "test-value", attrCtx.Value(key))
	})

	t.Run("multiple attributes", func(t *testing.T) {
		key1 := ctx.NewAttributeKey("key1")
		key2 := ctx.NewAttributeKey("key2")
		
		attrs := []ctx.Attribute{
			ctx.NewAttribute(key1, "value1"),
			ctx.NewAttribute(key2, 42),
		}
		
		attrCtx := ctx.WithAttributes(baseCtx, attrs...)
		assert.Equal(t, "value1", attrCtx.Value(key1))
		assert.Equal(t, 42, attrCtx.Value(key2))
	})

	t.Run("overwrite attribute", func(t *testing.T) {
		key := ctx.NewAttributeKey("test-key")
		attr1 := ctx.NewAttribute(key, "value1")
		attr2 := ctx.NewAttribute(key, "value2")
		
		ctx1 := ctx.WithAttribute(baseCtx, attr1)
		ctx2 := ctx.WithAttribute(ctx1, attr2)
		
		assert.Equal(t, "value2", ctx2.Value(key))
	})
}

func TestComplexFlow(t *testing.T) {
	// Simulate a complete request flow
	baseCtx := context.Background()
	
	// Setup service context
	svcCtx := ctx.WithService(baseCtx, "test-service", "1.0.0", "instance-1")
	
	// Start request
	reqID := uuid.New().String()
	reqCtx := ctx.WithRequestID(svcCtx, reqID)
	reqCtx = ctx.WithUserID(reqCtx, "user-123")
	reqCtx = ctx.WithTenantID(reqCtx, "tenant-456")
	
	// Start main operation
	opCtx := ctx.StartSpan(reqCtx, "main-operation", ctx.SpanKindServer)
	
	// Add custom attributes
	attrKey := ctx.NewAttributeKey("request.type")
	opCtx = ctx.WithAttribute(opCtx, ctx.NewAttribute(attrKey, "api"))
	
	// Start sub-operation
	subOpCtx := ctx.StartSpan(opCtx, "sub-operation", ctx.SpanKindInternal)
	
	// Verify complete context chain
	assert.Equal(t, "test-service", subOpCtx.Value(ctx.KeyServiceName))
	assert.Equal(t, reqID, subOpCtx.Value(ctx.KeyRequestID))
	assert.Equal(t, "user-123", subOpCtx.Value(ctx.KeyUserID))
	assert.Equal(t, "tenant-456", subOpCtx.Value(ctx.KeyTenantID))
	assert.Equal(t, "api", subOpCtx.Value(attrKey))
	assert.NotNil(t, subOpCtx.Value(ctx.KeySpanParentID))
	assert.Equal(t, "sub-operation", subOpCtx.Value(ctx.KeySpanName))
	
	// End operations
	ctx.EndSpan(subOpCtx)
	ctx.EndSpan(opCtx)
}
