// Package secret provides the core domain models for the secret management system.
package secret

import (
	"errors"

	"github.com/google/uuid"
)

// Field-specific errors.
var (
	ErrFieldIDRequired = errors.New("field ID is required")
	ErrEmptyPayload    = errors.New("field payload cannot be empty")
	ErrEmptyIV         = errors.New("field IV cannot be empty")
	ErrEmptyHash       = errors.New("field hash is required")
)

// EncryptedField represents an encrypted value with its encryption metadata.
type EncryptedField struct {
	// ID uniquely identifies this field
	ID uuid.UUID
	// Payload contains the encrypted field data
	Payload []byte
	// IV is the initialization vector used for encryption
	IV []byte
	// Hash is the hash of the unencrypted payload
	Hash string
}

// Validate checks if the field has all required attributes.
func (f *EncryptedField) Validate() error {
	if f.ID == uuid.Nil {
		return ErrFieldIDRequired
	}

	if len(f.Payload) == 0 {
		return ErrEmptyPayload
	}

	if len(f.IV) == 0 {
		return ErrEmptyIV
	}

	if f.Hash == "" {
		return ErrEmptyHash
	}

	return nil
}
