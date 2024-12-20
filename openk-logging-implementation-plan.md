# OpenK Logging Implementation Plan - Complete

## Core Architecture Decision

We will use `slog` as our logging abstraction layer, with optional structured logging via zerolog:
- All internal packages use only `slog` interface
- Applications can choose to use zerolog's handler for structured/JSON output
- Keep logging configuration clean and simple

## Implementation Status

### Phase 0: Core Setup ✓

#### A. Base Logging Package ✓
```go
internal/logging/
├── config.go      # Logging configuration - DONE
└── logging.go     # slog initialization - DONE
```

Completed:
- [x] Create minimal Config struct for slog
- [x] Implement DefaultConfig()
- [x] Create InitLogger() using slog's handlers

#### B. Example Implementation ✓
```go
internal/logging/examples/
└── main.go        # Usage examples - DONE
```

Completed:
- [x] Show default slog usage
- [x] Demo development setup (console output)
- [x] Demo production setup (JSON output)
- [x] Document output formats

## Phase 1: Enhanced Logging Foundation

### Stage 1A: Context Enhancement
Priority: High
Estimated Time: 1 day

Tasks:
- [ ] Design context keys for logging metadata
- [ ] Implement request ID handling in context
- [ ] Add trace context support
- [ ] Create helper functions for context extraction
- [ ] Write tests for context utilities

### Stage 1B: Basic Sanitization
Priority: High
Estimated Time: 1 day

Tasks:
- [ ] Define sanitizer interface
- [ ] Implement basic PII sanitizer
- [ ] Add secret value redaction
- [ ] Create sanitizer registry
- [ ] Write tests for sanitizers

### Stage 1C: Configuration Updates
Priority: Medium
Estimated Time: 0.5 day

Tasks:
- [ ] Add sanitization options to config
- [ ] Add context handling options
- [ ] Update documentation
- [ ] Write tests for new config options

## Phase 2: Middleware Integration

### Stage 2A: Basic Request Logging
Priority: High
Estimated Time: 1 day

Tasks:
- [ ] Implement basic request logging middleware
- [ ] Add timing information
- [ ] Add status code logging
- [ ] Write tests for request logging

### Stage 2B: Correlation Support
Priority: Medium
Estimated Time: 1 day

Tasks:
- [ ] Add request ID middleware
- [ ] Implement header extraction
- [ ] Add trace context propagation
- [ ] Write correlation tests

### Stage 2C: Error Integration
Priority: Medium
Estimated Time: 0.5 day

Tasks:
- [ ] Integrate with OpenE error system
- [ ] Add error logging middleware
- [ ] Implement error sanitization
- [ ] Write error handling tests

## Implementation Notes

### Stage Selection Criteria
- Each stage should be completable in 1 day or less
- Stages can be implemented independently when possible
- Each stage should include its own tests
- Documentation should be updated with each stage

### Testing Strategy
- Unit tests for all new functionality
- Integration tests for middleware
- Performance benchmarks for critical paths
- Security tests for sanitization

### Security Considerations
- Never log raw secrets or credentials
- Sanitize PII consistently
- Maintain audit trail requirements
- Consider log level security implications

### Suggested First Implementation
1. Stage 1A: Context Enhancement
2. Stage 2A: Basic Request Logging

This provides the foundation for:
- Request tracking
- Basic operational visibility
- Further logging enhancements