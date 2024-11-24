// Package storage defines the core storage interface for the secret management system.
package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/neox5/openk/internal/secret"
)

// Common errors that may be returned by storage implementations
var (
	// ErrNotFound indicates the requested secret or version doesn't exist
	ErrNotFound = errors.New("secret not found")

	// ErrVersionNotFound indicates the specific version doesn't exist
	ErrVersionNotFound = errors.New("version not found")

	// ErrConcurrentMod indicates a concurrent modification conflict
	ErrConcurrentMod = errors.New("concurrent modification")

	// ErrInvalidOperation indicates an invalid operation attempt
	ErrInvalidOperation = errors.New("invalid operation")
)

// StorageBackend defines the interface that any storage implementation must satisfy.
// It provides strongly consistent operations for secret management with built-in
// versioning support and transactional capabilities.
type StorageBackend interface {
	// Put stores an encrypted secret, handling versioning internally.
	// For new secrets, it generates and assigns both id and secretId.
	// For existing secrets (identified by secretId), it creates a new version.
	// Returns the stored encrypted secret with its assigned identifiers and version.
	Put(ctx context.Context, secret *secret.EncryptedSecret) (*secret.EncryptedSecret, error)

	// Get retrieves a specific or latest version of an encrypted secret.
	// If version is 0, returns the latest version.
	// Returns ErrNotFound if the secret or version doesn't exist.
	Get(ctx context.Context, secretID uuid.UUID, version uint64) (*secret.EncryptedSecret, error)

	// GetAll retrieves all versions of an encrypted secret.
	// Returns empty slice if no versions exist.
	GetAll(ctx context.Context, secretID uuid.UUID) ([]*secret.EncryptedSecret, error)

	// GetByID retrieves an encrypted secret using its specific instance ID.
	// Returns ErrNotFound if the secret doesn't exist.
	GetByID(ctx context.Context, id uuid.UUID) (*secret.EncryptedSecret, error)

	// Delete marks a secret and all its versions as logically deleted.
	// The data remains recoverable until destroyed.
	Delete(ctx context.Context, secretID uuid.UUID) error

	// Destroy permanently removes a secret and all its versions.
	// This operation is irreversible.
	Destroy(ctx context.Context, secretID uuid.UUID) error

	// Transaction executes the given function within a transaction.
	// The transaction commits if the function returns nil, otherwise rolls back.
	// The provided StorageBackend instance is transaction-aware.
	Transaction(ctx context.Context, fn func(tx StorageBackend) error) error
}
