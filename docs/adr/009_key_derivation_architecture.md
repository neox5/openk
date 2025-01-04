# ADR-009: Key Derivation Architecture

## Status
Revised (supersedes previous version)

## Context
The system needs a clear approach for deriving cryptographic keys from user credentials. We implement parallel key derivation for authentication and encryption operations while maintaining strong security properties.

## Decision

### 1. Key Derivation Process

#### 1.1 Master Key Derivation (Client-side)
```go
// Master Key used for encryption operations
MasterKey = PBKDF2(
    password,
    username,    // Used as salt
    iterations,  // Default: 100,000
    keyLength    // 256 bits
)
```

#### 1.2 Auth Key Derivation and Storage
```go
// Auth Key derived from Master Key (Client-side)
AuthKey = PBKDF2(
    MasterKey,
    "openk4auth",    // Constant salt
    iterations,      // Same iteration count
    keyLength        // 256 bits
)

// Auth Key Hashing (Service Layer)
AuthKeyHash = SHA256(AuthKey)    // Computed before storage
```

### 2. Implementation Requirements

#### 2.1 Client Operations
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
}
```

#### 2.2 Service Layer Operations
```go
type UserService interface {
    // Creates user record with hashed auth key
    CreateUser(ctx context.Context, username string, authKey []byte, ...) error
}
```

#### 2.3 Key Usage
- Master Key: Primary encryption key for protecting KeyPairs
- Auth Key: Transmitted for API authentication
- Auth Key Hash: Stored form after service-layer hashing
- Independent derivation paths
- Separate rotation cycles

### 3. Security Properties

#### 3.1 Advantages
- Independent key derivation paths
- Deterministic salt for master key
- No reversible auth key storage
- Simple and secure auth key hashing
- Clear separation of auth/encryption
- Proven approach

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

## Consequences

### Positive
- Independent keys for auth and encryption
- No key derivation chaining
- Proven industry approach
- Simple key recovery path
- Clear separation of concerns
- Minimal state management
- Auth key never stored in raw form

### Negative
- Two key derivation operations required
- Fixed to PBKDF2 for initial derivation
- Raw auth key transmitted (protected by TLS)
- Multiple keys need memory protection

## Notes
- Regular review of iteration counts
- Consider adding parameters for quantum resistance
- Monitor for new key derivation standards
- Document clear recovery procedures
- Service layer responsible for auth key hashing

## References
- NIST SP 800-132: Key Derivation Using PBKDF
- Secure Hash Standard (SHA-256)
