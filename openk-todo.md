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

## KMS Package Implementation

### 1. Core Interface Updates
- [x] Update Encrypter interface with ID() method
- [x] Update all existing Encrypter implementations:
  - [x] MasterKey
  - [x] UnsealedKeyPair
  - [x] Add tests for ID functionality

### 2. DEK Implementation
- [ ] Create basic DEK structures
  - [ ] InitialDEK
  - [ ] DEK
  - [ ] UnsealedDEK
- [ ] Implement DEK methods
  - [ ] GenerateDEK()
  - [ ] InitialSeal()
  - [ ] Clear() for memory protection
- [ ] Add DEK validation
- [ ] Add comprehensive tests for DEK

### 3. Envelope Implementation
- [ ] Create Envelope structures
  - [ ] InitialEnvelope
  - [ ] Envelope
- [ ] Implement CreateEnvelope method on UnsealedDEK
- [ ] Add state validation for envelope operations
- [ ] Add comprehensive tests for Envelope

### 4. KeyPair Updates
- [ ] Update KeyPair to use DEK/Envelope for protection
- [ ] Modify InitialSeal to create and use DEK
- [ ] Update Unseal to handle DEK/Envelope pattern
- [ ] Update existing KeyPair tests
- [ ] Add new tests for DEK protection

### 5. Integration Tests
- [ ] Test end-to-end KeyPair creation with DEK
- [ ] Test KeyPair unsealing flow
- [ ] Test multiple envelopes for same DEK
- [ ] Test state management across components
- [ ] Test memory protection and cleanup

### 6. Documentation
- [ ] Update comments for all new types
- [ ] Add examples for common operations
- [ ] Document state management rules
- [ ] Document memory handling requirements

### 7. Future Considerations
- Storage layer implementations
- Authentication flows
- Session management
- Device sync protocols
- Recovery procedures

## Notes
- Maintain zero-knowledge architecture throughout
- Follow established memory protection patterns
- Keep consistent error handling
- Ensure comprehensive test coverage
