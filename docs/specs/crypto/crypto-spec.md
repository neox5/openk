# OpenK Cryptographic Specification

## Key Hierarchy

### Core Concepts
1. User Credentials -> Master Key (via PBKDF2)
   - Authentication and identity proof
   - Rotates with password changes

2. Independent KEK (via PBKDF2 or external source)
   - Storage protection
   - Independent rotation cycle
   - Multiple active KEKs possible

3. KeyPair (stable identity)
   - Long-term cryptographic identity
   - Protected by KEK(s)
   - Revocation rather than rotation

4. DEK (via random generation)
   - Secret encryption
   - Wrapped via KeyPair (envelope)
   - Supports rotation for containment

Each layer provides distinct security properties with clear separation of concerns.

## Core Cryptographic Components

### Data Structures
```go
// Ciphertext represents encrypted data with authentication (AES-256-GCM)
type Ciphertext struct {
    Nonce   []byte // 96 bits
    Data    []byte // Encrypted data
    Tag     []byte // 128 bits
}

type Algorithm int
const (
    AlgorithmRSAOAEPSHA256  Algorithm = iota
    AlgorithmAESGCM256
)

// KeyPair represents a long-term identity through an asymmetric RSA key pair
type KeyPair struct {
    ID        string     
    Algorithm Algorithm  // AlgorithmRSA
    PublicKey []byte     // X.509/SPKI format
    Private   Ciphertext // Encrypted with KEK
    Created   time.Time
    State     KeyState
}

// KEK represents a key encryption key for protecting KeyPair private keys
type KEK struct {
    ID        string
    Algorithm Algorithm // AlgorithmAES
    Salt      []byte    // PBKDF2 salt if derived
    Created   time.Time
    State     KeyState
}

// DEK (Data Encryption Key) is a symmetric key that encrypts data
type DEK struct {
    ID        string    
    Algorithm Algorithm // AlgorithmAES
    Key       []byte    // In memory only, wrapped by KeyPair for distribution
    Created   time.Time
    State     KeyState
}

// Envelope wraps a DEK encrypted with a KeyPair's public key
type Envelope struct {
    ID        string     
    Algorithm Algorithm  // AlgorithmRSA
    Key       Ciphertext // DEK encrypted with recipient's public key
    Created   time.Time
    State     KeyState
    OwnerID   string     // References recipient's KeyPair.ID
}

type KeyState int
const (
    // Key is enabled and available for use
    KeyStateActive KeyState = iota
    // Key is scheduled for rotation or awaiting a state change
    KeyStatePendingRotation
    // Key is disabled and temporarily unavailable but can be re-enabled
    KeyStateInactive
    // Key is permanently deactivated and cannot be re-enabled
    KeyStateDestroyed
)
```

### Algorithm Parameters

#### Identity Keys (RSA)
- Algorithm: RSA-2048 with OAEP
- Hash Algorithm: SHA-256 (for both message digest and MGF1)
- Key Format: PKCS#8 (Private Key), X.509/SPKI (Public Key)
- Use: Long-term identity and key wrapping operations
- Rotation: Not intended for rotation, focus on protection

#### Storage Protection (KEK)
- Algorithm: AES-256-GCM
- Key Size: 256 bits
- Independent derivation parameters
- Multiple active keys supported
- Regular rotation capability

#### Data Encryption (DEK)
- Algorithm: AES-256-GCM
- Key Size: 256 bits
- Nonce Size: 96 bits (random)
- Auth Tag Size: 128 bits
- Regular rotation supported

### Key Derivation

#### Master Key (Authentication)
- Algorithm: PBKDF2-HMAC-SHA256
- Salt Size: 128 bits (16 bytes)
- Output Size: 256 bits (32 bytes)
- Iteration Count: 100,000
- Use: Authentication proof
- Rotation: On password change only

#### KEK Derivation (if not externally provided)
- Algorithm: PBKDF2-HMAC-SHA256
- Independent salt and iteration count
- Output Size: 256 bits (32 bytes)
- Use: KeyPair protection
- Rotation: Independent cycle

## Identity Management

### KeyPair Lifecycle
1. **Creation**
   - Generate during user/entity initialization
   - Strong protection with KEK
   - Record creation metadata

2. **Usage**
   - Long-term stable identity
   - Used for DEK wrapping
   - Trust anchor for system

3. **Protection**
   - Multiple KEK support
   - Regular KEK rotation
   - Secure storage requirements

4. **Revocation**
   - Clear revocation status
   - Timestamp and reason
   - No recovery after revocation

### KEK Management

#### Independent KEK Lifecycle
1. **Creation**
   - Either derived or externally provided
   - Unique identifier and metadata
   - Clear protection scope

2. **Rotation**
   - Independent rotation schedule
   - Can maintain multiple active KEKs
   - Clean transition process

3. **Emergency Procedures**
   - Fast rotation capability
   - Multiple KEK support
   - Clear recovery process

## Data Protection

### DEK Lifecycle
1. **Generation**
   - Random generation
   - Envelope creation
   - Distribution to authorized users

2. **Rotation Triggers**
   - Regular schedule
   - Security incidents
   - Access revocation
   - Compliance requirements

3. **Rotation Process**
   - New DEK generation
   - Envelope updates
   - Data re-encryption
   - Clean state transition

## Security Requirements

### Key Material Handling
- Clear Master Key after authentication
- Protect KEKs during use
- Clear unwrapped private keys
- Secure memory handling
- Anti-swapping measures

### Error Handling
- Clear sensitive data
- Constant-time operations
- Rate limiting on critical operations
- Appropriate error types

## Standards Compliance
- FIPS 197 (AES)
- FIPS 198-1 (HMAC)
- NIST SP 800-38D (GCM)
- NIST SP 800-56B (RSA)
- NIST SP 800-108 (Key Derivation)
- PKCS#1 v2.2 (RSA OAEP)
- PKCS#8 (Private Key Info)
