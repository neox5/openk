# ADR-009: Key Derivation Architecture

## Status
Revised (supersedes original version)

## Context
The system needs clear separation between authentication and key protection mechanisms. Our revised architecture separates the Master Key (used for authentication) from the Key Encryption Key (KEK, used for protecting key material) to provide better security boundaries and more flexible key management.

### Current State
- PBKDF2 implementation exists for key derivation
- Need clear separation of authentication and protection
- Support for multiple active KEKs
- Need for flexible key protection strategies

### Industry Standards
- NIST SP 800-108: Key separation and derivation
- NIST SP 800-57: Key hierarchies and purpose binding
- PKCS#11: Key protection mechanisms
- Common practices in HSM and key protection

## Decision

### 1. Authentication Layer (Master Key)

#### 1.1 Master Key Derivation
```go
type MasterKeyParams struct {
    ID            string
    PBKDF2Salt    []byte    // 128 bits (16 bytes)
    Iterations    int       // 100,000
    OutputSize    int       // 256 bits (32 bytes)
    CreatedAt     time.Time
}
```

#### 1.2 Derivation Process
- Algorithm: PBKDF2-HMAC-SHA256
- Purpose: Authentication proof only
- Rotation: Only with password changes
- Never stored, exists only in memory
- Cleared after authentication

### 2. Key Protection Layer (KEK)

#### 2.1 KEK Sources
1. **Derived KEK**
```go
type DerivedKEKParams struct {
    ID            string
    PBKDF2Salt    []byte    // 128 bits
    Iterations    int       // Configurable
    OutputSize    int       // 256 bits
    CreatedAt     time.Time
    RotationDue   time.Time
}
```

2. **External KEK**
```go
type ExternalKEKParams struct {
    ID            string
    Source        string    // e.g., "AWS-KMS", "HSM"
    Reference     string    // Provider-specific reference
    CreatedAt     time.Time
    RotationDue   time.Time
}
```

#### 2.2 KEK Management
- Independent lifecycle from Master Key
- Support for multiple active KEKs
- Clear rotation schedule
- Status tracking for each KEK

### 3. Implementation Requirements

#### 3.1 Master Key Implementation
```go
type MasterKeyDerivation interface {
    // DeriveKey generates a master key from password
    DeriveKey(password []byte, params MasterKeyParams) ([]byte, error)
    
    // ValidateKey checks if a master key is valid
    ValidateKey(masterKey []byte) error
    
    // ClearKey securely wipes key material
    ClearKey(masterKey []byte) error
}
```

#### 3.2 KEK Implementation
```go
type KEKManager interface {
    // CreateKEK generates or imports a new KEK
    CreateKEK(params interface{}) (*KEK, error)
    
    // GetActiveKEKs returns all active KEKs
    GetActiveKEKs() ([]*KEK, error)
    
    // RotateKEK initiates KEK rotation
    RotateKEK(kekID string) error
    
    // ProtectKey encrypts key material with KEK
    ProtectKey(keyMaterial []byte, kekID string) (*Ciphertext, error)
    
    // UnprotectKey decrypts protected key material
    UnprotectKey(protected *Ciphertext, kekID string) ([]byte, error)
}
```

### 4. Key State Management

#### 4.1 Master Key States
- Derived (temporary)
- In-use (memory only)
- Cleared (after use)

#### 4.2 KEK States
```go
const (
    KEKStateActive           = "Active"
    KEKStatePendingRotation = "PendingRotation"
    KEKStateRetired         = "Retired"
    KEKStateRevoked         = "Revoked"
)
```

### 5. Rotation Procedures

#### 5.1 Master Key Rotation
- Triggered by password change
- New salt generation
- No impact on KEKs
- Clear audit trail

#### 5.2 KEK Rotation
1. Regular Rotation
   - Scheduled based on policy
   - Graceful transition
   - Maintain multiple active KEKs
   - Re-protect affected material

2. Emergency Rotation
   - Immediate new KEK
   - Fast transition
   - May force re-protection
   - Clear incident record

### 6. Security Requirements

#### 6.1 Memory Protection
- Secure allocation when possible
- Clear Master Key after auth
- Protect KEKs during operations
- Prevent key material swapping
- Memory encryption where available

#### 6.2 Error Handling
- Clear sensitive data on errors
- Constant-time comparisons
- Rate limiting on operations
- Appropriate error types

## Consequences

### Positive
- Clear separation of concerns
- Independent key lifecycles
- Flexible KEK management
- Multiple KEK support
- Better security boundaries
- Clear rotation procedures

### Negative
- More complex key management
- Additional state tracking
- Need for rotation coordination
- Multiple KEKs to manage

## Notes
- Monitor KEK usage patterns
- Define clear rotation triggers
- Document emergency procedures
- Regular security reviews
- Maintain audit trails
- Define recovery procedures

## References
- NIST SP 800-108: Key Derivation
- NIST SP 800-57: Key Management
- PKCS#11: Cryptographic Token Interface
- crypto-spec.md: Core cryptographic requirements
- ADR-005: Encryption Architecture