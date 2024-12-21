package ctx

// ctxKey provides type-safety for context values using a concrete struct type
type ctxKey struct {
    id int32
}

var (
    // Service context (set once at startup)
    KeyServiceName    = ctxKey{0}
    KeyServiceVersion = ctxKey{1}
    KeyServiceInstance = ctxKey{2}

    // Trace context (from headers or generated at request start)
    KeyTraceID  = ctxKey{3}
    KeySpanID   = ctxKey{4}

    // Request identification
    KeyRequestID = ctxKey{5}
    KeyUserID    = ctxKey{6}
    KeyTenantID  = ctxKey{7}

    // Span tracking (set when span starts)
    KeySpanName      = ctxKey{8}
    KeySpanKind      = ctxKey{9}
    KeySpanStartTime = ctxKey{10}
    KeySpanParentID  = ctxKey{11}

    // Span state (updated during execution)
    KeySpanStatus = ctxKey{12}

    // Error details (set when errors occur)
    KeyErrorType    = ctxKey{13}
    KeyErrorMessage = ctxKey{14}
    KeyErrorStack   = ctxKey{15}
)

// SpanKind defines the type of span being performed
type SpanKind int32

const (
    SpanKindUnspecified SpanKind = iota
    SpanKindInternal
    SpanKindServer
    SpanKindClient
    SpanKindProducer
    SpanKindConsumer
)

// SpanStatus defines the state of a span
type SpanStatus int32

const (
    SpanStatusUnset SpanStatus = iota
    SpanStatusOK
    SpanStatusError
)

// AttributeKey represents a key for custom context attributes
type AttributeKey struct {
    name string
}

// NewAttributeKey creates a new custom attribute key
func NewAttributeKey(name string) AttributeKey {
    return AttributeKey{name: name}
}
