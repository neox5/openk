# ADR-009: Key Derivation Architecture

## Status
Revised (supersedes previous version)

## Context
The system needs a clear approach for deriving cryptographic keys from user credentials. Taking inspiration from proven systems like Bitwarden, we've simplified our approach while maintaining strong security properties.

## Decision

### 1. Key Derivation Process

#### 1.1 Master Key Derivation
```go
type MasterKeyDerivation struct {
    Username    string    // Used as salt
    Iterations  int      // Default: 100,000
}

// Master Key used for encryption operations
MasterKey = PBKDF2(
    password,
    username,
    iterations
)
```

#### 1.2 Auth Key Derivation
```go
// Auth Key used for API authentication
AuthKey = PBKDF2(
    masterKey,
    "auth4openk",
    100000  // Fixed server-side iterations
)
```

### 2. Implementation Requirements

#### 2.1 Master Key Operations
```go
type MasterKeyOps interface {
    // DeriveKey generates master key from password
    DeriveKey(password []byte, username string) ([]byte, error)
    
    // DeriveAuthKey generates auth key from master key
    DeriveAuthKey(masterKey []byte) ([]byte, error)
    
    // ClearKey securely wipes key material
    ClearKey(key []byte) error
}
```

#### 2.2 Key Usage
- Master Key: Encrypts/decrypts KeyPair
- Auth Key: API authentication
- Clear separation of concerns
- Independent rotation cycles

### 3. Security Properties

#### 3.1 Advantages
- Deterministic salt (username)
- No additional salt storage needed
- Server-side key stretching
- Clear separation of auth/encryption
- Proven approach (similar to Bitwarden)

#### 3.2 Memory Protection
- Clear keys after use
- Secure memory wiping
- Prevent key material swapping

#### 3.3 Error Handling
- Clear sensitive data on errors
- Constant-time comparisons
- Rate limiting on authentication
- Clear error messages

## Consequences

### Positive
- Simpler architecture
- Proven approach
- No salt management needed
- Strong security properties
- Clear separation of concerns

### Negative
- Username changes affect key derivation
- Fixed to PBKDF2 for both operations
- Server must maintain auth stretching

## Notes
- Regular review of iteration counts
- Consider adding parameters for quantum resistance
- Monitor for new key derivation standards
- Document recovery procedures

## References
- crypto-spec.md: Core cryptographic specifications
- NIST SP 800-132: Key Derivation Using PBKDF
- Bitwarden Key Derivation Model
