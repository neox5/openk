# OpenK - Open Source Secret Management System

## Overview
OpenK ("openKey") is an end-to-end encrypted secret management system for enterprise and personal use. The system implements a zero-knowledge architecture with end-to-end encryption.

## Core Use Cases

### Enterprise Secret Management
- Zero-knowledge architecture with end-to-end encryption
- Role-based access through envelope encryption
- Flexible authentication through OAuth/OIDC providers and Enterprise SSO
- Support for technical clients with secure key derivation
- Comprehensive audit logging and version control

### Personal Password Management
- Unified cryptographic foundation shared with enterprise deployment
- Personal vaults with secure device synchronization 
- Secure sharing through envelope encryption
- User-friendly interfaces while maintaining security

## Architecture Components

### Secret Organization
| Component | Implementation |
|-----------|---------------|
| Path Structure | Hierarchical with forward slashes |
| Unique ID | path + labels + name |
| Privacy | HMAC-based metadata protection |

Labels provide identity qualification (environments, regions) while optional tags allow flexible categorization.

### Service Architecture
The system is split into two main services:

1. Organization Service
   - Client interaction handling
   - Metadata management with privacy preservation
   - Input validation
   - Coordination of operations

2. Secret Service
   - Core secret management
   - Version control
   - Secret storage operations
   - Transaction management

### Security Features
The system implements comprehensive security through multiple layers:

1. Authentication Layer
   - Pluggable authentication providers
   - Enterprise SSO integration
   - Multi-factor authentication support
   - Session management

2. Encryption Layer
   - Zero-knowledge architecture
   - Client-side encryption/decryption
   - Standard envelope encryption
   - Key rotation capabilities

3. Privacy Features
   - Private metadata through HMAC
   - Encrypted path segments
   - Protected organizational structure
   - Secure search capabilities

## Technical Implementation
- Implementation in Go
- CBOR for payload serialization
- Standard cryptographic primitives
- Pluggable authentication and storage

## Documentation Structure
All major decisions and specifications are documented through ADRs and technical specifications, covering topics from secret field formats to encryption architecture.

## Standards Compliance
- Standard cryptographic primitives
- Modern cryptographic practices
- Enterprise security requirements
- Privacy regulations

## Repository
[github.com/neox5/openk](https://github.com/neox5/openk)