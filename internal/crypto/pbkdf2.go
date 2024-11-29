package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// PBKDF2 parameters
	SaltSize       = 16 // 128 bits
	MasterKeySize  = 32 // 256 bits
	IterationCount = 100_000
)

var (
	ErrInvalidSaltSize = errors.New("salt must be 16 bytes")
	ErrEmptyPassword   = errors.New("password cannot be empty")
)

// GenerateSalt generates a cryptographically secure random salt
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
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
	if len(salt) != SaltSize {
		return nil, ErrInvalidSaltSize
	}

	return pbkdf2.Key(password, salt, iterations, keyLen, sha256.New), nil
}

// DeriveMasterKey derives a master key using predefined parameters
func DeriveMasterKey(password, salt []byte) ([]byte, error) {
	return DeriveKey(password, salt, IterationCount, MasterKeySize)
}
