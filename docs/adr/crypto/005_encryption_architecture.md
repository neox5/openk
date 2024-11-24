# ADR-005: Secret Management Encryption Architecture

## Status
Proposed

## Context
Following the establishment of our secret management system through ADR-002, ADR-003, and ADR-004, we need to define a comprehensive encryption architecture that ensures end-to-end encryption (E2EE) while maintaining flexibility in authentication methods and supporting technical clients.

## Decision

### 1. Core Encryption Strategy

#### 1.1 Cryptographic Foundation
This ADR builds upon the cryptographic primitives and structures defined in crypto-spec.md. All cryptographic operations must conform to these specifications, including:
- Key formats and structures (KeyPair, DEK, Envelope)
- Algorithm parameters 
- Key state transitions
- Error handling requirements

For implementation details, refer to crypto-spec.md.

#### 1.2 Authentication Layer (Pluggable)
* Support for OAuth/OIDC providers (e.g., Keycloak, Zitadel)
* Enterprise SSO integration
* Handles user authentication, MFA, session management
* Independent from encryption layer

#### 1.3 Encryption Layer
* Independent encryption credentials (separate from authentication)
* Client-side encryption/decryption using crypto-spec structures
* Zero-knowledge architecture - server never has access to:
  - Plaintext secrets
  - Encryption keys
  - Client private keys

### 2. User Key Management

#### 2.1 Initial Setup
1. User completes authentication
2. User provides encryption password
3. Client performs key derivation per crypto-spec
4. Generate KeyPair following crypto-spec requirements
5. Store encrypted private key and public key on server

#### 2.2 Server Storage
The server stores:
- User authentication details
- Public key from KeyPair
- Encrypted private key per crypto-spec Ciphertext structure
- Key state information

### 3. Technical Client Access

#### 3.1 Client Authentication & Key Derivation
* Requires:
  - clientId
  - clientSecret (high entropy)
* Two-stage process:
  1. Authentication using derived key
  2. Encryption using separate derived key

#### 3.2 Client Key Management
* Generate and manage KeyPair per crypto-spec
* Store encrypted KeyPair on server
* Support key rotation per crypto-spec state transitions

### 4. Secret Sharing

#### 4.1 Direct Sharing Process
Uses envelope encryption as defined in crypto-spec:
1. Owner decrypts DEK using their KeyPair
2. Creates new Envelope for recipient using crypto-spec structure
3. Stores encrypted Envelope alongside secret

#### 4.2 Shared Secret Access
* Each recipient has individual Envelope
* Original encrypted secret remains unchanged
* Recipients use their KeyPair to decrypt DEK
* Follows envelope structure defined in crypto-spec

### 5. Key Management Service (KMS)

#### 5.1 KMS Provider Requirements
Must support all operations defined in crypto-spec:
- Key generation
- Encryption/decryption
- Key rotation
- State management

#### 5.2 Implementation Options
* Cloud provider KMS
* Self-hosted HSM
* OpenK's own implementation
* Hybrid approaches

All implementations must conform to crypto-spec requirements.

### 6. Security Requirements

#### 6.1 Memory Protection
* Clear sensitive key material from memory after use
* No serialization of raw key material
* Secure memory wiping when available
* Memory encryption where possible

#### 6.2 Error Handling
Must follow error handling requirements from crypto-spec:
* Use constant-time comparisons
* Clear sensitive data on errors
* Generic error messages
* Rate limiting on decryption operations

## Consequences

### Positive
* Clear separation between authentication and encryption
* Flexible authentication options
* Strong E2EE guarantees
* Proven cryptographic approaches
* Support for both user and technical client access
* Pluggable KMS architecture
* Comprehensive security requirements

### Negative
* Users need to manage separate encryption password
* More complex client implementation
* Performance overhead from key derivation
* Additional client-side compute requirements

## Notes
* Consider proxy re-encryption (PRE) as future optimization
* Performance impact of crypto operations needs evaluation
* Client SDK development required
* Browser compatibility testing needed
* Consider automation tools for password management
* Regular security audits required
* Need clear key rotation procedures
* Documentation needed for recovery procedures

## References
* crypto-spec.md - Core cryptographic specifications
* ADR-002 - Secret management model
* ADR-003 - Input processing
* ADR-004 - Service architecture
