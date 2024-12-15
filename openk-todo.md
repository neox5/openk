# OpenK Implementation Todo

## Already Implemented
- [x] Core Algorithm types (RSA-OAEP-SHA256, AES-256-GCM)
- [x] Ciphertext structure with validation
- [x] Key State management
- [x] Basic PBKDF2 key derivation with salt generation
- [x] AES-256-GCM encryption operations
- [x] RSA-2048-OAEP-SHA256 operations
- [x] KeyPair type using encryption interfaces
  - [x] Generation using RSA
  - [x] State management
  - [x] Memory protection through Clear()
  - [x] Full test coverage
- [x] Update Encrypter interface with ID() method
- [x] Update all existing Encrypter implementations
  - [x] MasterKey
  - [x] UnsealedKeyPair
  - [x] Add tests for ID functionality
- [x] Create basic DEK structures
  - [x] InitialDEK
  - [x] DEK
  - [x] UnsealedDEK
- [x] Implement DEK methods
  - [x] GenerateDEK()
  - [x] Seal()
  - [x] CreateEnvelope()
  - [x] Clear() for memory protection
  - [x] Encrypt/Decrypt
- [x] Add DEK validation
- [x] Add comprehensive tests for DEK
- [x] Create Envelope structures
  - [x] InitialEnvelope
  - [x] Envelope
- [x] Define MiniSecret type
  - [x] Basic secret structure (key/value)
  - [x] Three-state model (Initial/Stored/Unsealed)
  - [x] Encryption using Encrypter interface
  - [x] DEK tracking through EncrypterID
- [x] Key Derivation Implementation
  - [x] Define structures (KeyDerivation)
  - [x] Implement methods
  - [x] Add validation
  - [x] Add tests
- [x] Key Derivation Storage - Phase 1
  - [x] Implement InMemoryMiniBackend for KeyDerivation
    - [x] StoreDerivationParams
    - [x] GetDerivationParams
- [x] Basic HTTP server planning
  - [x] Server package structure defined
  - [x] RFC 7807 error handling design (ADR-010)
  - [x] Updated CODE_STYLE.md
- [x] Error handling package
  - [x] Protocol-agnostic API error types
  - [x] API error tests
  - [x] HTTP error mapping
  - [x] Clean separation of concerns

## Next Up

### 1. Basic HTTP server implementation
- [x] Error Middleware
  - [x] Implement middleware using new error types
  - [x] Add HTTP error tests
- [x] Key Derivation endpoints
  - [x] POST /api/v1/derivation/params
  - [x] GET /api/v1/derivation/params/{username}
- [x] Update server implementation
  - [x] Modern routing with Go 1.22 features
  - [x] Consistent HandlerFunc usage
  - [x] Clean middleware composition
- [ ] Tests
  - [ ] Server tests
  - [ ] Handler tests
  - [ ] Error handling tests
  - [ ] Middleware tests

### 2. CLI Support - Phase 1
- [ ] Basic CLI implementation
  - [ ] KeyDerivation commands
  - [ ] Server interaction

### 3. Integration Testing - Phase 1
- [ ] Test KeyDerivation/MasterKey flow
  - [ ] End-to-end tests
  - [ ] Error scenarios

### 4. Authentication Implementation
- [ ] Design authentication flow
  - [ ] Define auth endpoints
  - [ ] Implement AuthKey validation
  - [ ] Session management
  - [ ] Rate limiting
  - [ ] Token management
- [ ] Update existing endpoints for auth requirements
- [ ] CLI authentication support
- [ ] Integration tests for auth flow

### 5. KMS Storage Implementation - Phase 2
- [ ] Extend InMemoryMiniBackend for KMS
  - [ ] KeyPair operations
  - [ ] DEK operations
  - [ ] Envelope operations
  - [ ] Transaction support
  - [ ] Comprehensive tests

### 6. Server Support - Phase 2
- [ ] Extend HTTP server
  - [ ] KMS endpoints
  - [ ] Error handling
  - [ ] Tests

### 7. CLI Support - Phase 2
- [ ] Extend CLI
  - [ ] KMS commands
  - [ ] Key management operations
  - [ ] Tests

### 8. Integration Testing - Phase 2
- [ ] Test complete KMS flow
  - [ ] End-to-end tests
  - [ ] Error scenarios
  - [ ] CLI interaction tests

### 9. Secret Storage Implementation - Phase 3
- [ ] Extend InMemoryMiniBackend for MiniSecret
  - [ ] Secret operations
  - [ ] Transaction support
  - [ ] Tests

### 10. Server Support - Phase 3
- [ ] Extend HTTP server
  - [ ] Secret endpoints
  - [ ] Error handling
  - [ ] Tests

### 11. CLI Support - Phase 3
- [ ] Extend CLI
  - [ ] Secret management commands
  - [ ] Tests

### 12. Integration Testing - Phase 3
- [ ] Test complete system
  - [ ] End-to-end tests
  - [ ] Performance tests
  - [ ] Security validation

### Future Work

#### Integration & Testing
- [ ] Integration Tests
  - [ ] State transition testing
  - [ ] Key lifecycle scenarios
  - [ ] Error handling validation
  - [ ] Memory leak detection
- [ ] Storage Integration Tests
  - [ ] Backend interface validation
  - [ ] Transaction handling
  - [ ] Concurrent access testing

#### Storage & Authentication
- [ ] Production Storage Layer
  - [ ] Define storage interfaces
  - [ ] Implement PostgreSQL backend
  - [ ] Add Redis caching support
  - [ ] Create MongoDB adapter
- [ ] Authentication System
  - [ ] Implement authentication flows
  - [ ] Add session management
  - [ ] Support multiple auth providers
  - [ ] MFA integration

#### Extended CLI Framework
- [ ] Advanced CLI Features
  - [ ] Configuration management
  - [ ] Advanced CRUD operations
  - [ ] Rich output formatting
- [ ] Terminal UI (TUI)
  - [ ] Secret browser interface
  - [ ] Real-time updates
  - [ ] Keyboard shortcuts

#### Sync & Recovery
- [ ] Device Synchronization
  - [ ] Define sync protocol
  - [ ] Implement conflict resolution
  - [ ] Add change tracking
  - [ ] Support offline operations
- [ ] Recovery Procedures
  - [ ] Implement key recovery
  - [ ] Add emergency access
  - [ ] Create backup systems
  - [ ] Document recovery processes

## Implementation Notes
- Maintain zero-knowledge architecture throughout
- Follow established memory protection patterns
- Keep consistent error handling
- Ensure comprehensive test coverage
- Align with vision document priorities
- Consider progressive complexity in feature rollout
- Maintain focus on security fundamentals

## Future Considerations
- Performance optimization
- Additional storage backends
- Extended authentication methods
- Advanced synchronization features
- Enterprise integration patterns
- Community contribution guidelines
- Security audit framework
