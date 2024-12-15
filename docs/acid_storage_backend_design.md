# Design Specification: In-Memory and Extensible Storage Backend with ACID Capabilities

## Overview

This document outlines the design for an in-memory storage backend capable of providing **ACID** guarantees, supporting multiple resource types (e.g., KeyPairs, DEKs), and allowing easy extensibility to other storage backends such as key-value (KV) stores or SQL databases. The backend is optimized for read-heavy workloads.

---

## Objectives

1. **Support for Multiple Resource Types**: Handle diverse data structures like KeyPairs and DEKs while maintaining clean separation of resource-specific logic.
2. **ACID Compliance**:
   - **Atomicity**: Ensure transactions are all-or-nothing.
   - **Consistency**: Validate resources to ensure data integrity.
   - **Isolation**: Prevent conflicts in concurrent transactions.
   - **Durability**: Simulated through snapshots or logs for in-memory, inherent for SQL.
3. **Read Optimization**: Prioritize frequent reads over writes with efficient data structures.
4. **Extensibility**:
   - Abstract interface to support different storage backends.
   - Dynamic resource handling for easy onboarding of new types.

---

## High-Level Design

### 1. **Resource Abstraction**

All resource types (e.g., KeyPairs, DEKs) must implement the following interface:

```go
type Resource interface {
  GetID() string                  // Unique identifier for the resource
  Serialize() ([]byte, error)     // Serialize resource to bytes (e.g., JSON, protobuf)
  Deserialize([]byte) error       // Deserialize bytes into the resource
  Validate() error                // Validate the resource for consistency
}
```

Each resource type encapsulates its own logic for:
- Data serialization/deserialization
- Validation to ensure consistency
- Identifier management (e.g., resource-specific IDs)

---

### 2. **General Storage Interface**

The storage backend exposes a unified interface to abstract different implementations (in-memory, KV store, SQL):

```go
type Storage interface {
  Get(resourceType string, id string) (Resource, error)  // Fetch resource by type and ID
  Put(resourceType string, resource Resource) error     // Insert or update resource
  Delete(resourceType string, id string) error          // Delete resource by ID
  BeginTransaction() Transaction                        // Start a transaction
}
```

#### Transaction Interface

Transactions allow ACID compliance with atomic commits and rollbacks:

```go
type Transaction interface {
  Get(resourceType string, id string) (Resource, error)
  Put(resourceType string, resource Resource) error
  Delete(resourceType string, id string) error
  Commit() error    // Persist all changes
  Rollback() error  // Revert all changes
}
```

---

### 3. **In-Memory Storage Implementation**

#### Structure

Data is stored in a nested map for efficient organization:

```go
type InMemoryStorage struct {
  data map[string]map[string]Resource  // resourceType -> id -> Resource
  lock sync.RWMutex
}
```

#### Operations

- **Get**: Fetch a resource by type and ID.
- **Put**: Add or update a resource, validating it before saving.
- **Delete**: Remove a resource by ID.
- **Transaction Support**: Implement isolation via per-transaction copies of data.

#### Transaction Management

Transactions operate on an isolated map of changes:

```go
type InMemoryTransaction struct {
  storage   *InMemoryStorage
  changes   map[string]map[string]Resource  // Staged changes
  deleted   map[string]map[string]bool      // Staged deletions
  committed bool
}
```

---

### 4. **SQL-Based Storage Implementation**

#### Structure

Each resource type maps to a dedicated SQL table. Table names are dynamically derived from resource types:

```sql
CREATE TABLE KeyPair (
  id TEXT PRIMARY KEY,
  data BLOB NOT NULL
);

CREATE TABLE DEK (
  id TEXT PRIMARY KEY,
  data BLOB NOT NULL
);
```

#### Operations

- **Get**: Fetch the serialized resource by ID and deserialize it.
- **Put**: Serialize the resource and perform an `INSERT` or `UPDATE`.
- **Delete**: Remove the resource by ID.

#### Dynamic Resource Handling

The SQL storage backend uses helper functions to:
- Map resource types to table names.
- Create resource instances dynamically based on type.

---

### 5. **Resource-Specific Structs**

Each resource type implements the `Resource` interface. Example structures:

#### KeyPair
```go
type KeyPair struct {
  ID         string
  PublicKey  string
  PrivateKey string
}

func (k *KeyPair) GetID() string { return k.ID }

func (k *KeyPair) Serialize() ([]byte, error) {
  return json.Marshal(k)
}

func (k *KeyPair) Deserialize(data []byte) error {
 # Design Specification: In-Memory and Extensible Storage Backend with ACID Capabilities

## Overview

This document outlines the design for an in-memory storage backend capable of providing **ACID** guarantees, supporting multiple resource types (e.g., KeyPairs, DEKs), and allowing easy extensibility to other storage backends such as key-value (KV) stores or SQL databases. The backend is optimized for read-heavy workloads.

---

## Objectives

1. **Support for Multiple Resource Types**: Handle diverse data structures like KeyPairs and DEKs while maintaining clean separation of resource-specific logic.
2. **ACID Compliance**:
   - **Atomicity**: Ensure transactions are all-or-nothing.
   - **Consistency**: Validate resources to ensure data integrity.
   - **Isolation**: Prevent conflicts in concurrent transactions.
   - **Durability**: Simulated through snapshots or logs for in-memory, inherent for SQL.
3. **Read Optimization**: Prioritize frequent reads over writes with efficient data structures.
4. **Extensibility**:
   - Abstract interface to support different storage backends.
   - Dynamic resource handling for easy onboarding of new types.

---

## High-Level Design

### 1. **Resource Abstraction**

All resource types (e.g., KeyPairs, DEKs) must implement the following interface:

```go
type Resource interface {
  GetID() string                  // Unique identifier for the resource
  Serialize() ([]byte, error)     // Serialize resource to bytes (e.g., JSON, protobuf)
  Deserialize([]byte) error       // Deserialize bytes into the resource
  Validate() error                // Validate the resource for consistency
}
```

Each resource type encapsulates its own logic for:
- Data serialization/deserialization
- Validation to ensure consistency
- Identifier management (e.g., resource-specific IDs)

---

### 2. **General Storage Interface**

The storage backend exposes a unified interface to abstract different implementations (in-memory, KV store, SQL):

```go
type Storage interface {
  Get(resourceType string, id string) (Resource, error)  // Fetch resource by type and ID
  Put(resourceType string, resource Resource) error     // Insert or update resource
  Delete(resourceType string, id string) error          // Delete resource by ID
  BeginTransaction() Transaction                        // Start a transaction
}
```

#### Transaction Interface

Transactions allow ACID compliance with atomic commits and rollbacks:

```go
type Transaction interface {
  Get(resourceType string, id string) (Resource, error)
  Put(resourceType string, resource Resource) error
  Delete(resourceType string, id string) error
  Commit() error    // Persist all changes
  Rollback() error  // Revert all changes
}
```

---

### 3. **In-Memory Storage Implementation**

#### Structure

Data is stored in a nested map for efficient organization:

```go
type InMemoryStorage struct {
  data map[string]map[string]Resource  // resourceType -> id -> Resource
  lock sync.RWMutex
}
```

#### Operations

- **Get**: Fetch a resource by type and ID.
- **Put**: Add or update a resource, validating it before saving.
- **Delete**: Remove a resource by ID.
- **Transaction Support**: Implement isolation via per-transaction copies of data.

#### Transaction Management

Transactions operate on an isolated map of changes:

```go
type InMemoryTransaction struct {
  storage   *InMemoryStorage
  changes   map[string]map[string]Resource  // Staged changes
  deleted   map[string]map[string]bool      // Staged deletions
  committed bool
}
```

---

### 4. **SQL-Based Storage Implementation**

#### Structure

Each resource type maps to a dedicated SQL table. Table names are dynamically derived from resource types:

```sql
CREATE TABLE KeyPair (
  id TEXT PRIMARY KEY,
  data BLOB NOT NULL
);

CREATE TABLE DEK (
  id TEXT PRIMARY KEY,
  data BLOB NOT NULL
);
```

#### Operations

- **Get**: Fetch the serialized resource by ID and deserialize it.
- **Put**: Serialize the resource and perform an `INSERT` or `UPDATE`.
- **Delete**: Remove the resource by ID.

#### Dynamic Resource Handling

The SQL storage backend uses helper functions to:
- Map resource types to table names.
- Create resource instances dynamically based on type.

---

### 5. **Resource-Specific Structs**

Each resource type implements the `Resource` interface. Example structures:

#### KeyPair
```go
type KeyPair struct {
  ID         string
  PublicKey  string
  PrivateKey string
}

func (k *KeyPair) GetID() string { return k.ID }

func (k *KeyPair) Serialize() ([]byte, error) {
  return json.Marshal(k)
}

func (k *KeyPair) Deserialize(data []byte) error {
  return json.Unmarshal(data, k)
}

func (k *KeyPair) Validate() error {
  if k.ID == "" || k.PublicKey == "" || k.PrivateKey == "" {
    return fmt.Errorf("invalid KeyPair")
  }
  return nil
}
```

#### DEK
```go
type DEK struct {
  ID     string
  Key    string
  Expiry time.Time
}

func (d *DEK) GetID() string { return d.ID }

func (d *DEK) Serialize() ([]byte, error) {
  return json.Marshal(d)
}

func (d *DEK) Deserialize(data []byte) error {
  return json.Unmarshal(data, d)
}

func (d *DEK) Validate() error {
  if d.ID == "" || d.Key == "" {
    return fmt.Errorf("invalid DEK")
  }
  if time.Now().After(d.Expiry) {
    return fmt.Errorf("DEK expired")
  }
  return nil
}
```

---

## Features

### ACID Guarantees

1. **Atomicity**: Transactions either commit all changes or roll back completely.
2. **Consistency**: Validate resources before committing.
3. **Isolation**: Each transaction operates on isolated copies of data.
4. **Durability**:
   - **In-Memory**: Optional periodic snapshots or write-ahead logs (WAL).
   - **SQL**: Native durability through database persistence.

### Extensibility

1. **Backend Flexibility**: Swap storage backends (in-memory, KV store, SQL) by implementing the `Storage` interface.
2. **Dynamic Resource Support**: Register new resource types dynamically via factories or registries.

---

## Summary of Components

| **Component**         | **Responsibility**                                                                 |
|------------------------|-------------------------------------------------------------------------------------|
| `Resource` Interface   | Defines operations for resource types (e.g., serialization, validation).           |
| `Storage` Interface    | Abstracts storage operations for in-memory, KV, or SQL backends.                   |
| In-Memory Storage      | Read-optimized, thread-safe, transactional, nested-map-based implementation.       |
| SQL Storage            | Maps resource types to tables, dynamically handles schema and operations.          |
| Transaction Interface  | Provides ACID guarantees through staged changes, isolation, commit, and rollback.  |
| Resource Implementations | Specific logic for each resource type (e.g., KeyPair, DEK).                        |

---

## Future Enhancements

1. **Indexing**: Add indexing support for faster queries in in-memory and KV backends.
2. **Snapshotting**: Implement periodic persistence for in-memory storage durability.
3. **Event Logging**: Provide a mechanism to log changes for debugging or auditing.
4. **Configuration API**: Enable runtime configuration of resource types and backends.

---

This design ensures a robust, extensible, and high-performance storage backend for various resource types with ACID guarantees. return json.Unmarshal(data, k)
}

func (k *KeyPair) Validate() error {
  if k.ID == "" || k.PublicKey == "" || k.PrivateKey == "" {
    return fmt.Errorf("invalid KeyPair")
  }
  return nil
}
```

#### DEK
```go
type DEK struct {
  ID     string
  Key    string
  Expiry time.Time
}

func (d *DEK) GetID() string { return d.ID }

func (d *DEK) Serialize() ([]byte, error) {
  return json.Marshal(d)
}

func (d *DEK) Deserialize(data []byte) error {
  return json.Unmarshal(data, d)
}

func (d *DEK) Validate() error {
  if d.ID == "" || d.Key == "" {
    return fmt.Errorf("invalid DEK")
  }
  if time.Now().After(d.Expiry) {
    return fmt.Errorf("DEK expired")
  }
  return nil
}
```

---

## Features

### ACID Guarantees

1. **Atomicity**: Transactions either commit all changes or roll back completely.
2. **Consistency**: Validate resources before committing.
3. **Isolation**: Each transaction operates on isolated copies of data.
4. **Durability**:
   - **In-Memory**: Optional periodic snapshots or write-ahead logs (WAL).
   - **SQL**: Native durability through database persistence.

### Extensibility

1. **Backend Flexibility**: Swap storage backends (in-memory, KV store, SQL) by implementing the `Storage` interface.
2. **Dynamic Resource Support**: Register new resource types dynamically via factories or registries.

---

## Summary of Components

| **Component**         | **Responsibility**                                                                 |
|------------------------|-------------------------------------------------------------------------------------|
| `Resource` Interface   | Defines operations for resource types (e.g., serialization, validation).           |
| `Storage` Interface    | Abstracts storage operations for in-memory, KV, or SQL backends.                   |
| In-Memory Storage      | Read-optimized, thread-safe, transactional, nested-map-based implementation.       |
| SQL Storage            | Maps resource types to tables, dynamically handles schema and operations.          |
| Transaction Interface  | Provides ACID guarantees through staged changes, isolation, commit, and rollback.  |
| Resource Implementations | Specific logic for each resource type (e.g., KeyPair, DEK).                        |

---

## Future Enhancements

1. **Indexing**: Add indexing support for faster queries in in-memory and KV backends.
2. **Snapshotting**: Implement periodic persistence for in-memory storage durability.
3. **Event Logging**: Provide a mechanism to log changes for debugging or auditing.
4. **Configuration API**: Enable runtime configuration of resource types and backends.

---

This design ensures a robust, extensible, and high-performance storage backend for various resource types with ACID guarantees.
