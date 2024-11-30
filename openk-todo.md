# OpenK Implementation Todo

## Already Implemented
- [x] Core Algorithm types (RSA-OAEP-SHA256, AES-256-GCM)
- [x] Ciphertext structure with validation
- [x] Key State management
- [x] PBKDF2 key derivation with salt generation
- [x] AES-256-GCM encryption operations with comprehensive tests
- [x] Implement RSA-2048-OAEP-SHA256 operations with comprehensive tests

## Next Steps

### 1. Key Management Structures
- [ ] Implement KeyPair type
  - [ ] Generation using RSA functions
  - [ ] Protection with Master Key using AES
  - [ ] Export/Import functions for keys
  - [ ] State management
  - [ ] Tests
- [ ] Implement DEK (Data Encryption Key) type
  - [ ] Random generation using AES functions
  - [ ] Envelope encryption
  - [ ] Secure deletion
  - [ ] State transitions
  - [ ] Tests
- [ ] Implement Envelope type
  - [ ] Creation with recipient key
  - [ ] Unwrapping logic
  - [ ] State management
  - [ ] Tests

### 2. Derived Key Management
- [ ] Extend PBKDF2 implementation
  - [ ] Auth Key derivation
  - [ ] Memory protection
  - [ ] Key validation
  - [ ] Tests
- [ ] Key memory handling
  - [ ] Secure wiping
  - [ ] Anti-swap protections
  - [ ] Clear key material
  - [ ] Tests

### 3. Integration Testing
- [ ] End-to-end encryption tests
- [ ] Key lifecycle tests
- [ ] Memory handling tests
- [ ] Performance benchmarks

## Later Phases
- Storage layer implementations
- Authentication flows
- Session management
- Device sync protocols
- Recovery procedures

## Notes
- Focus on core cryptographic operations first
- Ensure thorough testing with standard test vectors
- Implement secure memory handling throughout
- Maintain zero-knowledge architecture principles
