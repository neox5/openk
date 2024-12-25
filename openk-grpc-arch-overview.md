# OpenK gRPC Architecture Overview

## Introduction
This document outlines the high-level architecture for OpenK's gRPC implementation, focusing on core patterns, best practices, and integration with existing components.

## Core Architecture Components

### 1. Service Organization
- Clear separation between API definition (proto) and implementation
- Versioned APIs with backward compatibility
- Service implementations follow domain boundaries
- Reuse of core domain logic between HTTP and gRPC

### 2. Middleware Architecture
Implemented through gRPC interceptors:
- Authentication & Authorization
- Logging & Tracing
- Error Translation
- Metrics Collection
- Request/Response Validation

### 3. Observability
- Structured logging through interceptors
- Metrics for both unary and streaming operations
- Distributed tracing across service boundaries
- Health checks and readiness probes

### 4. Security Model
- End-to-end encryption maintained
- Zero-knowledge architecture preserved
- Authentication through unified credential model
- Support for both human and machine identities
- TLS for transport security

## REST Gateway Integration

### 1. Overview
gRPC-Gateway provides REST API access to gRPC services by:
- Generating REST endpoints from proto definitions
- Translating HTTP/REST to gRPC calls
- Providing OpenAPI/Swagger documentation
- Managing content type negotiation
- Handling request/response validation

### 2. Integration Strategy
Initially deployed as part of the main service:
```
┌────────── OpenK Service ──────────┐
│                                   │
│   ┌─────────┐      ┌─────────┐   │
│   │  REST   │ ---> │  gRPC   │   │
│   │ Gateway │      │ Server  │   │
│   └─────────┘      └─────────┘   │
│                                   │
└───────────────────────────────────┘
```

Benefits of integrated deployment:
- Simplified operations
- Unified logging and monitoring
- Shared configuration
- Direct server access
- Single deployment unit

Future considerations for separate deployment:
- Independent scaling needs
- Client proximity requirements
- Load pattern differences
- Deployment flexibility

### 3. Request Flow
Example trace propagation:
```
Web Client -> HTTP Request -> Gateway -> gRPC Request -> gRPC Server
     |            |             |            |             |
     └─────────────────── Same Trace ID ──────────────────┘
```

## Service Patterns

### 1. Request/Response Types
- Unary calls for standard operations
- Server streaming for event notifications
- Bidirectional streaming for node synchronization
- Clear error models aligned with HTTP APIs

### 2. Data Flow
```
Client → Interceptors → Service Implementation → Domain Logic → Storage
  ↑          ↓                    ↓                  ↓           ↓
  └──────── Response ←─── Domain Objects ←─── Domain Events ←────┘
```

### 3. Connection Management
- Connection pooling
- Load balancing support
- Graceful shutdown handling
- Health checking and service discovery

## Implementation Strategy

### 1. Proto Organization
```
api/
└── v1/
    ├── health.proto     # Health checks
    ├── common.proto     # Shared definitions
    ├── service1.proto   # Domain services
    └── service2.proto
```

### 2. Service Structure
```
internal/
├── server/
    ├── grpc_server.go    # gRPC server setup
    ├── gateway.go        # REST gateway
    ├── interceptors/     # gRPC interceptors
    ├── services/        # Service implementations
    └── health/          # Health checks
└── domain/             # Core business logic
```

### 3. Error Handling
- Consistent error mapping between domain and gRPC
- Structured error details
- Clear error categorization
- Proper error propagation

## Operational Considerations

### 1. Monitoring
- Request rate and latency
- Error rates and types
- Connection pool status
- Resource utilization
- Custom business metrics

### 2. Debug Capabilities
- gRPC reflection support
- Clear logging patterns
- Trace correlation
- Error context preservation

### 3. Performance
- Connection management
- Message size limits
- Streaming thresholds
- Resource constraints

## Integration Points

### 1. Existing Components
- Authentication system
- Key management
- Storage backends
- Event system
- Logging infrastructure

### 2. Client Support
- Official client libraries
- Authentication flows
- Error handling
- Connection management
- Retry strategies

## Future Considerations

### 1. Extensibility
- New service addition
- Version upgrades
- Protocol evolution
- Feature flagging

### 2. Scalability
- Load balancing
- Service discovery
- Caching strategies
- Resource management

## Success Criteria

### 1. Functional
- Complete feature parity with HTTP API
- Consistent error handling
- Proper authentication flows
- Data consistency

### 2. Non-Functional
- Clear monitoring
- Performance metrics
- Debug capabilities
- Security compliance

### 3. Operational
- Easy deployment
- Clear documentation
- Proper logging
- Maintenance procedures
