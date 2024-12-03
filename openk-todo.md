# OpenK Implementation Todo

## Already Implemented
- [x] Core Algorithm types (RSA-OAEP-SHA256, AES-256-GCM)
- [x] Ciphertext structure with validation
- [x] Key State management
- [x] Basic PBKDF2 key derivation with salt generation
- [x] AES-256-GCM encryption operations
- [x] RSA-2048-OAEP-SHA256 operations
- [x] KeyPair type (using encryption interfaces)
  - [x] Generation using RSA
  - [x] Encryption with Encrypter interface
  - [x] State management without transitions
  - [x] Memory protection through Clear()
  - [x] Full test coverage
- [x] Test Standardization
  - [x] pbkdf2_test.go updates
  - [x] secure_wipe_test.go updates
  - [x] aes_gcm_test.go updates
  - [x] master_key_test.go updates
  - [x] rsa_test.go updates
  - [x] algorithm_test.go updates
  - [x] key_pair_test.go updates

## Next Steps

### 1. Key Management Structures
- [x] Implement DEK (Data Encryption Key) type
  - [x] Random generation
  - [x] Basic envelope encryption
  - [x] State representation
  - [x] Tests
- [ ] Implement Envelope type
  - [ ] Creation with recipient key
  - [ ] Basic unwrapping logic
  - [ ] State representation
  - [ ] Tests

### 2. Advanced Integration Testing
- [ ] End-to-end encryption flows
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
- Current focus: Key Management Structures implementation
- Keep zero-knowledge architecture throughout
- Maintain interface-based design for flexibility
