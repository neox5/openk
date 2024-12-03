package kms

import (
	"time"

	"github.com/neox5/openk/internal/crypto"
)

// InitialDEK represents a newly generated DEK before server storage
type InitialDEK struct {
	Algorithm crypto.Algorithm
	Key       *crypto.Ciphertext // Encrypted with protection key
	Created   time.Time
	State     crypto.KeyState
}

// DEK represents a Data Encryption Key as stored in the backend
type DEK struct {
	ID        string
	Algorithm crypto.Algorithm
	Key       *crypto.Ciphertext
	Created   time.Time
	State     crypto.KeyState
}

// UnsealedDEK represents an active DEK with access to key material
type UnsealedDEK struct {
	key []byte // Raw key material, protected in memory
}

// GenerateDEK creates a new random AES-256 key
func GenerateDEK() (*UnsealedDEK, error) {
	// Generate random AES-256 key
	key, err := crypto.AESGenerateKey()
	if err != nil {
		return nil, err
	}

	return &UnsealedDEK{
		key: key,
	}, nil
}

// InitialSeal creates an InitialDEK by encrypting the key using the provided encrypter
func (dk *UnsealedDEK) InitialSeal(enc crypto.Encrypter) (*InitialDEK, error) {
	if enc == nil {
		return nil, ErrNilEncrypter
	}

	// Encrypt the key with provided encrypter
	encKey, err := enc.Encrypt(dk.key)
	if err != nil {
		return nil, err
	}

	return &InitialDEK{
		Algorithm: crypto.AlgorithmAESGCM256,
		Key:       encKey,
		Created:   time.Now(),
		State:     crypto.KeyStateActive,
	}, nil
}

// Unseal decrypts the key using the provided decrypter and returns an UnsealedDEK
func (dk *DEK) Unseal(dec crypto.Decrypter) (*UnsealedDEK, error) {
	if dec == nil {
		return nil, ErrNilDecrypter
	}

	if dk.State == crypto.KeyStateDestroyed {
		return nil, ErrKeyRevoked
	}

	// Decrypt the key
	key, err := dec.Decrypt(dk.Key)
	if err != nil {
		return nil, err
	}

	// Validate key length
	if len(key) != crypto.AESKeySize {
		return nil, ErrInvalidDEK
	}

	return &UnsealedDEK{
		key: key,
	}, nil
}

// Encrypt performs AES-256-GCM encryption of data using the key
func (dk *UnsealedDEK) Encrypt(data []byte) (*crypto.Ciphertext, error) {
	if len(dk.key) != crypto.AESKeySize {
		return nil, ErrInvalidDEK
	}

	return crypto.AESEncrypt(dk.key, data)
}

// Decrypt performs AES-256-GCM decryption of data using the key
func (dk *UnsealedDEK) Decrypt(ct *crypto.Ciphertext) ([]byte, error) {
	if len(dk.key) != crypto.AESKeySize {
		return nil, ErrInvalidDEK
	}

	return crypto.AESDecrypt(dk.key, ct)
}

// Clear wipes the key material from memory
func (dk *UnsealedDEK) Clear() {
	if dk.key != nil {
		crypto.SecureWipe(dk.key)
		dk.key = nil
	}
}
