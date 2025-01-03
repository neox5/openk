# Auth Models - User Management

## Domain Models

### Core Types
```go
// UserCreate represents client-side registration data
type UserCreate struct {
    Username           string
    KeyDerivationID   string    // From initial key derivation
    AuthKey           []byte    // Derived from Master Key
    PublicKey         []byte    // From generated KeyPair
    EncryptedKeyPair  Ciphertext // Private key protected via DEK/Envelope
}

// User represents a stored user with their auth data
type User struct {
    ID                string
    Username          string
    KeyDerivationID   string    // References key derivation params
    AuthKey           []byte    // For future authentication
    PublicKey         []byte
    EncryptedKeyPair  Ciphertext
    CreatedAt         time.Time
}
```

## Storage Models

### SQL (PostgreSQL)
```sql
CREATE TABLE users (
    id              UUID PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE,
    deriv_id        UUID NOT NULL REFERENCES key_derivation_params(id),
    auth_key        BYTEA NOT NULL,
    public_key      BYTEA NOT NULL,
    -- Encrypted key pair components
    key_nonce       BYTEA NOT NULL,
    key_data        BYTEA NOT NULL,
    key_tag         BYTEA NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_deriv_id ON users(deriv_id);
```

### Key-Value (Redis)
```
# User records
user:{id} -> {
    username: string,
    deriv_id: string,
    auth_key: bytes,
    public_key: bytes,
    key_nonce: bytes,
    key_data: bytes,
    key_tag: bytes,
    created_at: timestamp
}

# Indexes
username_to_user:{username} -> id
deriv_to_user:{deriv_id} -> id
```

### Document (MongoDB)
```javascript
// Users Collection
{
    _id: UUID,
    username: String,
    derivId: UUID,      // References key derivation params
    authKey: Binary,    // For future authentication
    publicKey: Binary,
    encryptedKeyPair: {
        nonce: Binary,
        data: Binary,
        tag: Binary
    },
    createdAt: Timestamp
}

// Indexes
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "derivId": 1 });
```

## Notes

1. **Storage Requirements**
   - Username must be unique across the system
   - References to key derivation parameters must be valid
   - All components must be stored atomically
   - No sensitive data storage (follows zero-knowledge principle)

2. **Common Operations**
   - Create new user
   - Check username availability
   - Retrieve user by username or ID
   - List users with key derivation info
