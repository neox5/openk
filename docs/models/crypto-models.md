# Cryptographic Models

This document defines the domain and storage models based on crypto-spec.md.

## Domain Models

### Core Types
```go
type KeyState int
type Algorithm int

// Ciphertext represents encrypted data with authentication
type Ciphertext struct {
    Nonce []byte // 96 bits
    Data  []byte // Encrypted data
    Tag   []byte // 128 bits
}

// KeyPair represents an asymmetric RSA key pair used for key wrapping
type KeyPair struct {
    ID        string     
    Algorithm Algorithm  // AlgorithmRSAOAEPSHA256
    PublicKey []byte     // X.509/SPKI format
    Private   Ciphertext // Encrypted with KEK
    Created   time.Time
    State     KeyState   
}

// DEK (Data Encryption Key) is a symmetric key that encrypts data
type DEK struct {
    ID        string    
    Algorithm Algorithm // AlgorithmAESGCM256
    Key       []byte    // In memory only, wrapped by KeyPair
    Created   time.Time
    State     KeyState  
}

// Envelope wraps a DEK encrypted with a KeyPair's public key
type Envelope struct {
    ID        string     
    Algorithm Algorithm  // AlgorithmRSAOAEPSHA256
    Key       Ciphertext // DEK encrypted with recipient's public key
    Created   time.Time
    State     KeyState  
    OwnerID   string    // References recipient's KeyPair.ID
}

// KeyDerivation represents parameters for Master Key derivation
type KeyDerivation struct {
    Salt          []byte    // 128-bit random salt
    Iterations    int       // PBKDF2 iteration count
    MasterKeySize int       // 256 bits
}

// HKDFParams represents parameters for KEK derivation
type HKDFParams struct {
    Salt     []byte // Optional, 256 bits if used
    Info     []byte // "OpenK-KEK-v1"
    KeySize  int    // 256 bits
}
```

## Storage Models

### SQL (PostgreSQL)

```sql
-- Algorithms using numeric values from crypto-spec.md
CREATE TYPE crypto_algorithm AS ENUM ('0', '1');
COMMENT ON TYPE crypto_algorithm IS 'Algorithm: 0=RSA-2048-OAEP-SHA256, 1=AES-256-GCM';

-- Key states using numeric values from crypto-spec.md
CREATE TYPE key_state AS ENUM ('0', '1', '2', '3');
COMMENT ON TYPE key_state IS 'KeyState: 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed';

-- Key Derivation Parameters
CREATE TABLE key_derivation_params (
    id          UUID PRIMARY KEY,
    pbkdf2_salt BYTEA NOT NULL,       -- PBKDF2 salt
    iterations  INTEGER NOT NULL,      -- PBKDF2 iterations
    hkdf_salt   BYTEA,                -- Optional HKDF salt
    hkdf_info   BYTEA NOT NULL,       -- HKDF context info
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Key Pairs
CREATE TABLE key_pairs (
    id          UUID PRIMARY KEY,
    algorithm   SMALLINT NOT NULL DEFAULT 0 CHECK (algorithm >= 0 AND algorithm <= 1),
    public_key  BYTEA NOT NULL,
    nonce       BYTEA NOT NULL,    -- Private key encryption nonce
    data        BYTEA NOT NULL,    -- Encrypted private key data
    tag         BYTEA NOT NULL,    -- Private key encryption tag
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    state       SMALLINT NOT NULL DEFAULT 0 CHECK (state >= 0 AND state <= 3)
);

-- Data Encryption Keys (only metadata, key material in memory)
CREATE TABLE deks (
    id          UUID PRIMARY KEY,
    algorithm   SMALLINT NOT NULL DEFAULT 1 CHECK (algorithm >= 0 AND algorithm <= 1),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    state       SMALLINT NOT NULL DEFAULT 0 CHECK (state >= 0 AND state <= 3)
);

-- Key Envelopes
CREATE TABLE envelopes (
    id          UUID PRIMARY KEY,
    algorithm   SMALLINT NOT NULL DEFAULT 0 CHECK (algorithm >= 0 AND algorithm <= 1),
    dek_id      UUID NOT NULL REFERENCES deks(id),
    owner_id    UUID NOT NULL REFERENCES key_pairs(id),
    nonce       BYTEA NOT NULL,    -- DEK encryption nonce
    data        BYTEA NOT NULL,    -- Encrypted DEK data
    tag         BYTEA NOT NULL,    -- DEK encryption tag
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    state       SMALLINT NOT NULL DEFAULT 0 CHECK (state >= 0 AND state <= 3),
    UNIQUE(dek_id, owner_id)
);

-- Indexes
CREATE INDEX idx_key_pairs_state ON key_pairs(state);
CREATE INDEX idx_deks_state ON deks(state);
CREATE INDEX idx_envelopes_dek_id ON envelopes(dek_id);
CREATE INDEX idx_envelopes_owner_id ON envelopes(owner_id);
CREATE INDEX idx_envelopes_state ON envelopes(state);
```

### Key-Value (Redis)

```
# Key Derivation Parameters
kd:{id} -> {
    pbkdf2_salt: bytes,    # PBKDF2 salt
    iterations: number,    # PBKDF2 iterations
    hkdf_salt: bytes,     # Optional HKDF salt
    hkdf_info: bytes,     # HKDF context info
    created_at: timestamp
}

# Key Pairs
kp:{id} -> {
    algorithm: number,   # 0=RSA-2048-OAEP-SHA256
    public_key: bytes,  # X.509/SPKI format
    nonce: bytes,       # Private key encryption nonce
    data: bytes,        # Encrypted private key data
    tag: bytes,         # Private key encryption tag
    created_at: timestamp,
    state: number       # 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
}

# DEKs (metadata only)
dek:{id} -> {
    algorithm: number,   # 1=AES-256-GCM
    created_at: timestamp,
    state: number       # 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
}

# Envelopes
env:{id} -> {
    algorithm: number,   # 0=RSA-2048-OAEP-SHA256
    dek_id: string,     # Reference to DEK
    owner_id: string,   # Reference to KeyPair
    nonce: bytes,       # DEK encryption nonce
    data: bytes,        # Encrypted DEK data
    tag: bytes,         # DEK encryption tag
    created_at: timestamp,
    state: number       # 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
}

# Indexes
dek_envelopes:{dek_id} -> Set[env_id]    # Envelope IDs for DEK
keypair_envelopes:{owner_id} -> Set[env_id]  # Envelope IDs for KeyPair
```

### Document (MongoDB)

```javascript
// Key Derivation Parameters Collection
{
    _id: UUID,
    pbkdf2Salt: Binary,    // PBKDF2 salt
    iterations: Number,    // PBKDF2 iterations
    hkdfSalt: Binary,     // Optional HKDF salt
    hkdfInfo: Binary,     // HKDF context info
    createdAt: Timestamp
}

// Key Pairs Collection
{
    _id: UUID,
    algorithm: Number,  // 0=RSA-2048-OAEP-SHA256
    publicKey: Binary,  // X.509/SPKI format
    private: {
        nonce: Binary,
        data: Binary,
        tag: Binary
    },
    createdAt: Timestamp,
    state: Number      // 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
}

// DEKs Collection (metadata only)
{
    _id: UUID,
    algorithm: Number,  // 1=AES-256-GCM
    createdAt: Timestamp,
    state: Number,     // 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
    // Envelopes embedded for efficient access
    envelopes: [{
        id: UUID,
        algorithm: Number,  // 0=RSA-2048-OAEP-SHA256
        ownerId: UUID,     // References KeyPair._id
        key: {
            nonce: Binary,
            data: Binary,
            tag: Binary
        },
        createdAt: Timestamp,
        state: Number      // 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
    }]
}

// Indexes
db.keyPairs.createIndex({ "state": 1 });
db.deks.createIndex({ "state": 1 });
db.deks.createIndex({ "envelopes.ownerId": 1 });
db.deks.createIndex({ "envelopes.state": 1 });
```

## Notes

1. **Algorithm Values**
   - 0: RSA-2048-OAEP-SHA256 (used for key wrapping)
   - 1: AES-256-GCM (used for data encryption)

2. **State Values**
   - 0: Active 
   - 1: PendingRotation
   - 2: Inactive
   - 3: Destroyed

3. **Storage Considerations**
   - Master Key and KEK never stored
   - Only derivation parameters are persisted
   - SQL uses SMALLINT with CHECK constraints
   - Redis and MongoDB use numeric values
   - All backends should validate values on write

4. **Memory Protection**
   - Key material cleared after use
   - Master Key cleared after KEK derivation
   - KEK cleared after private key encryption
   - Secure memory wiping when available