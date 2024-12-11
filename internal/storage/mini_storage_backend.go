package storage

import (
	"context"

	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/kms"
	"github.com/neox5/openk/internal/secret"
)

type MiniStorageBackend interface {
	// Key Derivation Parameters
	StoreDerivationParams(ctx context.Context, params *kms.KeyDerivation) error
	GetDerivationParams(ctx context.Context, username string) (*kms.KeyDerivation, error)

	// KeyPair operations
	StoreKeyPair(ctx context.Context, kp *kms.InitialKeyPair) (*kms.KeyPair, error)
	GetKeyPairByID(ctx context.Context, id string) (*kms.KeyPair, error)
	UpdateKeyPairState(ctx context.Context, id string, state crypto.KeyState) error

	// DEK operations
	StoreDEK(ctx context.Context, dek *kms.InitialDEK) (*kms.DEK, error)
	GetDEKByID(ctx context.Context, id string) (*kms.DEK, error)
	UpdateDEKState(ctx context.Context, id string, state crypto.KeyState) error

	// Envelope operations
	StoreEnvelope(ctx context.Context, dekID string, env *kms.InitialEnvelope) (*kms.Envelope, error)
	GetEnvelopesByDEK(ctx context.Context, dekID string) ([]*kms.Envelope, error)
	UpdateEnvelopeState(ctx context.Context, id string, state crypto.KeyState) error

	// Mini Secret operations
	StoreMiniSecret(ctx context.Context, secret *secret.InitialMiniSecret) (*secret.MiniSecret, error)
	GetMiniSecretByID(ctx context.Context, id string) (*secret.MiniSecret, error)
	DeleteMiniSecret(ctx context.Context, id string) error

	// Transaction support
	BeginTx(ctx context.Context) (Tx, error)
}

type Tx interface {
	MiniStorageBackend
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
