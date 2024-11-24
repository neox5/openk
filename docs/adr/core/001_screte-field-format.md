# ADR-001: Secret Field Payload Format

## Status
Accepted

## Context
Secret fields need to store various data types efficiently with type safety and without schema requirements.

## Considered Options

1. Raw Bytes
   - No type information
   - Direct binary storage
   - Application must know data interpretation

2. JSON with Base64
   - Explicit type handling
   - Base64 overhead
   - Additional conversion required

3. CBOR (RFC 8949)
   - Native binary support
   - Type preservation
   - No schema required
   - Used in COSE

## Decision
Use CBOR for secret field payload serialization.

## Consequences
+ Efficient binary handling
+ Type safety without schema
+ Standard format, good tooling
- Less human-readable than JSON
- Requires CBOR library