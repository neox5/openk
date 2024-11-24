# ADR-008: Personal Vault Secret Management

## Status
Proposed

## Context
Following our encryption architecture (ADR-005), privacy-preserving metadata model (ADR-006), and organization secret management (ADR-007), we need to define how personal vault metadata is protected. Similar to organizations using org_secret for HMAC generation, personal vaults need their own vault_secret for private metadata operations.

## Decision

### 1. Vault Secret Implementation

#### 1.1 Core Structure
The vault_secret is implemented using the DEK structure defined in crypto-spec.md, ensuring:
- Consistent cryptographic properties
- Standard state management
- Single-owner envelope encryption
- Key rotation capabilities

Refer to crypto-spec.md section "DEK (Data Encryption Key)" for structure and requirements.

#### 1.2 Access Configuration
Personal vault access uses the Envelope structure from crypto-spec.md with single-owner configuration:
- One envelope per user
- Bound to user's KeyPair
- Follows crypto-spec state management
- Supports device synchronization

### 2. Access Patterns

#### 2.1 Initial Vault Setup
1. On user account creation:
   - Generate new DEK for vault_secret following crypto-spec requirements
   - Set key state to KeyStateActive
   - Create envelope using user's public key
   - Store encrypted vault_secret and envelope
   
2. Device Support:
   - Same vault_secret across user's devices
   - Access through user's KeyPair
   - Consistent HMAC generation

#### 2.2 Vault Access
1. Retrieve encrypted vault_secret
2. Decrypt envelope using KeyPair
3. Maintain vault_secret in memory per crypto-spec requirements
4. Clear from memory following crypto-spec guidelines

### 3. Personal Sharing Model

#### 3.1 Share Records
```plaintext
ShareRecord {
    id: uuid,
    secret_id: uuid,        // Reference to original secret
    owner_id: uuid,         // Original secret owner
    recipient_id: uuid,     // Share recipient
    status: "pending" | "accepted" | "rejected" | "revoked",
    envelope: Envelope,     // As defined in crypto-spec.md
    created_at: timestamp,
    updated_at: timestamp
}
```

#### 3.2 Share Creation Flow
1. Owner Operation:
   - Retrieve secret's DEK
   - Create new Envelope per crypto-spec requirements
   - Create ShareRecord with envelope
   - Set status to "pending"

2. Share Acceptance:
   - Recipient reviews pending share
   - Optional preview (if allowed)
   - Accept/reject decision:
     * Accept: Can decrypt using envelope
     * Reject: Share marked as rejected

3. Share Management:
   Owner capabilities:
   - Monitor share status
   - Revoke access (follows crypto-spec key state transitions)
   - Update original secret
   
   Recipient capabilities:
   - Accept/reject sharing
   - Remove share from view
   - Reference in personal context

### 4. Moving to Organization

When moving items to an organization:

1. Key Transition:
   - Generate new DEK following crypto-spec
   - Re-encrypt item with new DEK
   - Create envelopes per organization structure

2. Metadata Transition:
   - Generate new HMACs using org_secret
   - Create organization metadata structure
   - Remove from personal vault
   - Maintain crypto-spec key states

Note: This is a copy-then-delete operation due to different HMAC key contexts

### 5. Key Management

#### 5.1 Key Rotation
Following crypto-spec key state transitions:
1. Generate new vault_secret (DEK)
2. Set new key to KeyStateActive
3. Set old key to KeyStatePendingRotation
4. Create new envelope
5. Re-compute personal vault HMACs
6. Update stored HMACs atomically
7. Set old key to KeyStateInactive

#### 5.2 Device Management
1. New Device Addition:
   - Uses existing vault_secret
   - Accesses via user's KeyPair
   - No HMAC recomputation needed

2. Device Removal:
   - Optional vault_secret rotation
   - Follows crypto-spec rotation process if needed

### 6. Implementation Requirements

#### 6.1 Memory Protection
Following crypto-spec memory handling requirements:
* vault_secret in memory only
* Clear on:
  - Explicit vault lock
  - Session timeout
  - Tab/window close
  - System sleep
* No disk caching
* Memory encryption when available

#### 6.2 Error Handling
Implement crypto-spec error handling:
* Handle decryption failures
* Detect invalid envelopes
* Manage HMAC computation errors
* Handle rotation failures
* Process device sync conflicts

### 6. Recovery Scenarios

#### 6.1 Standard Recovery
* Based on user's KeyPair recovery
* No central recovery mechanism
* Follow crypto-spec security properties

#### 6.2 Recovery Options
* Account recovery procedures
* Recovery key management
* Backup strategies
* Data loss scenarios

All recovery methods must maintain crypto-spec security guarantees.

## Consequences

### Positive
* Consistent with crypto-spec
* Simple sharing model
* Multiple device support
* Clean separation of concerns
* Strong security properties

### Negative
* No central recovery mechanism
* Device synchronization complexity
* Memory management requirements
* All-or-nothing vault access

## Notes
* Configure vault unlock timeouts
* Implement secure memory handling
* Define sync conflict resolution
* Monitor large vault performance
* Consider partial unlock options
* Document recovery procedures

## References
* crypto-spec.md - Core cryptographic specifications
* ADR-005 - Encryption architecture
* ADR-006 - Privacy-preserving metadata
* ADR-007 - Organization secret management