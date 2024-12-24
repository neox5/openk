package client

import (
	"context"

	"github.com/neox5/openk/internal/opene"
)

// RegisterOptions contains parameters for user registration
type RegisterOptions struct {
	Username string
	Password string
}

// Register creates a new user account
func Register(ctx context.Context, opts RegisterOptions) error {
	// TODO: Implement key derivation (PBKDF2)
	// TODO: Generate initial KeyPair
	// TODO: Create envelope
	// TODO: Send registration request

	return opene.NewInternalError("auth", "register", "not implemented").WithMetadata(
		opene.Metadata{
			"username": opts.Username,
			"password": opts.Password,
		},
	)
}
