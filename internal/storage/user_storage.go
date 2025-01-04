package storage

import (
	"context"
	"time"

	"github.com/neox5/openk/internal/crypto"
)

// User represents a stored user with their key derivation parameters
type User struct {
	ID              string
	Username        string
	Iterations      int       // PBKDF2 iteration count
	Salt           []byte    // Salt for key derivation
	AuthKeyHash    []byte    // Pre-computed hash of auth key
	PublicKey      []byte    // User's public key
	EncryptedKeyPair crypto.Ciphertext // Protected private key
	CreatedAt      time.Time
}

// UserCreate represents the data required to store a new user
type UserCreate struct {
	Username        string
	Iterations      int
	Salt           []byte
	AuthKeyHash    []byte    // Pre-computed hash of auth key
	PublicKey      []byte
	EncryptedKeyPair crypto.Ciphertext
}

// UserStore defines the interface for user storage operations
type UserStore interface {
	// CreateUser stores a new user record
	CreateUser(ctx context.Context, input *UserCreate) (*User, error)

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*User, error)

	// GetUserByUsername retrieves a user by username
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}
