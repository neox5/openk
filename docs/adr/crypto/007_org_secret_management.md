# ADR-007: Organization Secret Management

## Status
Revised (supersedes original version)

## Context
Following our revised encryption architecture (ADR-005) and key derivation model (ADR-009), we need to define how organization secrets (org_secret) are managed. The org_secret is crucial for HMAC generation and must be available to all authorized members while maintaining our end-to-end encryption guarantees.

## Decision

### 1. Organization Identity

#### 1.1 Organization KeyPair
The organization's identity is represented by a long-term KeyPair:
- Follows crypto-spec KeyPair structure
- Serves as trust anchor for organization
- Protected by organization KEK(s)
- Subject to revocation, not rotation

#### 1.2 Access Management
Each member's access is managed through:
- Individual Envelope for org_secret
- Member's KeyPair for verification
- Clear access records
- Audit capabilities

### 2. Key Protection

#### 2.1 Organization KEK Management
```go
type OrgKEKManager interface {
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
- Multiple active KEKs supported
- Clear rotation schedule
- Emergency rotation procedures
- Audit requirements

### 3. Access Patterns

#### 3.1 Initial Organization Setup
1. Generate organization KeyPair
2. Create initial organization KEK
3. Protect organization private key
4. Generate initial org_secret (DEK)
5. Create admin envelope

#### 3.2 Member Access
1. Member authenticates via Master Key
2. System verifies member's KeyPair
3. Member receives their org_secret envelope
4. Member maintains org_secret in memory
5. Used for HMAC operations as needed

### 4. Member Management

#### 4.1 Adding Members
```go
type MemberAccess struct {
    MemberID     string
    KeyPairID    string    // Member's KeyPair ID
    EnvelopeID   string    // org_secret envelope
    AccessLevel  string    // e.g., "Admin", "Member"
    CreatedAt    time.Time
    RevokedAt    *time.Time
}
```

Process:
1. Verify member's KeyPair
2. Create new envelope for org_secret
3. Record access grant
4. Enable member operations

#### 4.2 Removing Members
1. Revoke member's envelope
2. Record access revocation
3. If required by policy:
   - Rotate org_secret
   - Create new envelopes for remaining members
   - Update affected HMACs

### 5. Security Controls

#### 5.1 Access Enforcement
```go
type AccessControl interface {
    // VerifyAccess checks member's access rights
    VerifyAccess(memberID string, operation string) error
    
    // RecordAccess logs access attempts
    RecordAccess(memberID string, operation string, granted bool)
    
    // GetAccessHistory retrieves access logs
    GetAccessHistory(filter AccessFilter) ([]*AccessLog, error)
}
```

#### 5.2 Audit Requirements
- Access grant/revocation
- KEK rotation events
- org_secret usage
- Operation timestamps
- Access attempts

### 6. Emergency Procedures

#### 6.1 Key Compromise Response
1. Immediate KEK rotation
2. New org_secret generation
3. New envelopes for all members
4. HMAC recomputation
5. Incident recording

#### 6.2 Recovery Procedures
- Minimum admin requirement
- Backup key procedures
- Recovery validation
- Audit trail maintenance

### 7. Implementation Requirements

#### 7.1 Memory Protection
Following crypto-spec requirements:
- org_secret in memory only
- Clear on session end
- No disk caching
- Memory encryption when available

#### 7.2 Error Handling
- Cryptographic operation errors
- Access control violations
- State transition failures
- Clear error responses

## Consequences

### Positive
- Clear identity model
- Flexible key protection
- Strong access controls
- Emergency procedures
- Comprehensive audit
- Multiple KEK support

### Negative
- Complex key management
- Multiple state tracking
- Performance implications
- Recovery complexity

## Notes
- Define minimum admin count
- Document recovery procedures
- Monitor HMAC performance
- Regular security reviews
- Maintain clear procedures
- Consider backup strategies

## References
- crypto-spec.md: Core cryptographic specifications
- ADR-005: Revised encryption architecture
- ADR-009: Key derivation architecture
- NIST SP 800-57: Key management guidelines