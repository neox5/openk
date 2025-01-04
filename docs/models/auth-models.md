# Auth Models - User Management

## Domain Models

### Core Types
```go
// UserCreate represents client-side registration data
type UserCreate struct {
    Username         string
    Iterations      int       // PBKDF2 iteration count
    Salt           []byte    // Salt for key derivation
    AuthKey        []byte    // Derived from Master Key, transmitted to service
    PublicKey      []byte    // From generated KeyPair
    EncryptedKeyPair crypto.Ciphertext // Private key protected via DEK/Envelope
}

// User represents a stored user with their auth data
type User struct {
    ID              string
    Username        string
    Iterations      int       // PBKDF2 iteration count
    Salt           []byte    // Salt for key derivation
    AuthKeyHash    []byte    // SHA256(AuthKey)
    PublicKey      []byte
    EncryptedKeyPair crypto.Ciphertext
    CreatedAt      time.Time
}
```

## Storage Models

### SQL (PostgreSQL)
```sql
CREATE TABLE users (
    id              UUID PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE,
    iterations      INTEGER NOT NULL,
    salt            BYTEA NOT NULL,
    auth_key_hash   BYTEA NOT NULL,
    public_key      BYTEA NOT NULL,
    -- Encrypted key pair components
    key_nonce       BYTEA NOT NULL,
    key_data        BYTEA NOT NULL,
    key_tag         BYTEA NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_username ON users(username);
```

### Key-Value (Redis)
```
# User records
user:{id} -> {
    username: string,
    iterations: number,
    salt: bytes,
    auth_key_hash: bytes,
    public_key: bytes,
    key_nonce: bytes,
    key_data: bytes,
    key_tag: bytes,
    created_at: timestamp
}

# Indexes
username_to_user:{username} -> id
```

### Document (MongoDB)
```javascript
// Users Collection
{
    _id: UUID,
    username: String,
    iterations: Number,
    salt: Binary,
    authKeyHash: Binary,    // SHA256(AuthKey)
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
```

## Notes

1. **Storage Requirements**
   - Username must be unique across the system
   - Key derivation parameters stored with user
   - All components must be stored atomically
   - No sensitive data storage (follows zero-knowledge principle)

2. **Common Operations**
   - Create new user
   - Check username availability
   - Retrieve user by username or ID