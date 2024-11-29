# ADR-008: Personal Vault Management

## Status
Revised (supersedes original version)

## Context
Following our revised encryption architecture (ADR-005), key derivation model (ADR-009), and organization secret management (ADR-007), we need to define how personal vault metadata is protected. Similar to organizations, personal vaults need their own vault_secret for private metadata operations while maintaining consistent security properties.

## Decision

### 1. Personal Identity

#### 1.1 Personal KeyPair
The user's vault identity is represented by a long-term KeyPair:
- Follows crypto-spec KeyPair structure
- Single owner configuration
- Protected by personal KEK(s)
- Subject to revocation, not rotation

#### 1.2 Access Configuration
Personal vault access uses the standard KeyPair and envelope structure:
- One primary KeyPair per user
- Bound to user's identity
- Support for multiple devices
- Clear access tracking

### 2. Key Protection

#### 2.1 Personal KEK Management
```go
type VaultKEKManager interface {
    // CreateKEK generates or imports a new KEK
    CreateKEK(params interface{}) (*KEK, error)
    
    // GetActiveKEKs returns all active KEKs
    GetActiveKEKs() ([]*KEK, error)
    
    // RotateKEK initiates KEK rotation
    RotateKEK(kekID string) error
    
    // EmergencyRotate performs immediate KEK rotation
    EmergencyRotate() (*KEK, error)
}
```

#### 2.2 KEK Policies
- Support for multiple active KEKs
- Independent rotation schedule
- Emergency rotation capability
- Device synchronization support

### 3. Access Patterns

#### 3.1 Initial Vault Setup
1. User completes authentication (Master Key)
2. Generate personal KeyPair
3. Create initial personal KEK
4. Protect private key
5. Generate vault_secret (DEK)
6. Record vault creation

#### 3.2 Vault Access
1. User authenticates
2. System verifies KeyPair
3. Access vault_secret via envelope
4. Maintain in memory per requirements
5. Clear on session end

### 4. Device Management

#### 4.1 Device Registration
```go
type DeviceRegistration struct {
    DeviceID        string
    KeyPairID       string    // User's KeyPair ID
    EnvelopeID      string    // vault_secret envelope
    DevicePublicKey []byte    // Device-specific public key
    RegisteredAt    time.Time
    RevokedAt       *time.Time
    LastSyncAt      time.Time
}
```

Process:
1. Verify user authentication
2. Register device public key
3. Create device-specific envelope
4. Enable device synchronization

#### 4.2 Device Removal
1. Revoke device envelope
2. Record removal
3. If required by policy:
   - Rotate vault_secret
   - Update remaining device envelopes
   - Mark for synchronization

### 5. Sharing Model

#### 5.1 Share Records
```go
type ShareRecord struct {
    ID            string
    SecretID      string     // Reference to shared secret
    OwnerID       string     // Original owner
    RecipientID   string     // Share recipient
    EnvelopeID    string     // Secret envelope for recipient
    Status        ShareStatus
    CreatedAt     time.Time
    ExpiresAt     *time.Time
    RevokedAt     *time.Time
}

type ShareStatus string
const (
    ShareStatusPending   ShareStatus = "Pending"
    ShareStatusAccepted  ShareStatus = "Accepted"
    ShareStatusRejected  ShareStatus = "Rejected"
    ShareStatusRevoked   ShareStatus = "Revoked"
    ShareStatusExpired   ShareStatus = "Expired"
)
```

#### 5.2 Sharing Process
1. Owner Operation:
   - Verify secret ownership
   - Create recipient envelope
   - Record share intent
   - Notify recipient

2. Recipient Operation:
   - Review share details
   - Accept/reject share
   - Access via envelope if accepted
   - Clear rejection if declined

### 6. Synchronization

#### 6.1 State Synchronization
```go
type SyncManager interface {
    // PushChanges sends local changes to server
    PushChanges(deviceID string, changes []Change) error
    
    // PullChanges retrieves remote changes
    PullChanges(deviceID string, lastSync time.Time) ([]Change, error)
    
    // ResolveSyncConflict handles conflicting changes
    ResolveSyncConflict(change1, change2 Change) (*Change, error)
}
```

#### 6.2 Conflict Resolution
- Last-write-wins for metadata
- Merge strategy for non-conflicting changes
- User resolution for conflicts
- Clear audit trail

### 7. Security Requirements

#### 7.1 Memory Protection
Following crypto-spec requirements:
- vault_secret in memory only
- Clear on vault lock
- No disk caching
- Memory encryption when available

#### 7.2 Error Handling
- Cryptographic operation errors
- Synchronization failures
- Access control violations
- Clear error responses

### 8. Recovery Procedures

#### 8.1 Key Recovery
- Based on user's recovery process
- Clear recovery validation
- Device re-registration
- Audit trail maintenance

#### 8.2 Emergency Access
- Optional trusted contacts
- Recovery key support
- Clear activation process
- Access time limitations

## Consequences

### Positive
- Clear identity model
- Flexible key protection
- Strong device management
- Comprehensive sharing
- Recovery options
- Multi-device support

### Negative
- Complex sync requirements
- Multiple state tracking
- Recovery complexity
- Performance implications

## Notes
- Monitor sync patterns
- Define clear recovery paths
- Document emergency procedures
- Regular security reviews
- Consider backup strategies
- Test sync edge cases

## References
- crypto-spec.md: Core cryptographic specifications
- ADR-005: Revised encryption architecture
- ADR-009: Key derivation architecture
- ADR-007: Organization secret management