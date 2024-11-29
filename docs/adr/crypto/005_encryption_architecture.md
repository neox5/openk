# ADR-005: Secret Management Encryption Architecture

## Status
Revised (supersedes original version)

## Context
Following ADR-002, ADR-003, and ADR-004, we need to define a comprehensive encryption architecture that ensures end-to-end encryption (E2EE) while maintaining flexibility in authentication methods and supporting technical clients. Based on implementation experience and industry standards review, we've revised our approach to key management and identity.

## Decision

### 1. Core Encryption Strategy

#### 1.1 Key Hierarchy
This ADR builds on ADR-009 (Key Derivation Architecture) and crypto-spec.md:

1. Authentication Layer
   - User Credentials → Master Key (PBKDF2)
   - Purely for authentication
   - Rotates with password changes

2. Protection Layer
   - Independent KEK(s) (via PBKDF2 or external)
   - Storage protection focus
   - Independent rotation cycle
   - Multiple active keys possible

3. Identity Layer
   - Long-term stable KeyPair
   - Protected by KEK(s)
   - Revocation rather than rotation
   - Trust anchor for system

4. Data Protection Layer
   - Random DEK generation
   - Envelope-based distribution
   - Support for authorized rotation

#### 1.2 Authentication Layer (Pluggable)
* Support for OAuth/OIDC providers (e.g., Keycloak, Zitadel)
* Enterprise SSO integration
* Handles user authentication, MFA, session management
* Master Key derivation for authentication proof
* Independent from key protection layer

#### 1.3 Encryption Layer
* Independent encryption credentials
* Client-side encryption/decryption using crypto-spec structures
* Zero-knowledge architecture - server never has access to:
  - Plaintext secrets
  - Encryption keys
  - Private keys
  - KEK material

### 2. Key Management

#### 2.1 Initial Setup
1. User completes authentication
2. User provides encryption password
3. Client performs key derivation:
   - Generates PBKDF2 salt for Master Key
   - Derives Master Key using PBKDF2
   - Either derives or receives KEK
   - Records KEK metadata
4. Generate long-term KeyPair
5. Encrypt private key with KEK
6. Clear sensitive material from memory
7. Store encrypted private key and public key

#### 2.2 Server Storage
The server stores:
- User authentication details
- PBKDF2 parameters for Master Key
- KEK metadata and parameters
- Public key from KeyPair
- Encrypted private key
- Key state information

### 3. Technical Client Access

#### 3.1 Client Authentication & Key Protection
* Based on proven PKI patterns
* Requires:
  - clientId
  - High entropy credentials
  - Clear key protection strategy
* Separation between:
  1. Authentication credentials
  2. Key protection mechanisms

#### 3.2 Client Key Management
* Generate and protect long-term KeyPair
* Store encrypted KeyPair on server
* Support KEK rotation
* Clear revocation process

### 4. Secret Sharing

#### 4.1 Direct Sharing Process
Uses envelope encryption as defined in crypto-spec:
1. Owner decrypts DEK using their KeyPair
2. Creates new Envelope for recipient
3. Stores encrypted Envelope

#### 4.2 Shared Secret Access
* Each recipient has individual Envelope
* Original encrypted secret unchanged
* Recipients use their KeyPair to decrypt DEK
* Revocation through Envelope state

### 5. Key Protection Service (KPS)

#### 5.1 KPS Provider Requirements
Must support all operations defined in crypto-spec:
- KEK management
- Key protection
- State management
- Rotation capabilities

#### 5.2 Implementation Options
* Cloud provider KMS
* Self-hosted HSM
* OpenK's own implementation
* Hybrid approaches

All implementations must conform to crypto-spec requirements.

### 6. Security Requirements

#### 6.1 Memory Protection
* Clear Master Key after authentication
* Protect KEK during operations
* No serialization of raw key material
* Secure memory wiping
* Memory encryption where possible
* Prevent key material swapping

#### 6.2 Error Handling
Must follow error handling requirements from crypto-spec:
* Use constant-time comparisons
* Clear sensitive data on errors
* Generic error messages
* Rate limiting on critical operations

## Consequences

### Positive
* Clear separation of authentication and key protection
* Stable cryptographic identities
* Flexible KEK management
* Strong E2EE guarantees
* Proven cryptographic approaches
* Support for both user and technical clients
* Comprehensive security model

### Negative
* More complex KEK management
* Need for clear KEK rotation procedures
* Additional state tracking requirements
* Multiple active KEKs to manage

## Notes
* Consider backup procedures for key material
* Define clear KEK rotation policies
* Document emergency procedures
* Regular security audits required
* Need clear revocation procedures
* Define recovery processes

## References
* crypto-spec.md - Core cryptographic specifications
* ADR-009 - Key Derivation Architecture
* ADR-002 - Secret management model
* ADR-003 - Input processing
* ADR-004 - Service architecture
