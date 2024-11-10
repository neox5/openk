# Secret

```go
type Secret struct {
    ID         string                 // Unique identifier for the secret.
    Type       string                 // Type of secret (e.g., "password", "certificate", "api_key").
    Data       map[string]string      // Core data of the secret as key-value pairs.
    Metadata   map[string]string      // Additional attributes, such as tags, created-by info, etc.
    Version    int                    // Version number for version control.
    CreatedAt  time.Time              // Timestamp for creation.
    ExpiresAt  *time.Time             // Optional expiry timestamp.
    RevokedAt  *time.Time             // Optional revocation timestamp.
    Status     string                 // Status of the secret, e.g., "active", "expired", "revoked".
}
```
