package kms

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neox5/openk/internal/crypto"
)

var ErrNoValidEnvelope = errors.New("no valid envelope found for decrypter")

type InitialDEK struct {
	Algorithm crypto.Algorithm
	Created   time.Time
	State     crypto.KeyState
	Envelopes []*InitialEnvelope
}

type InitialEnvelope struct {
	Algorithm   crypto.Algorithm
	Key         *crypto.Ciphertext
	Created     time.Time
	State       crypto.KeyState
	EncrypterID string
}

type DEK struct {
	ID        string
	Algorithm crypto.Algorithm
	Created   time.Time
	State     crypto.KeyState
	Envelopes map[string]*Envelope
}

type Envelope struct {
	ID          string
	DEKID       string
	Algorithm   crypto.Algorithm
	Key         *crypto.Ciphertext
	Created     time.Time
	State       crypto.KeyState
	EncrypterID string
}

type UnsealedDEK struct {
	id    *uuid.UUID
	key   []byte
	state crypto.KeyState
}

func GenerateDEK() (*UnsealedDEK, error) {
	key, err := crypto.AESGenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %w", err)
	}

	return &UnsealedDEK{
		key:   key,
		state: crypto.KeyStateActive,
	}, nil
}

func (dek *UnsealedDEK) Seal(enc crypto.Encrypter) (*InitialDEK, error) {
	if enc == nil {
		return nil, ErrNilEncrypter
	}
	if len(dek.key) != crypto.AESKeySize {
		return nil, ErrInvalidDEK
	}
	if dek.state == crypto.KeyStateDestroyed {
		return nil, ErrKeyRevoked
	}

	// Create initial envelope
	env, err := dek.CreateEnvelope(enc)
	if err != nil {
		return nil, err
	}

	return &InitialDEK{
		Algorithm: crypto.AlgorithmAESGCM256,
		Created:   time.Now(),
		State:     dek.state,
		Envelopes: []*InitialEnvelope{env},
	}, nil
}

func (dek *UnsealedDEK) CreateEnvelope(enc crypto.Encrypter) (*InitialEnvelope, error) {
	if enc == nil {
		return nil, ErrNilEncrypter
	}
	if len(dek.key) != crypto.AESKeySize {
		return nil, ErrInvalidDEK
	}
	if dek.state == crypto.KeyStateDestroyed {
		return nil, ErrKeyRevoked
	}

	ciphertext, err := enc.Encrypt(dek.key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt DEK: %w", err)
	}

	return &InitialEnvelope{
		Algorithm:   crypto.AlgorithmAESGCM256,
		Key:         ciphertext,
		Created:     time.Now(),
		State:       crypto.KeyStateActive,
		EncrypterID: enc.ID(),
	}, nil
}

func (dek *UnsealedDEK) Encrypt(data []byte) (*crypto.Ciphertext, error) {
	if len(dek.key) != crypto.AESKeySize {
		return nil, ErrInvalidDEK
	}
	if dek.state == crypto.KeyStateDestroyed {
		return nil, ErrKeyRevoked
	}

	return crypto.AESEncrypt(dek.key, data)
}

func (dek *UnsealedDEK) Decrypt(ct *crypto.Ciphertext) ([]byte, error) {
	if len(dek.key) != crypto.AESKeySize {
		return nil, ErrInvalidDEK
	}
	if dek.state == crypto.KeyStateDestroyed {
		return nil, ErrKeyRevoked
	}

	return crypto.AESDecrypt(dek.key, ct)
}

func (dek *UnsealedDEK) Clear() {
	if dek.key != nil {
		crypto.SecureWipe(dek.key)
		dek.key = nil
	}
}

func (dek *UnsealedDEK) ID() string {
	if dek.id == nil {
		return ""
	}
	return fmt.Sprintf("dek-%s", dek.id.String())
}

func (dek *DEK) Unseal(dec crypto.Decrypter) (*UnsealedDEK, error) {
	if dec == nil {
		return nil, ErrNilDecrypter
	}
	if dek.State == crypto.KeyStateDestroyed {
		return nil, ErrKeyRevoked
	}

	envelope, exists := dek.Envelopes[dec.ID()]
	if !exists {
		return nil, ErrNoValidEnvelope
	}

	if envelope.State == crypto.KeyStateDestroyed {
		return nil, ErrKeyRevoked
	}

	// Parse ID from string
	id, err := uuid.Parse(dek.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid DEK ID: %w", err)
	}

	// Decrypt the key material
	key, err := dec.Decrypt(envelope.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt DEK: %w", err)
	}

	// Validate key length
	if len(key) != crypto.AESKeySize {
		crypto.SecureWipe(key)
		return nil, ErrInvalidDEK
	}

	return &UnsealedDEK{
		id:    &id,
		key:   key,
		state: dek.State,
	}, nil
}
