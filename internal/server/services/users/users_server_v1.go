package users

import (
	"context"
	"crypto/sha256"
	"log/slog"

	usersv1 "github.com/neox5/openk/internal/api_gen/openk/users/v1"
	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/opene"
	"github.com/neox5/openk/internal/storage"
	"github.com/neox5/openk/internal/storage/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UsersServerV1 implements the Users service.
type UsersServerV1 struct {
	usersv1.UnimplementedUserServiceServer
	store  storage.UserStore
	logger *slog.Logger
}

// NewUsersServerV1 creates a new Users service instance.
func NewUsersServerV1(store storage.UserStore, logger *slog.Logger) *UsersServerV1 {
	if logger == nil {
		logger = slog.Default()
	}
	return &UsersServerV1{
		store:  store,
		logger: logger,
	}
}

// CreateUser implements the user creation RPC endpoint.
func (s *UsersServerV1) RegisterUser(ctx context.Context, req *usersv1.RegisterUserRequest) (*usersv1.RegisterUserResponse, error) {
	// Log request with context
	s.logger.LogAttrs(ctx, slog.LevelInfo, "handling user creation",
		slog.String("username", req.Username),
		slog.Int("iterations", int(req.Iterations)),
	)

	// Validate request
	if err := s.validateRegisterRequest(req); err != nil {
		s.logger.LogAttrs(ctx, slog.LevelError, "invalid creation request",
			slog.String("error", err.Error()),
		)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	encryptedKeyPair, err := crypto.NewCiphertext(
		req.EncryptedKeypair.Tag,
		req.EncryptedKeypair.Data,
		req.EncryptedKeypair.Tag,
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid encrypted keypair format")
	}

	// Create user input from request
	input := &models.UserCreate{
		Username:         req.Username,
		Iterations:       int(req.Iterations),
		Salt:            req.Salt,
		AuthKeyHash:      hashAuthKey(req.AuthKey),
		PublicKey:        req.PublicKey,
		EncryptedKeyPair: *encryptedKeyPair,
	}

	// Attempt user creation
	user, err := s.store.Create(ctx, input)
	if err != nil {
		// Map storage errors to appropriate gRPC status codes
		if e, ok := err.(*opene.Error); ok {
			switch e.Code {
			case opene.CodeConflict:
				return nil, status.Error(codes.AlreadyExists, "username already exists")
			case opene.CodeValidation:
				return nil, status.Error(codes.InvalidArgument, e.Message)
			default:
				s.logger.LogAttrs(ctx, slog.LevelError, "user creation failed",
					slog.String("error", err.Error()),
					slog.String("code", string(e.Code)),
				)
				return nil, status.Error(codes.Internal, "internal error")
			}
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	// Log successful creation
	s.logger.LogAttrs(ctx, slog.LevelInfo, "user created successfully",
		slog.String("user_id", user.ID),
		slog.String("username", user.Username),
	)

	// Create response
	return &usersv1.RegisterUserResponse{
		Identity: &usersv1.UserIdentity{
			Id:        user.ID,
			Username:  user.Username,
			PublicKey: user.PublicKey,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}

func (s *UsersServerV1) validateRegisterRequest(req *usersv1.RegisterUserRequest) error {
	if req == nil {
		return opene.NewValidationError("users", "validate_create", "request cannot be nil")
	}

	// Username validation
	if req.Username == "" {
		return opene.NewValidationError("users", "validate_create", "username cannot be empty")
	}

	// Check username length (from auth-models.md)
	if len(req.Username) > 256 {
		return opene.NewValidationError("users", "validate_create", "username too long").
			WithMetadata(opene.Metadata{
				"max_length": 256,
				"actual_length": len(req.Username),
			})
	}

	// Iteration count validation (from key_derivation_architecture.md)
	if req.Iterations < 100_000 {
		return opene.NewValidationError("users", "validate_create", "iterations count too low").
			WithMetadata(opene.Metadata{
				"min_iterations": 100_000,
				"actual": req.Iterations,
			})
	}

	// Cryptographic parameter validation
	if len(req.Salt) == 0 {
		return opene.NewValidationError("users", "validate_create", "salt is required")
	}

	if len(req.AuthKey) == 0 {
		return opene.NewValidationError("users", "validate_create", "auth key is required")
	}

	if len(req.PublicKey) == 0 {
		return opene.NewValidationError("users", "validate_create", "public key is required")
	}

	if req.EncryptedKeypair == nil {
		return opene.NewValidationError("users", "validate_create", "encrypted keypair is required")
	}

	// Validate encrypted keypair structure
	if len(req.EncryptedKeypair.Nonce) == 0 {
		return opene.NewValidationError("users", "validate_create", "keypair nonce is required")
	}

	if len(req.EncryptedKeypair.Data) == 0 {
		return opene.NewValidationError("users", "validate_create", "keypair data is required")
	}

	if len(req.EncryptedKeypair.Tag) == 0 {
		return opene.NewValidationError("users", "validate_create", "keypair tag is required")
	}

	return nil
}

// hashAuthKey creates a SHA-256 hash of the auth key
func hashAuthKey(authKey []byte) []byte {
	hasher := sha256.New()
	hasher.Write(authKey)
	return hasher.Sum(nil)
}
