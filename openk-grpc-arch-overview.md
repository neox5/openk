# OpenK gRPC Architecture Overview

## Directory Organization

### API Layer
```
proto/                      # Proto definitions
├── buf.yaml               # Buf configuration
├── buf.gen.yaml           # Code generation config
├── openk/                 # API namespace
│   └── <service>/        # Service directories
│       ├── v1/           # Version directories
│       │   ├── service.proto  # Service definitions
│       │   └── types.proto    # Type definitions
│       └── v2/
└── vendor/               # Vendored dependencies
    └── google/
        └── protobuf/

internal/
├── api_gen/             # Generated code
│   └── openk/
│       └── <service>/
│           └── v{n}/
└── server/              # Server implementation
    ├── grpc_server.go   # gRPC server setup
    ├── gateway.go       # REST gateway
    ├── interceptors/    # gRPC interceptors
    ├── services/       # Service implementations
    └── health/         # Health checks
```

## Core Architecture Components

### 1. Service Organization
- Proto definitions in `proto/openk/<service>/v{n}/`
- Generated code in `internal/api_gen/openk/<service>/v{n}/`
- Service implementations in `internal/server/services/`
- Clear separation between API and implementation
- Core logic reuse between HTTP and gRPC

### 2. Middleware Architecture
Implemented through interceptors in `internal/server/interceptors/`:
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

## Gateway Integration

### 1. Overview
Gateway integration (`internal/server/gateway.go`) provides:
- REST API access to gRPC services
- OpenAPI/Swagger documentation
- Content type negotiation
- Request/response validation

### 2. Integration Strategy
Initially deployed as part of the main service:
```
┌────────── Internal Server ──────────┐
│                                    │
│   ┌─────────┐       ┌─────────┐    │
│   │ gateway │ ───── │  grpc   │    │
│   │  .go    │       │ server  │    │
│   └─────────┘       └─────────┘    │
│                                    │
└────────────────────────────────────┘
```

### 3. Request Flow
```
Client Request → Gateway → gRPC Server → Service Implementation
     │             │           │                │
     └─────────── Trace ID Propagation ────────┘
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

## Implementation Guidelines

### 1. Proto Organization
All protos in `proto/openk/<service>/v{n}/`:
- `service.proto`: Service definitions
- `types.proto`: Data structures
- Clear versioning
- Shared types in common

### 2. Service Structure
Each service in `internal/server/services/<service>`:
- Clear interface definition
- Separation from transport
- Reusable business logic
- Proper error handling

### 3. Error Handling
- Consistent error mapping
- Structured error details
- Clear error categorization
- Proper error propagation

## Operational Considerations

### 1. Monitoring
- Request rate and latency
- Error rates and types
- Connection pool status
- Resource utilization

### 2. Debug Capabilities
- gRPC reflection support
- Clear logging patterns
- Trace correlation
- Error context preservation

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