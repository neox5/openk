package secret

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/neox5/openk/internal/crypto"
)

var (
	ErrNilEncryption       = errors.New("encryption provider cannot be nil")
	ErrNilDecryption       = errors.New("decryption provider cannot be nil")
	ErrEmptyData           = errors.New("key or value cannot be empty")
	ErrDecrypterIDMismatch = errors.New("decrypter ID mismatch")
)

// InitialMiniSecret represents a newly created secret before storage
type InitialMiniSecret struct {
	Name        string             // Reference name for the secret
	Key         *crypto.Ciphertext // Encrypted key
	Value       *crypto.Ciphertext // Encrypted value
	EncrypterID string             // ID of the encryption provider (e.g., DEK ID)
}

// MiniSecret represents a secret as stored in the backend
type MiniSecret struct {
	ID          string             // Unique identifier
	Name        string             // Reference name for the secret
	Key         *crypto.Ciphertext // Encrypted key
	Value       *crypto.Ciphertext // Encrypted value
	EncrypterID string             // ID of the encryption provider (e.g., DEK ID)
}

// UnsealedMiniSecret represents an active secret with access to decrypted data
type UnsealedMiniSecret struct {
	id    *uuid.UUID // Optional, set when received from storage
	name  string
	key   []byte
	value []byte
}

// CreateMiniSecret creates a new unsealed secret
func CreateMiniSecret(name string, key, value []byte) (*UnsealedMiniSecret, error) {
	if len(key) == 0 || len(value) == 0 {
		return nil, ErrEmptyData
	}

	return &UnsealedMiniSecret{
		name:  name,
		key:   key,
		value: value,
	}, nil
}

// InitialSeal creates an InitialMiniSecret by encrypting the data
func (s *UnsealedMiniSecret) InitialSeal(enc crypto.Encrypter) (*InitialMiniSecret, error) {
	if enc == nil {
		return nil, ErrNilEncryption
	}

	// Encrypt key and value
	encKey, err := enc.Encrypt(s.key)
	if err != nil {
		return nil, err
	}

	encValue, err := enc.Encrypt(s.value)
	if err != nil {
		return nil, err
	}

	return &InitialMiniSecret{
		Name:        s.name,
		Key:         encKey,
		Value:       encValue,
		EncrypterID: enc.ID(), // Track which DEK encrypted this secret
	}, nil
}

// GetSecret returns the decrypted key and value
func (s *UnsealedMiniSecret) GetSecret() ([]byte, []byte) {
	return s.key, s.value
}

// String implements the Stringer interface for human-readable output
func (s *UnsealedMiniSecret) String() string {
	return fmt.Sprintf("%s: %s", string(s.key), string(s.value))
}

// Unseal decrypts the secret data using the provided decrypter
func (s *MiniSecret) Unseal(dec crypto.Decrypter) (*UnsealedMiniSecret, error) {
	if dec == nil {
		return nil, ErrNilDecryption
	}

	// Verify decrypter ID matches
	if s.EncrypterID != dec.ID() {
		return nil, ErrDecrypterIDMismatch
	}

	// Decrypt key
	key, err := dec.Decrypt(s.Key)
	if err != nil {
		return nil, err
	}

	// Decrypt value
	value, err := dec.Decrypt(s.Value)
	if err != nil {
		return nil, err
	}

	// Parse UUID from string
	id, err := uuid.Parse(s.ID)
	if err != nil {
		return nil, err
	}

	return &UnsealedMiniSecret{
		id:    &id,
		name:  s.Name,
		key:   key,
		value: value,
	}, nil
}
