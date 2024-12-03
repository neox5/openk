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
- [x] pbkdf2_test.go updates:
  - [x] Rename `TestGenerateSalt` to `TestPBKDF2_GenerateSalt`
  - [x] Rename `TestDeriveKey` to `TestPBKDF2_DeriveKey`
  - [x] Convert to table-driven tests
  - [x] Standardize subtest naming
  - [x] Review with new TESTING_GUIDE

- [x] secure_wipe_test.go updates:
  - [x] Rename `TestSecureWipe` to `TestMemory_SecureWipe`
  - [x] Add table-driven tests for different input sizes
  - [x] Add more edge cases
  - [x] Focus tests on zero-verification only
  - [x] Review with new TESTING_GUIDE

- [x] aes_gcm_test.go updates:
  - [x] Rename `TestAESGenerateKey` to `TestAES_GenerateKey`
  - [x] Rename `TestAESGenerateNonce` to `TestAES_GenerateNonce`
  - [x] Rename `TestAESEncryptDecrypt` to `TestAES_EncryptDecrypt`
  - [x] Improve test organization with logical grouping
  - [x] Review with new TESTING_GUIDE

- [x] master_key_test.go updates:
  - [x] Add more table-driven tests
  - [x] Standardize subtest naming
  - [x] Improve error case coverage
  - [x] Review with new TESTING_GUIDE

- [x] rsa_test.go updates:
  - [x] Improve test grouping
  - [x] Standardize error testing patterns
  - [x] Add missing edge cases

- [x] algorithm_test.go updates:
  - [x] Convert simple tests to table-driven format
  - [x] Add negative test cases

- [x] key_pair_test.go updates:
  - [x] Improve test organization
  - [x] Add more edge cases
  - [x] Standardize setup/teardown patterns

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
- Current focus: Complete test standardization
- Keep zero-knowledge architecture throughout
- Maintain interface-based design for flexibility
