package storage

import (
	"context"

	"github.com/neox5/openk/internal/storage/models"
)

// Store provides access to all storage operations
type Store interface {
    Users() UserStore
}

// UserStore defines user storage operations
type UserStore interface {
    Create(ctx context.Context, input *models.UserCreate) (*models.User, error)
    GetByID(ctx context.Context, id string) (*models.User, error)
    GetByUsername(ctx context.Context, username string) (*models.User, error)
}

