# ADR-009: Key Derivation Architecture

## Status
Revised (supersedes previous version)

## Context
The system needs a clear approach for deriving cryptographic keys from user credentials. Taking inspiration from proven systems like Bitwarden, we implement parallel key derivation for authentication and encryption operations while maintaining strong security properties.

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
    username,    // Used as salt
    iterations,
    keyLength    // 256 bits
)
```

#### 1.2 Auth Key Derivation
```go
// Auth Key used for API authentication, derived in parallel
AuthKey = PBKDF2(
    password,
    "openk4auth",    // Constant salt
    iterations,      // Same iteration count
    keyLength        // 256 bits
)
```

### 2. Implementation Requirements

#### 2.1 Key Operations
```go
type MasterKeyOps interface {
    // Derive both master and auth keys from password
    Derive(password, username []byte) error
    
    // Check if keys are available
    HasKey() bool
    
    // Clear sensitive key material
    Clear()
    
    // Encrypt data using master key
    Encrypt(data []byte) (*Ciphertext, error)
    
    // Decrypt data using master key
    Decrypt(ct *Ciphertext) ([]byte, error)
    
    // Get derived auth key for API operations
    GetAuthKey() ([]byte, error)
    
    // Get stable identifier for this encryption provider
    ID() string
}
```

#### 2.2 Key Usage
- Master Key: Primary encryption key for protecting KeyPairs
- Auth Key: Used exclusively for API authentication
- Independent derivation paths
- Separate rotation cycles

### 3. Security Properties

#### 3.1 Advantages
- Independent key derivation paths
- Deterministic salt for master key
- No additional salt storage needed
- Server-side key stretching
- Clear separation of auth/encryption
- Proven approach (follows Bitwarden model)

#### 3.2 Memory Protection
- Clear both keys after use
- Secure memory wiping
- Prevent key material swapping
- Clear on process exit

#### 3.3 Error Handling
- Clear sensitive data on errors
- Constant-time comparisons
- Rate limiting on authentication
- Clear error messages

### 4. Implementation Notes

#### 4.1 Constants
```go
const (
    // PBKDF2 parameters
    DefaultIterations = 100_000
    MasterKeySize    = 32       // 256 bits
    AuthSalt         = "openk4auth"
)
```

#### 4.2 Error Types
```go
var (
    ErrInvalidPassword = errors.New("invalid password")
    ErrInvalidUsername = errors.New("invalid username")
    ErrKeyNotDerived   = errors.New("master key not derived")
    ErrKeyAlreadySet   = errors.New("master key already set")
)
```

## Consequences

### Positive
- Independent keys for auth and encryption
- No key derivation chaining
- Proven industry approach
- Simple key recovery path
- Clear separation of concerns
- Minimal state management

### Negative
- Two key derivation operations required
- Fixed to PBKDF2 for both operations
- Server must maintain auth stretching
- Both keys need memory protection

## Notes
- Regular review of iteration counts
- Consider adding parameters for quantum resistance
- Monitor for new key derivation standards
- Document clear recovery procedures

## References
- crypto-spec.md: Core cryptographic specifications
- NIST SP 800-132: Key Derivation Using PBKDF
- Bitwarden Security Whitepaper