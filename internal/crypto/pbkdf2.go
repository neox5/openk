package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// PBKDF2 parameters
	DefaultSaltSize = 16 // 128 bits (for GenerateSalt only)
	MasterKeySize   = 32 // 256 bits
	IterationCount  = 100_000
	MinSaltLength   = 1 // Minimum allowed salt length
)

var (
	ErrEmptyPassword = errors.New("password cannot be empty")
	ErrInvalidSalt   = errors.New("salt cannot be empty")
)

// GenerateSalt generates a cryptographically secure random salt
// This is used when we need a random salt (not for username-based salts)
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, DefaultSaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKey derives a key using PBKDF2 with SHA256
func DeriveKey(password, salt []byte, iterations, keyLen int) ([]byte, error) {
	if len(password) == 0 {
		return nil, ErrEmptyPassword
	}
	if len(salt) < MinSaltLength {
		return nil, ErrInvalidSalt
	}

	return pbkdf2.Key(password, salt, iterations, keyLen, sha256.New), nil
}

// DeriveMasterKey derives a master key using predefined parameters
func DeriveMasterKey(password, salt []byte) ([]byte, error) {
	return DeriveKey(password, salt, IterationCount, MasterKeySize)
}
