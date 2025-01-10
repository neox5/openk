package memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/opene"
	"github.com/neox5/openk/internal/storage/models"
)

const maxUsernameLength = 256 // defines the maximum length as 256 characters per RFC5321

type userMemoryStore struct {
	mu    sync.RWMutex
	users map[string]*models.User // ID -> User
	names map[string]string       // Username -> ID
}

func NewUserMemoryStore() *userMemoryStore {
	return &userMemoryStore{
		users: make(map[string]*models.User),
		names: make(map[string]string),
	}
}

func (s *userMemoryStore) Create(ctx context.Context, input *models.UserCreate) (*models.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if input == nil {
		return nil, opene.NewValidationError("storage", "create_user", "input cannot be nil")
	}

	if input.Username == "" {
		return nil, opene.NewValidationError("storage", "create_user", "username cannot be empty")
	}

	if len(input.Username) > maxUsernameLength {
		return nil, opene.NewValidationError("storage", "create_user", "username too long").
			WithMetadata(opene.Metadata{
				"max_length":    maxUsernameLength,
				"actual_length": len(input.Username),
			})
	}

	// Validate required byte arrays
	if input.Salt == nil || input.AuthKeyHash == nil || input.PublicKey == nil ||
		input.EncryptedKeyPair.Nonce == nil || input.EncryptedKeyPair.Data == nil || input.EncryptedKeyPair.Tag == nil {
		return nil, opene.NewValidationError("storage", "create_user", "all byte arrays are required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.names[input.Username]; exists {
		return nil, opene.NewConflictError("storage", "create_user", "username already exists")
	}

	user := &models.User{
		ID:          uuid.New().String(),
		Username:    input.Username,
		Iterations:  input.Iterations,
		Salt:        make([]byte, len(input.Salt)),
		AuthKeyHash: make([]byte, len(input.AuthKeyHash)),
		PublicKey:   make([]byte, len(input.PublicKey)),
		EncryptedKeyPair: crypto.Ciphertext{
			Nonce: make([]byte, len(input.EncryptedKeyPair.Nonce)),
			Data:  make([]byte, len(input.EncryptedKeyPair.Data)),
			Tag:   make([]byte, len(input.EncryptedKeyPair.Tag)),
		},
		CreatedAt: time.Now(),
	}

	// Deep copy all byte slices
	copy(user.Salt, input.Salt)
	copy(user.AuthKeyHash, input.AuthKeyHash)
	copy(user.PublicKey, input.PublicKey)
	copy(user.EncryptedKeyPair.Nonce, input.EncryptedKeyPair.Nonce)
	copy(user.EncryptedKeyPair.Data, input.EncryptedKeyPair.Data)
	copy(user.EncryptedKeyPair.Tag, input.EncryptedKeyPair.Tag)

	s.users[user.ID] = user
	s.names[user.Username] = user.ID

	return copyUser(user), nil
}

func (s *userMemoryStore) GetByID(ctx context.Context, id string) (*models.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if id == "" {
		return nil, opene.NewValidationError("storage", "get_user", "id cannot be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, opene.NewNotFoundError("storage", "get_user", "user not found").
			WithMetadata(opene.Metadata{
				"id": id,
			})
	}

	return copyUser(user), nil
}

func (s *userMemoryStore) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if username == "" {
		return nil, opene.NewValidationError("storage", "get_user", "username cannot be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	id, exists := s.names[username]
	if !exists {
		return nil, opene.NewNotFoundError("storage", "get_user", "user not found").
			WithMetadata(opene.Metadata{
				"username": username,
			})
	}

	user := s.users[id]
	return copyUser(user), nil
}

// copyUser creates a deep copy of a user ensuring separate memory allocation
func copyUser(u *models.User) *models.User {
	if u == nil {
		return nil
	}

	userCopy := &models.User{
		ID:         u.ID,
		Username:   u.Username,
		Iterations: u.Iterations,
		CreatedAt:  u.CreatedAt,
	}

	// Deep copy byte slices with nil checks
	if len(u.Salt) > 0 {
		userCopy.Salt = make([]byte, len(u.Salt))
		copy(userCopy.Salt, u.Salt)
	}

	if len(u.AuthKeyHash) > 0 {
		userCopy.AuthKeyHash = make([]byte, len(u.AuthKeyHash))
		copy(userCopy.AuthKeyHash, u.AuthKeyHash)
	}

	if len(u.PublicKey) > 0 {
		userCopy.PublicKey = make([]byte, len(u.PublicKey))
		copy(userCopy.PublicKey, u.PublicKey)
	}

	// Deep copy EncryptedKeyPair fields
	if len(u.EncryptedKeyPair.Nonce) > 0 {
		userCopy.EncryptedKeyPair.Nonce = make([]byte, len(u.EncryptedKeyPair.Nonce))
		copy(userCopy.EncryptedKeyPair.Nonce, u.EncryptedKeyPair.Nonce)
	}

	if len(u.EncryptedKeyPair.Data) > 0 {
		userCopy.EncryptedKeyPair.Data = make([]byte, len(u.EncryptedKeyPair.Data))
		copy(userCopy.EncryptedKeyPair.Data, u.EncryptedKeyPair.Data)
	}

	if len(u.EncryptedKeyPair.Tag) > 0 {
		userCopy.EncryptedKeyPair.Tag = make([]byte, len(u.EncryptedKeyPair.Tag))
		copy(userCopy.EncryptedKeyPair.Tag, u.EncryptedKeyPair.Tag)
	}

	return userCopy
}
