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

## Test Standardization
- [ ] pbkdf2_test.go updates:
  - [ ] Rename `TestGenerateSalt` to `TestPBKDF2_GenerateSalt`
  - [ ] Rename `TestDeriveKey` to `TestPBKDF2_DeriveKey`
  - [ ] Convert to table-driven tests where appropriate
  - [ ] Standardize subtest naming
  
- [ ] secure_wipe_test.go updates:
  - [ ] Rename `TestSecureWipe` to `TestMemory_SecureWipe`
  - [ ] Add table-driven tests for different input sizes
  - [ ] Add more edge cases

- [ ] aes_gcm_test.go updates:
  - [ ] Rename `TestAESGenerateKey` to `TestAES_GenerateKey`
  - [ ] Rename `TestAESGenerateNonce` to `TestAES_GenerateNonce`
  - [ ] Rename `TestAESEncryptDecrypt` to `TestAES_EncryptDecrypt`
  - [ ] Improve test organization with logical grouping

- [ ] master_key_test.go updates:
  - [ ] Add more table-driven tests
  - [ ] Standardize subtest naming
  - [ ] Improve error case coverage

- [ ] rsa_test.go updates:
  - [ ] Improve test grouping
  - [ ] Standardize error testing patterns
  - [ ] Add missing edge cases

- [ ] algorithm_test.go updates:
  - [ ] Convert simple tests to table-driven format
  - [ ] Add negative test cases

- [ ] key_pair_test.go updates:
  - [ ] Improve test organization
  - [ ] Add more edge cases
  - [ ] Standardize setup/teardown patterns

## Next Steps

### 1. Key Management Structures
- [ ] Implement DEK (Data Encryption Key) type
  - [ ] Random generation
  - [ ] Basic envelope encryption
  - [ ] State representation
  - [ ] Tests
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
- Current focus: Complete key management fundamentals
- Keep zero-knowledge architecture throughout
- Maintain interface-based design for flexibility