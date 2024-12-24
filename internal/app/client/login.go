package client

import (
	"context"

	"github.com/neox5/openk/internal/opene"
)

// LoginOptions contains parameters for user login
type LoginOptions struct {
	Username string
	Password string
}

// Login authenticates a user and establishes a session
func Login(ctx context.Context, opts LoginOptions) error {
	// TODO: Implement key derivation (PBKDF2)
	// TODO: Retrieve and decrypt KeyPair
	// TODO: Establish session

	return opene.NewInternalError("auth", "login", "not implemented").WithMetadata(
		opene.Metadata{
			"username": opts.Username,
			"password": opts.Password,
		},
	)
}
