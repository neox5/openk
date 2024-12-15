package storage

import (
	"context"

	"github.com/neox5/openk/internal/kms"
)

type MiniStorageBackend interface {
    // Key Derivation Parameters
    StoreDerivationParams(ctx context.Context, params *kms.InitialKeyDerivation) error
    GetDerivationParams(ctx context.Context, username string) (*kms.KeyDerivation, error)
}

