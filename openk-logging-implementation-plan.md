# OpenK Logging Implementation Plan

## Core Architecture Decision

We will use `slog` as our logging abstraction layer, following its design principles:
- Use `slog` interfaces throughout our codebase
- Implement handlers for specific backends (zerolog)
- Keep the logging abstraction clean and framework-agnostic

## Implementation Phases

### Phase 1: Core Setup (Week 1)

#### A. Base Configuration
- [ ] Define default `slog` logger configuration
- [ ] Implement environment-based level configuration
- [ ] Create logging initialization package
- [ ] Add basic test coverage

```go
// Basic configuration setup
type LogConfig struct {
    Level      slog.Level
    AddSource  bool
    TimeFormat string
    JSONOutput bool
}

func InitLogger(config LogConfig) *slog.Logger {
    var handler slog.Handler
    
    // Default to text handler for development
    if !config.JSONOutput {
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level:     config.Level,
            AddSource: config.AddSource,
        })
    } else {
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level:     config.Level,
            AddSource: config.AddSource,
        })
    }
    
    return slog.New(handler)
}
```

### Phase 2: Zerolog Integration (Week 1)

#### A. Zerolog Handler Implementation
- [ ] Create zerolog handler implementing `slog.Handler`
- [ ] Add configuration options
- [ ] Implement level mapping
- [ ] Add performance tests

```go
// Zerolog handler implementation
type ZerologHandler struct {
    logger  zerolog.Logger
    attrs   []slog.Attr
    groups  []string
}

func NewZerologHandler(logger zerolog.Logger) slog.Handler {
    return &ZerologHandler{
        logger: logger,
    }
}

func (h *ZerologHandler) Handle(ctx context.Context, r slog.Record) error {
    // Implement slog.Handler interface
}
```

### Phase 3: OpenTelemetry Integration (Week 2)

#### A. Context Handling
- [ ] Add trace/span extraction from context
- [ ] Implement automatic trace ID logging
- [ ] Create middleware for HTTP servers

```go
func WithTraceContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        if span := trace.SpanFromContext(ctx); span != nil {
            spanCtx := span.SpanContext()
            slog.InfoContext(ctx, "request started",
                "trace_id", spanCtx.TraceID().String(),
                "span_id", spanCtx.SpanID().String(),
            )
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### Phase 4: Error Integration (Week 2)

#### A. OpenE Integration
- [ ] Implement error field extraction
- [ ] Add sensitive data handling
- [ ] Create error logging helpers

```go
// Error logging helper
func LogError(ctx context.Context, err error) {
    var openErr *opene.Error
    if errors.As(err, &openErr) {
        // Log using structured fields from OpenE
        slog.ErrorContext(ctx, openErr.Message,
            "error_code", openErr.Code,
            "domain", openErr.Domain,
            "operation", openErr.Operation,
        )
    } else {
        slog.ErrorContext(ctx, err.Error())
    }
}
```

## Testing Strategy

### 1. Handler Tests
- [ ] Test zerolog handler implementation
- [ ] Verify level mapping
- [ ] Test context propagation
- [ ] Performance benchmarks

### 2. Integration Tests
- [ ] Test OpenTelemetry context extraction
- [ ] Verify error logging
- [ ] Test configuration loading
- [ ] End-to-end logging tests

## Documentation

### 1. Usage Guidelines
- [ ] Basic logging patterns
- [ ] Context usage
- [ ] Error logging
- [ ] Configuration options

### 2. Handler Documentation
- [ ] Zerolog handler setup
- [ ] Performance considerations
- [ ] Best practices

## Future Considerations

1. **Additional Handlers**
   - Consider zap implementation
   - Cloud logging integrations
   - Custom formatting needs

2. **Advanced Features**
   - Structured logging patterns
   - Log aggregation
   - Performance optimization

3. **Monitoring Integration**
   - Metrics extraction
   - Health checks
   - Dashboard integration

## Success Criteria

1. All code uses `slog` interface for logging
2. Zero direct dependencies on specific logging frameworks in packages
3. Comparable performance to direct zerolog usage
4. Clean integration with OpenTelemetry
5. Comprehensive test coverage

## Next Steps

1. Begin with Phase 1: Core Setup
2. Create proof-of-concept for zerolog handler
3. Review and adjust implementation plan based on findings