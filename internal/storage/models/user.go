package models

import (
	"time"

	"github.com/neox5/openk/internal/crypto"
)

// User represents a stored user
type User struct {
	ID               string
	Username         string
	Iterations       int
	Salt             []byte
	AuthKeyHash      []byte
	PublicKey        []byte
	EncryptedKeyPair crypto.Ciphertext
	CreatedAt        time.Time
}

// UserCreate represents user creation parameters
type UserCreate struct {
	Username         string
	Iterations       int
	Salt             []byte
	AuthKeyHash      []byte
	PublicKey        []byte
	EncryptedKeyPair crypto.Ciphertext
}
