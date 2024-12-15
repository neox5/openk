package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/neox5/openk/internal/kms"
)

var (
	// Storage errors
	ErrParamsNil        = errors.New("derivation parameters cannot be nil")
	ErrParamsNotFound   = errors.New("derivation parameters not found")
	ErrUsernameNotFound = errors.New("username not found")
)

// InMemoryMiniStorage implements MiniStorageBackend with in-memory storage
type InMemoryMiniStorage struct {
	// Mutex for thread-safe access to maps
	mu sync.RWMutex

	// Map of username to key derivation parameters
	derivationParams map[string]*kms.KeyDerivation
}

// NewInMemoryMiniStorage creates a new in-memory storage instance
func NewInMemoryMiniStorage() *InMemoryMiniStorage {
	return &InMemoryMiniStorage{
		derivationParams: make(map[string]*kms.KeyDerivation),
	}
}

// StoreDerivationParams stores key derivation parameters
func (s *InMemoryMiniStorage) StoreDerivationParams(ctx context.Context, params *kms.InitialKeyDerivation) error {
	// Input validation
	if params == nil {
		return ErrParamsNil
	}
	if params.Username == "" {
		return kms.ErrUsernameEmpty
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Generate UUID for storage
	id := uuid.New()

	// Store parameters with write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert InitialKeyDerivation to KeyDerivation for storage
	stored := &kms.KeyDerivation{
		ID:         id.String(),
		Username:   params.Username,
		Iterations: params.Iterations,
		CreatedAt:  params.CreatedAt,
	}
	s.derivationParams[params.Username] = stored

	return nil
}

// GetDerivationParams retrieves key derivation parameters for a username
func (s *InMemoryMiniStorage) GetDerivationParams(ctx context.Context, username string) (*kms.KeyDerivation, error) {
	// Input validation
	if username == "" {
		return nil, kms.ErrUsernameEmpty
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Retrieve parameters with read lock
	s.mu.RLock()
	defer s.mu.RUnlock()

	params, exists := s.derivationParams[username]
	if !exists {
		return nil, ErrParamsNotFound
	}

	// Return copy to prevent external modification
	return &kms.KeyDerivation{
		ID:         params.ID,
		Username:   params.Username,
		Iterations: params.Iterations,
		CreatedAt:  params.CreatedAt,
	}, nil
}
