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

// KeyPair represents an asymmetric RSA key pair
type KeyPair struct {
    ID         string     
    Algorithm  Algorithm  // AlgorithmRSAOAEPSHA256
    PublicKey  []byte     // X.509/SPKI format
    PrivateKey Ciphertext // Encrypted by DEK
    Created    time.Time
    State      KeyState   
    DEKID      string     // References protecting DEK
}

// DEK (Data Encryption Key) is a symmetric key used for encryption
type DEK struct {
    ID        string    
    Algorithm Algorithm // AlgorithmAESGCM256
    Created   time.Time
    State     KeyState
    // Note: Envelopes managed separately in storage for flexibility
}

// Envelope wraps a DEK using an encryption provider
type Envelope struct {
    ID          string
    DEKID       string     // References wrapped DEK
    Algorithm   Algorithm  // AlgorithmRSAOAEPSHA256 or others
    Key         Ciphertext // DEK encrypted by provider
    Created     time.Time
    State       KeyState
    EncrypterID string    // References encryption provider
}

// KeyDerivation represents parameters for Master Key derivation
type KeyDerivation struct {
    Username     string    // Used as salt
    Iterations   int       // PBKDF2 iteration count
    CreatedAt    time.Time
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
    username    TEXT NOT NULL,
    iterations  INTEGER NOT NULL,
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
    state       SMALLINT NOT NULL DEFAULT 0 CHECK (state >= 0 AND state <= 3),
    dek_id      UUID NOT NULL REFERENCES deks(id)
);

-- Data Encryption Keys
CREATE TABLE deks (
    id          UUID PRIMARY KEY,
    algorithm   SMALLINT NOT NULL DEFAULT 1 CHECK (algorithm >= 0 AND algorithm <= 1),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    state       SMALLINT NOT NULL DEFAULT 0 CHECK (state >= 0 AND state <= 3)
);

-- Key Envelopes
CREATE TABLE envelopes (
    id          UUID PRIMARY KEY,
    dek_id      UUID NOT NULL REFERENCES deks(id),
    algorithm   SMALLINT NOT NULL DEFAULT 0 CHECK (algorithm >= 0 AND algorithm <= 1),
    encrypter_id TEXT NOT NULL,    -- References encryption provider
    nonce       BYTEA NOT NULL,    -- DEK encryption nonce
    data        BYTEA NOT NULL,    -- Encrypted DEK data
    tag         BYTEA NOT NULL,    -- DEK encryption tag
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    state       SMALLINT NOT NULL DEFAULT 0 CHECK (state >= 0 AND state <= 3)
);

-- Indexes
CREATE INDEX idx_key_pairs_state ON key_pairs(state);
CREATE INDEX idx_key_pairs_dek_id ON key_pairs(dek_id);
CREATE INDEX idx_deks_state ON deks(state);
CREATE INDEX idx_envelopes_dek_id ON envelopes(dek_id);
CREATE INDEX idx_envelopes_encrypter_id ON envelopes(encrypter_id);
CREATE INDEX idx_envelopes_state ON envelopes(state);
```

### Key-Value (Redis)

```
# Key Derivation Parameters
kd:{id} -> {
    username: string,
    iterations: number,
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
    state: number,      # 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
    dek_id: string     # Reference to protecting DEK
}

# DEKs
dek:{id} -> {
    algorithm: number,   # 1=AES-256-GCM
    created_at: timestamp,
    state: number       # 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
}

# Envelopes
env:{id} -> {
    dek_id: string,     # Reference to wrapped DEK
    algorithm: number,   # 0=RSA-2048-OAEP-SHA256
    encrypter_id: string, # References encryption provider
    nonce: bytes,       # DEK encryption nonce
    data: bytes,        # Encrypted DEK data
    tag: bytes,         # DEK encryption tag
    created_at: timestamp,
    state: number       # 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
}

# Indexes
dek_envelopes:{dek_id} -> Set[env_id]          # Envelope IDs for DEK
encrypter_envelopes:{encrypter_id} -> Set[env_id]  # Envelope IDs for encrypter
```

### Document (MongoDB)

```javascript
// Key Derivation Parameters Collection
{
    _id: UUID,
    username: String,
    iterations: Number,
    createdAt: Timestamp
}

// Key Pairs Collection
{
    _id: UUID,
    algorithm: Number,  // 0=RSA-2048-OAEP-SHA256
    publicKey: Binary,  // X.509/SPKI format
    private: {
        nonce: Binary,  // Private key encryption 
        data: Binary,
        tag: Binary
    },
    createdAt: Timestamp,
    state: Number,      // 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
    dekId: UUID        // References protecting DEK
}

// DEKs Collection
{
    _id: UUID,
    algorithm: Number,  // 1=AES-256-GCM
    createdAt: Timestamp,
    state: Number,     // 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
}

// Envelopes Collection
{
    _id: UUID,
    dekId: UUID,       // References wrapped DEK
    algorithm: Number, // Encryption algorithm used
    encrypterId: String, // References encryption provider
    key: {
        nonce: Binary,
        data: Binary,
        tag: Binary
    },
    createdAt: Timestamp,
    state: Number     // 0=Active, 1=PendingRotation, 2=Inactive, 3=Destroyed
}

// Indexes
db.keyPairs.createIndex({ "state": 1 });
db.keyPairs.createIndex({ "dekId": 1 });
db.deks.createIndex({ "state": 1 });
db.envelopes.createIndex({ "dekId": 1 });
db.envelopes.createIndex({ "encrypterId": 1 });
db.envelopes.createIndex({ "state": 1 });
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
   - Master Key never stored
   - Only derivation parameters persisted
   - SQL uses SMALLINT with CHECK constraints
   - Redis and MongoDB use numeric values
   - All backends should validate values on write

4. **Memory Protection**
   - Key material cleared after use
   - Master Key cleared after operations
   - Auth Key cleared after session establishment
   - Secure memory wiping when available

5. **Relationship Management**
   - DEK to Envelope is one-to-many
   - KeyPair references protecting DEK
   - Envelope references both DEK and encryption provider
   - Encryption provider ID may reference different entity types