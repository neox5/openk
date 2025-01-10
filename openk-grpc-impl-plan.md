# OpenK gRPC Implementation Plan - Updated

## Current Structure ✓
```
internal/
├── server/              
│   ├── grpc_server.go   # Core server implementation ✓
│   ├── grpc_options.go  # Server options ✓
│   ├── config.go        # Server configuration ✓
│   ├── interceptors/    
│   │   └── logging.go   # Logging interceptor ✓
│   └── services/        
│       ├── users/       # User service [NEXT]
│       └── health/      # Health service ✓
└── storage/             # Storage layer ✓
    ├── models/          # Domain models ✓
    │   └── user.go      # User models ✓
    ├── memory/          # Memory implementation ✓
    └── store.go         # Core interfaces ✓
```

## Implementation Progress

### Completed Components ✓

#### 1. Server Core ✓
- [x] Basic gRPC server setup
- [x] Option organization
- [x] Configuration management
- [x] Error handling
- [x] Logging integration
- [x] Connection management
- [x] Graceful shutdown

#### 2. Storage Layer ✓
- [x] Store interfaces
- [x] User models
- [x] Memory implementation
- [x] Comprehensive test suite
  - [x] Basic operations (Create, Get)
  - [x] Error handling
  - [x] Concurrency tests
  - [x] Immutability validation

## Next Priority: User Service Implementation

### 1. Service Structure
```go
// users/users_server_v1.go
type UsersServerV1 struct {
    usersv1.UnimplementedUserServiceServer
    store   storage.UserStore
    logger  *slog.Logger
}

// users/users_register.go
func RegisterUserServers(srv *grpc.Server, store storage.UserStore, logger *slog.Logger) error {
    // Version registration
}
```

### 2. Implementation Steps

#### Phase 1: Core Service Setup
- [ ] Create service structure
- [ ] Add store dependency injection
- [ ] Setup logging configuration
- [ ] Implement registration function
- [ ] Add basic health checks

#### Phase 2: RegisterUser Implementation
- [ ] Request validation
  - [ ] Username constraints
  - [ ] Key derivation parameters
  - [ ] Cryptographic material
- [ ] Error mapping
  - [ ] Storage errors → gRPC status
  - [ ] Validation errors → status.InvalidArgument
  - [ ] Conflict errors → status.AlreadyExists
- [ ] Success response building
  - [ ] User identity creation
  - [ ] Timestamp handling
  - [ ] Public key formatting

### 3. Testing Framework

#### Server Test Structure
```go
// users/users_server_v1_test.go
type UsersServerTestSuite struct {
    store  storage.UserStore
    server *UsersServerV1
}

func (s *UsersServerTestSuite) SetupTest() {
    s.store = memory.NewUserMemoryStore()
    s.server = NewUsersServerV1(s.store, testLogger)
}
```

#### Test Categories
Following storage test patterns:
1. Basic Operations
   - [ ] Valid user registration
   - [ ] Error cases (invalid input)
   - [ ] Duplicate username handling

2. Input Validation
   - [ ] Username constraints
   - [ ] Required fields
   - [ ] Parameter validation
   - [ ] Byte length checks

3. Concurrency
   - [ ] Parallel registration attempts
   - [ ] Race condition handling
   - [ ] Resource cleanup

4. Integration
   - [ ] End-to-end flow
   - [ ] Storage interaction
   - [ ] Error propagation

### 4. Metrics & Logging

#### Operation Metrics
- [ ] Registration attempts
- [ ] Success/failure rates
- [ ] Latency tracking
- [ ] Error type distribution

#### Structured Logging
- [ ] Operation entry/exit
- [ ] Error context capture
- [ ] User activity tracking
- [ ] Performance monitoring

## Future Enhancements

### 1. Additional Operations
- [ ] User lookup methods
- [ ] Profile updates
- [ ] Key rotation support
- [ ] Account recovery

### 2. Advanced Features
- [ ] Batch operations
- [ ] Streaming support
- [ ] Rate limiting
- [ ] Cache integration

## Success Criteria

### 1. Functionality
- Complete registration flow
- Proper error handling
- Secure key storage
- Clean shutdown

### 2. Performance
- Acceptable latency (<100ms)
- Resource efficiency
- Connection stability
- Proper timeout handling

### 3. Maintainability
- Clear documentation
- >80% test coverage
- Consistent patterns
- Easy to extend

## Implementation Notes
- Follow established testing patterns from userstore
- Maintain strict input validation
- Use structured logging consistently
- Focus on error handling clarity
- Keep crypto operations side-effect free
- Consider future version compatibility

## Immediate Next Steps
1. Create users_server_v1.go structure
2. Setup test infrastructure
3. Implement core registration flow
4. Add comprehensive tests