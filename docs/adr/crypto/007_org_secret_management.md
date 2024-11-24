# ADR-007: Organization Secret Management

## Status
Proposed

## Context
Following our encryption architecture (ADR-005) and privacy-preserving metadata model (ADR-006), we need to define how organization secrets (org_secret) are managed. The org_secret is crucial for HMAC generation and must be available to all authorized members while maintaining our end-to-end encryption guarantees.

## Decision

### 1. Org Secret Implementation

#### 1.1 Core Structure
The org_secret is implemented using the DEK structure defined in crypto-spec.md, ensuring:
- Uniform cryptographic properties
- Consistent state management
- Standard envelope encryption for member access
- Proper key rotation support

Refer to crypto-spec.md section "DEK (Data Encryption Key)" for detailed structure and requirements.

#### 1.2 Access Management
Each member's access to the org_secret is managed through the Envelope structure defined in crypto-spec.md, providing:
- Individual access control
- Safe key distribution
- Member-specific revocation
- Audit capabilities

### 2. Access Patterns

#### 2.1 Initial Organization Setup
1. Generate new DEK for org_secret following crypto-spec requirements
2. Set initial key state to KeyStateActive
3. Create envelope for admin (organization creator) using their public key
4. Store encrypted org_secret and envelope

#### 2.2 Member Access
1. Member retrieves their org_secret envelope
2. Uses their KeyPair (per crypto-spec) to decrypt the envelope
3. Maintains decrypted org_secret in memory following crypto-spec memory handling requirements
4. Uses for HMAC operations as needed

### 3. Member Management

#### 3.1 Adding Members
1. Admin retrieves and decrypts org_secret using their envelope
2. Creates new envelope for new member following crypto-spec Envelope structure:
   - Uses member's public key
   - Follows envelope creation requirements
   - Sets appropriate key state
3. Stores new envelope with appropriate member references

#### 3.2 Removing Members
1. Remove member's envelope
2. If high-security context requires:
   - Trigger org_secret rotation (see Section 4)
   - Member can't decrypt new org_secret
3. Update access records

### 4. Key Rotation

#### 4.1 Rotation Triggers
* Regularly scheduled rotation
* Member removal (when required)
* Security incident response
* Compliance requirements

#### 4.2 Rotation Process
Following crypto-spec key state transitions:
1. Generate new DEK for org_secret
2. Set new key to KeyStateActive
3. Set old key to KeyStatePendingRotation
4. Create new envelopes for all current members
5. Re-compute affected HMACs
6. Update stored HMACs atomically
7. Set old key to KeyStateInactive

### 5. Recovery Scenarios

#### 5.1 Standard Recovery
* Requires at least one admin with access
* Admin can:
  - Generate new envelopes
  - Grant access to new members
  - Re-establish access patterns

#### 5.2 Disaster Recovery
* Multiple admin requirement recommended
* Optional backup procedures for large organizations
* All recovery must maintain crypto-spec security properties

### 6. Implementation Requirements

#### 6.1 Memory Protection
Must follow crypto-spec memory handling requirements:
* Keep org_secret in memory only
* Clear on session end/vault lock
* No disk caching
* Use memory encryption when available

#### 6.2 Error Handling
Follow crypto-spec error handling requirements:
* Handle decryption failures appropriately
* Detect missing/invalid envelopes
* Handle invalid HMAC computation
* Manage rotation failures
* Provide appropriate error responses

## Consequences

### Positive
* Consistent with crypto-spec structures
* Same key rotation mechanisms
* Clear member management process
* Unified backup/restore procedures
* Strong security guarantees

### Negative
* System-wide updates during rotation
* Performance impact during rotation
* Complex recovery scenarios
* Critical for system operation

## Notes
* Monitor performance metrics during rotation
* Document recovery procedures
* Define minimum number of required admins
* Consider caching strategies for HMAC operations
* Plan for bulk HMAC updates during rotation

## References
* crypto-spec.md - Core cryptographic specifications and structures
* ADR-005 - Encryption architecture
* ADR-006 - Privacy-preserving metadata model