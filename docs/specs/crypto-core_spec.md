# OpenK Cryptographic Specification v1.0

## Key Hierarchy

### Core Concepts
1. Master Key (via PBKDF2)
   - Derived from password and username as salt
   - Used for authentication and encryption
   - Never stored, only in memory
   - Rotates with password changes

2. KeyPair (stable identity)
   - Long-term cryptographic identity
   - Protected by DEK/Envelope
   - Revocation rather than rotation
   - Trust anchor for system

3. DEK (Data Encryption Key)
   - Random key generation
   - Protected via envelopes
   - Used for data encryption
   - Multiple envelopes per DEK

4. Envelope
   - Wraps DEK for specific encryption provider
   - Supports different encryption methods
   - Independent lifecycle management
   - Access control unit

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

#### Master Key
- Algorithm: PBKDF2-HMAC-SHA256
- Salt: Username (deterministic)
- Output Size: 256 bits (32 bytes)
- Iteration Count: 100,000
- Use: Authentication and encryption
- Rotation: On password change only

#### Auth Key
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
    ID                  string     
    Algorithm           Algorithm  // AlgorithmRSA
    PublicKey           []byte     // X.509/SPKI format
    PrivateKey          Ciphertext // Encrypted with DEK
    Created             time.Time
    State               KeyState
    DEKID               string    // References protecting DEK
}

// DEK (Data Encryption Key) is a symmetric key used for encryption
type DEK struct {
    ID        string    
    Algorithm Algorithm // AlgorithmAES
    Created   time.Time
    State     KeyState
    Envelopes map[string]*Envelope // Map of envelopes by ID
}

// Envelope wraps a DEK using an encryption provider
type Envelope struct {
    ID          string     
    Algorithm   Algorithm  // Algorithm used for encryption
    DEKID       string     // References wrapped DEK
    Key         Ciphertext // Encrypted DEK
    Created     time.Time
    State       KeyState
    EncrypterID string    // References encryption provider
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
   - Generate protecting DEK
   - Create envelope for initial access
   - Record creation metadata

2. **Usage**
   - Long-term stable identity
   - Used for encryption operations
   - Trust anchor for system
   - Additional envelopes as needed

3. **Protection**
   - Protected by DEK/Envelope pattern
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
   - Initial envelope creation
   - Distribution through additional envelopes

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

## Encryption Providers

### Interface Requirements
```go
type Encrypter interface {
    // Encrypt encrypts data
    Encrypt(data []byte) (*Ciphertext, error)
    // ID returns the identifier of the encryption provider
    ID() string
}
```

### Provider Types
- Password-based (Master Key)
- Public Key (RSA)
- External KMS integration
- Hardware Security Modules

## Security Requirements

### Key Material Handling
- Clear Master Key after operations
- Clear DEKs after use
- Secure memory wiping
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
