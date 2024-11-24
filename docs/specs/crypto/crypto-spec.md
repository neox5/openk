# OpenK Cryptographic Specification

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

// KeyPair represents an asymmetric RSA key pair used for key wrapping
type KeyPair struct {
    ID        string     
    Algorithm Algorithm  // AlgorithmRSA
    PublicKey []byte     // X.509/SPKI format
    Private   Ciphertext // Encrypted with KEK derived from user credentials
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

#### RSA Key Wrapping
- Algorithm: RSA-2048 with OAEP
- Hash Algorithm: SHA-256 (for both message digest and MGF1)
- Key Format: PKCS#8 (Private Key), X.509/SPKI (Public Key)
- Use: Key wrapping operations

#### Symmetric Encryption
- Algorithm: AES-256-GCM
- Key Size: 256 bits
- Nonce Size: 96 bits (random)
- Auth Tag Size: 128 bits
- Use: Data encryption operations

#### Key Derivation
- Algorithm: PBKDF2-HMAC-SHA256
- Salt Size: 128 bits
- Iteration Count: 100,000
- Use: Deriving keys from user credentials

## Cryptographic Processes

### Key Generation
1. **DEK Generation**
   - Generate 256-bit random key using CSPRNG
   - Create DEK structure with AlgorithmAES
   - State set to KeyStateActive

2. **KeyPair Generation**
   - Generate RSA-2048 key pair using CSPRNG
   - Export private key in PKCS#8 format
   - Export public key in X.509/SPKI format
   - Create KeyPair structure with AlgorithmRSA
   - State set to KeyStateActive

### Key Wrapping
1. **Wrapping DEK**
   - Input: DEK, recipient's public key
   - Generate Envelope with new UUID
   - Encrypt DEK using RSA-OAEP
   - Store result in Envelope.Key as Ciphertext
   - Set OwnerID to recipient's KeyPair ID

2. **Wrapping Private Key**
   - Input: RSA private key, derived KEK
   - Encrypt private key using AES-256-GCM
   - Store result in KeyPair.Private as Ciphertext

### Key Rotation
1. **DEK Rotation**
   - Generate new DEK with KeyStateActive
   - Set old DEK state to KeyStatePendingRotation
   - Re-encrypt all data with new DEK
   - Create new Envelopes for all recipients
   - Set old DEK state to KeyStateInactive

2. **KeyPair Rotation**
   - Generate new KeyPair with KeyStateActive
   - Set old KeyPair state to KeyStatePendingRotation
   - Re-wrap all relevant DEKs
   - Set old KeyPair state to KeyStateInactive

### Key State Transitions
1. **Active → PendingRotation**
   - Triggered by rotation initiation
   - Key still usable for decryption
   - New operations use replacement key

2. **PendingRotation → Inactive**
   - All dependent data re-encrypted
   - All dependent keys re-wrapped
   - Key no longer used for operations
   - Can be reactivated if needed

3. **Any State → Destroyed**
   - Permanent, irreversible operation
   - Key material securely erased
   - All references marked as destroyed
   - Cannot be used for any operations

## Security Considerations

### Key Material Handling
- Clear DEK.Key from memory after use
- Never serialize DEK.Key to storage
- Clear unwrapped private keys after use
- Clear derived KEKs after use

### Nonce Requirements
- Always generate new random nonces
- Never reuse nonces with same key
- Use CSPRNG for nonce generation

### Error Handling
- Clear sensitive data on errors
- Never expose key material in errors
- Use constant-time comparisons
- Implement rate limiting on decryption

## Standards Compliance
- FIPS 197 (AES)
- FIPS 198-1 (HMAC)
- NIST SP 800-38D (GCM)
- NIST SP 800-56B (RSA KTS)
- PKCS#1 v2.2 (RSA OAEP)
- PKCS#8 (Private Key Format)
