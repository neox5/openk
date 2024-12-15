package kms

import (
	"errors"
	"time"
	"unicode"
)

var (
	// ErrUsernameEmpty indicates an empty username was provided
	ErrUsernameEmpty = errors.New("username cannot be empty")

	// ErrUsernameLength indicates username exceeds maximum length
	ErrUsernameLength = errors.New("username exceeds maximum length")

	// ErrUsernameInvalid indicates username contains invalid characters
	ErrUsernameInvalid = errors.New("username contains invalid characters")

	// ErrIterationsInvalid indicates iteration count is below minimum
	ErrIterationsInvalid = errors.New("iterations below minimum value")
)

const (
	// MinIterations defines minimum allowed PBKDF2 iterations
	MinIterations = 100_000

	// MaxUsernameLen defines maximum username length in bytes
	MaxUsernameLen = 255
)

// InitialKeyDerivation represents parameters before storage
type InitialKeyDerivation struct {
	Username   string
	Iterations int
}

// KeyDerivation represents parameters as stored in the backend
type KeyDerivation struct {
	ID         string
	Username   string
	Iterations int
	CreatedAt  time.Time
}

// NewKeyDerivation creates key derivation parameters
func NewKeyDerivation(username string, iterations int) (*InitialKeyDerivation, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}

	if iterations < MinIterations {
		return nil, ErrIterationsInvalid
	}

	return &InitialKeyDerivation{
		Username:   username,
		Iterations: iterations,
	}, nil
}

// validateUsername checks username validity
func validateUsername(username string) error {
	if username == "" {
		return ErrUsernameEmpty
	}

	if len(username) > MaxUsernameLen {
		return ErrUsernameLength
	}

	for _, r := range username {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return ErrUsernameInvalid
		}
	}

	return nil
}
