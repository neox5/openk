package kms

import "errors"

var (
	// Common key operation errors
	ErrKeyRevoked   = errors.New("key has been revoked")
	ErrNilEncrypter = errors.New("encrypter cannot be nil")
	ErrNilDecrypter = errors.New("decrypter cannot be nil")

	// DEK specific errors
	ErrInvalidDEK = errors.New("invalid DEK")
	ErrEmptyKey   = errors.New("key data is empty")
)
