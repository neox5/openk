# OpenK - Open Source Secret Management System

## Overview
OpenK ("openKey") is an end-to-end encrypted secret management system designed to provide a unified architecture for both enterprise and personal use cases. The system implements a comprehensive security model with zero-knowledge architecture, ensuring that secrets remain encrypted end-to-end regardless of the use case.

## Core Use Cases

### Enterprise Secret Management
- Zero-knowledge architecture with end-to-end encryption
- Role-based access through envelope encryption
- Flexible authentication through OAuth/OIDC providers and Enterprise SSO
- Support for technical clients with secure key derivation
- Comprehensive audit logging and version control
- Integration with existing enterprise infrastructure

### Personal Password Management
- Unified cryptographic foundation shared with enterprise deployment
- Personal vaults with secure device synchronization
- Secure sharing through envelope encryption
- User-friendly interfaces while maintaining security

## Architecture Components

### Cryptographic Foundation
All cryptographic operations in OpenK are defined and standardized in our crypto-spec. Key features include:
- Standard cryptographic structures (KeyPair, DEK, Envelope)
- Consistent key state management
- Unified approach to key rotation
- Standardized error handling
- Memory protection requirements

For detailed cryptographic specifications, see crypto-spec.md.

### Secret Organization
- **Hierarchical Structure**: 
  - Path-based organization using forward slashes
  - Labels for identity qualification (environments, regions)
  - Optional tags for flexible categorization
- **Unique Identification**: Path + labels + name combination
- **Privacy Preservation**:
  - HMAC-based metadata protection
  - Zero-knowledge search capabilities
  - Private organizational structure

### Service Architecture

#### Organization Service
- Client interaction handling
- Metadata management with privacy preservation
- Input validation
- Coordination of operations

#### Secret Service
- Core secret management
- Version control
- Secret storage operations
- Transaction management

### Key Management

#### Enterprise Deployment
- Organization secrets (org_secret) implemented as DEKs
- Role-based access through envelope encryption
- Key rotation and recovery procedures
- Multi-admin support

#### Personal Vaults
- Personal vault secrets (vault_secret) implemented as DEKs
- Device synchronization support
- Secure sharing capabilities
- Individual recovery procedures

### Security Features

#### Authentication Layer
- Pluggable authentication providers
- Enterprise SSO integration
- Multi-factor authentication support
- Session management

#### Encryption Layer
- Zero-knowledge architecture
- Client-side encryption/decryption
- Standard envelope encryption for sharing
- Key rotation capabilities

#### Privacy Features
- Private metadata through HMAC
- Encrypted path segments
- Protected organizational structure
- Secure search capabilities

## Technical Implementation

### Core Components
- Secret Model: Field-based with privacy preservation
- Storage Backend: Abstract interface for multiple implementations
- Encryption Layer: Standard cryptographic operations
- Key Management: Unified DEK-based approach
- Organization Model: Private path-based structure

### Technology Stack
- Implementation in Go
- CBOR for payload serialization
- Standard cryptographic primitives
- Pluggable authentication and KMS

## Documentation Structure

### Architecture Decision Records
- ADR-001: Secret Field Payload Format
- ADR-002: Secret Management Model
- ADR-003: Secret Input Processing
- ADR-004: Secret Management Service Architecture
- ADR-005: Secret Management Encryption Architecture
- ADR-006: Privacy-Preserving Metadata Model
- ADR-007: Organization Secret Management
- ADR-008: Personal Vault Secret Management

### Technical Specifications
- crypto-spec.md: Cryptographic foundations and requirements
- Additional specifications in development

## Standards Compliance
- Standard cryptographic primitives
- Modern cryptographic practices
- Enterprise security requirements
- Privacy regulations

## Repository
[github.com/neox5/openk](https://github.com/neox5/openk)