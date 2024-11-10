package secret

import (
	"time"

	"github.com/google/uuid"
)

// SecretType defines the type of secret using iota for easy management.
type SecretType int

const (
	Password SecretType = iota
	Certificate
	APIKey
	SSHKey
)

// Secret struct represents the data, metadata, and lifecycle information of a secret.
type Secret struct {
	ID           string
	PersistentID string
	Type         SecretType
	Data         FieldMap
	Metadata     FieldMap
	Version      int
	CreatedAt    time.Time
	ExpiresAt    *time.Time
	RevokedAt    *time.Time
	Deleted      bool
}

// NewSecret initializes a new Secret with a UUID and sets default fields.
func NewSecret(secretType SecretType, data map[string]string, metadata map[string]string, expiresAt *time.Time) *Secret {
	return &Secret{
		PersistentID: uuid.NewString(),
		Type:         secretType,
		Data:         *NewFieldMap(),
		Metadata:     *NewFieldMap(),
		Version:      1,
		CreatedAt:    time.Now(),
		ExpiresAt:    expiresAt,
		Deleted:      false,
	}
}

// IsExpired checks if the secret is expired based on ExpiresAt.
func (s *Secret) IsExpired() bool {
	return s.ExpiresAt != nil && time.Now().After(*s.ExpiresAt)
}

// IsRevoked checks if the secret is revoked based on RevokedAt.
func (s *Secret) IsRevoked() bool {
	return s.RevokedAt != nil && time.Now().After(*s.RevokedAt)
}

// IsActive checks if the secret is active, meaning it is not deleted, expired, or revoked.
func (s *Secret) IsActive() bool {
	return !s.Deleted && !s.IsExpired() && !s.IsRevoked()
}

// MarkDeleted sets the Deleted flag to true, marking the secret as deleted.
func (s *Secret) MarkDeleted() {
	s.Deleted = true
}

// Revoke sets the RevokedAt timestamp to the current time, revoking the secret.
func (s *Secret) Revoke() {
	now := time.Now()
	s.RevokedAt = &now
}
