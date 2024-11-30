# OpenK Implementation Todo

## Already Implemented
- [x] Core Algorithm types (RSA-OAEP-SHA256, AES-256-GCM)
- [x] Ciphertext structure with validation
- [x] Key State management
- [x] PBKDF2 key derivation with salt generation

## Next Steps

### 1. Core Crypto Operations
- [ ] Implement AES-256-GCM operations
  - [ ] Encryption function
  - [ ] Decryption function
  - [ ] Nonce generation
  - [ ] Tests with test vectors
- [ ] Implement RSA-2048-OAEP-SHA256 operations
  - [ ] Key generation
  - [ ] Encryption/Decryption
  - [ ] Key format handling (PKCS#8/SPKI)
  - [ ] Tests with test vectors

### 2. Key Management Structures
- [ ] Implement KeyPair type
  - [ ] Generation
  - [ ] Protection with Master Key
  - [ ] Export/Import functions
- [ ] Implement DEK type
  - [ ] Random generation
  - [ ] Envelope encryption
  - [ ] Secure deletion
- [ ] Implement Envelope type
  - [ ] Creation with recipient key
  - [ ] Unwrapping logic
  - [ ] State transitions

### 3. Derived Key Management
- [ ] Extend PBKDF2 implementation
  - [ ] Auth Key derivation
  - [ ] Memory protection
  - [ ] Key validation
- [ ] Key memory handling
  - [ ] Secure wiping
  - [ ] Anti-swap protections
  - [ ] Clear key material

### 4. Testing Infrastructure
- [ ] Test vectors for AES-GCM
- [ ] Test vectors for RSA-OAEP
- [ ] Key generation tests
- [ ] Memory handling tests
- [ ] End-to-end encryption tests

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