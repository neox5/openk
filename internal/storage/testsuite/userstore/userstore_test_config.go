package userstore

import (
	"github.com/neox5/openk/internal/crypto"
)

// Test configuration referencing crypto package parameters
const (
	SaltLength        = crypto.DefaultSaltSize // Salt size for key derivation
	AuthKeyHashLength = 32                     // SHA-256 output size
	NonceLength       = crypto.NonceSize       // GCM nonce size
	TagLength         = crypto.TagSize         // GCM authentication tag size
)
