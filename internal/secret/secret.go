// Package secret provides the core domain models for the secret management system.
package secret

import (
	"time"

	"github.com/google/uuid"
)

// EncryptedSecret represents the storage-level secret with all sensitive data encrypted.
type EncryptedSecret struct {
	// ID uniquely identifies this specific version of the secret
	ID uuid.UUID
	// SecretID identifies the logical secret across all versions
	SecretID uuid.UUID
	// Version is the sequential version number of this secret
	Version uint64
	// Actor identifies who/what performed the action
	Actor uuid.UUID
	// CreatedAt is the timestamp when this version was created
	CreatedAt time.Time
	// ExpiresAt defines when this secret version expires (optional)
	ExpiresAt *time.Time
	// RevokedAt defines when this secret version was revoked (optional)
	RevokedAt *time.Time
	
	// Fields contains the encrypted secret data
	Fields []*EncryptedField
}

