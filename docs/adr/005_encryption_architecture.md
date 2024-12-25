# ADR-005: Secret Management Encryption Architecture

## Status
Revised (supersedes previous version)

## Context
Following ADR-002, ADR-003, and ADR-004, we need to define a comprehensive encryption architecture that ensures end-to-end encryption (E2EE) while maintaining flexibility in authentication methods and supporting both human and machine identities. Based on implementation experience and industry standards review, we've revised our approach to key management and identity.

## Decision

### 1. Core Encryption Strategy

#### 1.1 Key Hierarchy
This ADR builds on ADR-009 (Key Derivation Architecture):

1. Authentication & Protection Layer
   - User/Machine Credentials → Master Key (PBKDF2)
   - Used for both authentication and key protection
   - Rotates with credential changes
   - Protected in memory, never stored

2. Identity Layer
   - Long-term stable KeyPair
   - Protected by Master Key
   - Revocation rather than rotation
   - Trust anchor for system

3. Data Protection Layer
   - Random DEK generation
   - Envelope-based distribution
   - Support for authorized rotation

### 2. Identity Types

#### 2.1 Human Identity
1. Core Components:
   - Username (identity)
   - Master password (protection)
   - KeyPair (cryptographic identity)
   
2. Key Derivation:
   - Master Key: PBKDF2(password, username as salt)
   - Auth Key: PBKDF2(Master Key, constant salt)

3. Authentication Methods:
   a. Direct:
      - Username + derived Auth Key
      - Master Key used locally for decryption
   
   b. SSO Integration:
      - SSO authentication first
      - Master password still required
      - Auth Key derived and used as in direct method
      - Zero-knowledge principle maintained

#### 2.2 Machine Identity
1. Core Components:
   - API Key ID (identity)
   - API Key Secret (protection)
   - Machine KeyPair (cryptographic identity)

2. Key Derivation (identical to human process):
   - Master Key: PBKDF2(apiKeySecret, apiKeyID as salt)
   - Auth Key: PBKDF2(Master Key, constant salt)

3. Creation Process:
   a. Identity Generation:
      - Generate API Key ID (unique identifier)
      - Generate API Key Secret (high entropy)
      - Create Machine KeyPair
   
   b. Protection:
      - Derive Machine Master Key
      - Encrypt Machine KeyPair
      - Store encrypted KeyPair
   
   c. Access Grant:
      - Human user decrypts target secret's DEK
      - Creates new envelope with machine's public key
      - Stores new envelope for access

4. Authentication Flow:
   - Client has API Key ID + Secret
   - Derives Auth Key for server authentication
   - Uses Master Key locally for decryption
   - Zero-knowledge principle maintained

### 3. Authentication Layer
* Support for all identity types
* Consistent derivation patterns
* Independent from key protection
* Zero-knowledge architecture

### 4. Encryption Layer
* Client-side encryption/decryption
* Unified envelope encryption model
* Support for multiple access grants
* Clear revocation process

## Consequences

### Positive
* Unified cryptographic model for all identities
* Consistent zero-knowledge architecture
* Simple key hierarchy
* Strong E2EE guarantees
* Clean separation of concerns
* Support for automation
* SSO integration without compromising security

### Negative
* Credential changes require re-encryption
* Must protect machine credentials
* Need robust key rotation procedures
* Higher complexity for automation tools

## Notes
* API Key Secret should have high entropy
* Consider API Key ID format standards
* Document revocation procedures
* Define rotation policies
* Consider access scope limitations
* Plan backup procedures

## References
* ADR-009: Key Derivation Architecture
* ADR-002: Secret management model
* ADR-003: Input processing
* ADR-004: Service architecture