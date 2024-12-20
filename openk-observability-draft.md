# Observability in OpenK: Logging, Metrics, and Tracing

## Logging
- **Purpose**: Record discrete events for debugging and audits.
- **Levels**:
  - **Internal Packages**: 
    - **Crypto**: Log encryption/decryption failures.
    - **KMS**: Log key retrievals, rotations, and errors.
  - **Service Layers**: 
    - Organization Service: Metadata access, validation failures.
    - Secret Service: Storage operations, version creation.
  - **Authentication**: Log login attempts, MFA challenges, and session expirations.
- **Example**: `User 'X' failed login due to invalid credentials.`

## Metrics
- **Purpose**: Monitor performance and operational health over time.
- **Levels**:
  - **Internal Packages**:
    - **Crypto**: Encryption/decryption throughput.
    - **KMS**: External key management latency.
  - **Service Layers**:
    - Organization Service: User activity metrics.
    - Secret Service: Storage latency, transaction commit rates.
  - **Authentication**: Success/failure rates, session expirations.
- **Example**: `auth_mfa_failures_total{method="otp"}` (MFA failure counter).

## Tracing
- **Purpose**: Provide end-to-end visibility into request flows.
- **Levels**:
  - **Service-Level Operations**:
    - Organization Service: Trace metadata retrieval and validation paths.
    - Secret Service: Trace secret retrieval, creation, and versioning workflows.
  - **Cross-Cutting Concerns**:
    - Encryption/Decryption: Sub-spans for cryptographic operations.
    - External Calls: Trace KMS and authentication provider latencies.
- **Example**: `GET /secrets/:id` spans include:
  1. User authentication span.
  2. Metadata retrieval span.
  3. Decryption span (child of metadata retrieval).

## Integration Summary

| **Feature**  | **Where**                              | **Example Tooling**   | **Example Use Case**                              |
|--------------|----------------------------------------|------------------------|--------------------------------------------------|
| **Logging**  | Internal packages, service layers      | `slog` or pluggable    | Track user actions, record encryption errors.    |
| **Metrics**  | High-frequency operations              | OpenTelemetry Metrics  | Measure API requests, failure rates, latencies. |
| **Tracing**  | Request flows, external dependencies   | OpenTelemetry Tracing  | Trace secret fetch, KMS interactions.           |

## Recommendations
- **Internal Packages**: Use logs for rare events, metrics for frequent insights, traces for workflows with external dependencies.
- **Service Layers**: Capture user actions (logs), monitor trends (metrics), and trace request flows (traces).
- **Authentication**: Focus on logs for security, metrics for trends, traces for debugging external integrations.
