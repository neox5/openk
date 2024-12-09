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

## Next Up
1. KeyPair Updates
   - [ ] Update KeyPair to use DEK/Envelope for protection
   - [ ] Modify InitialSeal to create and use DEK
   - [ ] Update Unseal to handle DEK/Envelope pattern
   - [ ] Update existing KeyPair tests
   - [ ] Add new tests for DEK protection

2. Integration Tests
   - [ ] Test end-to-end KeyPair creation with DEK
   - [ ] Test KeyPair unsealing flow
   - [ ] Test multiple envelopes for same DEK
   - [ ] Test state management across components
   - [ ] Test memory protection and cleanup

3. Future Considerations
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