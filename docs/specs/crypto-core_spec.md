# OpenK Cryptographic Specification v1.0

## Key Hierarchy

### Core Concepts
1. Password + Username -> Master Key (via PBKDF2)
   - Primary key for encryption operations
   - Uses username as deterministic salt
   - Used for KeyPair protection
   - Exists in memory only

2. Master Key -> Auth Key (via PBKDF2)
   - Used for API authentication
   - Server-side key stretching
   - Independent from encryption operations

3. KeyPair (stable identity)
   - Long-term cryptographic identity
   - Protected by Master Key
   - Revocation rather than rotation

4. DEK (via random generation)
   - Secret encryption
   - Wrapped via KeyPair (envelope)
   - Supports rotation for containment

## Core Cryptographic Components

### Algorithm Parameters

#### Identity Keys (RSA)
- Algorithm: RSA-2048 with OAEP
- Hash Algorithm: SHA-256 (for both message digest and MGF1)
- Key Format: PKCS#8 (Private Key), X.509/SPKI (Public Key)
- Use: Long-term identity and key wrapping operations
- Rotation: Not intended for rotation, focus on protection

#### Data Encryption (DEK)
- Algorithm: AES-256-GCM
- Key Size: 256 bits
- Nonce Size: 96 bits (random)
- Auth Tag Size: 128 bits
- Regular rotation supported

### Key Derivation

#### Master Key (Encryption)
- Algorithm: PBKDF2-HMAC-SHA256
- Salt: Username (deterministic)
- Output Size: 256 bits (32 bytes)
- Iteration Count: 100,000
- Use: KeyPair protection
- Rotation: On password change only

#### Auth Key (Authentication)
- Algorithm: PBKDF2-HMAC-SHA256
- Input: Master Key
- Salt: "auth4openk" (constant)
- Output Size: 256 bits (32 bytes)
- Iteration Count: 100,000 (server-side)
- Use: API authentication
- Rotation: With Master Key

### Data Structures
```go
// Ciphertext represents encrypted data with authentication
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
    Private   Ciphertext // Encrypted with Master Key
    Created   time.Time
    State     KeyState
}

// DEK (Data Encryption Key) is a symmetric key that encrypts data
type DEK struct {
    ID        string    
    Algorithm Algorithm // AlgorithmAES
    Key       []byte    // In memory only, wrapped by KeyPair
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
    OwnerID   string    // References recipient's KeyPair.ID
}

type KeyState int
const (
    KeyStateActive KeyState = iota
    KeyStatePendingRotation
    KeyStateInactive
    KeyStateDestroyed
)
```

## Identity Management

### KeyPair Lifecycle
1. **Creation**
   - Generate during user/entity initialization
   - Protect with Master Key
   - Record creation metadata

2. **Usage**
   - Long-term stable identity
   - Used for DEK wrapping
   - Trust anchor for system

3. **Protection**
   - Encrypted with Master Key
   - Clear protection requirements
   - Secure storage standards

4. **Revocation**
   - Clear revocation status
   - Timestamp and reason
   - No recovery after revocation

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
- Clear Auth Key after session establishment
- Clear unwrapped private keys
- Secure memory handling
- Anti-swapping measures

### Error Handling
- Clear sensitive data
- Constant-time operations
- Rate limiting on critical operations
- Appropriate error types

## Standards Compliance (2024)
| Standard | Version | Aspect |
|----------|---------|---------|
| FIPS 197 | 2001 | AES implementation |
| FIPS 198-1 | 2008 | HMAC usage |
| NIST SP 800-38D | 2007 | GCM mode |
| NIST SP 800-56B | Rev. 2, 2019 | RSA usage |
| NIST SP 800-108 | Rev. 1, 2022 | Key derivation |
| PKCS#1 | v2.2, RFC 8017 | RSA OAEP |
| PKCS#8 | RFC 5208 | Private key info |