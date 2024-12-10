package kms_test

import (
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyPair_Generate(t *testing.T) {
	unsealed, err := kms.GenerateKeyPair()
	require.NoError(t, err)
	require.NotNil(t, unsealed)
	defer unsealed.Clear()

	// Test basic encryption/decryption
	plaintext := []byte("test message")
	ct, err := unsealed.Encrypt(plaintext)
	require.NoError(t, err)
	require.NotNil(t, ct)

	decrypted, err := unsealed.Decrypt(ct)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestKeyPair_ID(t *testing.T) {
	t.Run("returns empty for new keypair", func(t *testing.T) {
		unsealed, err := kms.GenerateKeyPair()
		require.NoError(t, err)
		defer unsealed.Clear()

		assert.Empty(t, unsealed.ID())
	})

	t.Run("returns valid ID after unsealing", func(t *testing.T) {
		// Setup: create and seal a keypair
		masterKey := setupMasterKey(t)
		unsealed, err := kms.GenerateKeyPair()
		require.NoError(t, err)
		initial, err := unsealed.InitialSeal(masterKey)
		require.NoError(t, err)

		id := "123e4567-e89b-12d3-a456-426614174000" // Valid UUID
		stored := &kms.KeyPair{
			ID:          id,
			Algorithm:   initial.Algorithm,
			PublicKey:   initial.PublicKey,
			PrivateKey:  initial.PrivateKey,
			Created:     initial.Created,
			State:       initial.State,
			EncrypterID: initial.EncrypterID,
		}

		unsealed, err = stored.Unseal(masterKey)
		require.NoError(t, err)
		defer unsealed.Clear()

		assert.Equal(t, "keypair-"+id, unsealed.ID())
	})
}

func TestKeyPair_InitialSeal(t *testing.T) {
	masterKey := setupMasterKey(t)

	t.Run("success cases", func(t *testing.T) {
		t.Run("creates valid sealed keypair", func(t *testing.T) {
			unsealed, err := kms.GenerateKeyPair()
			require.NoError(t, err)

			initial, err := unsealed.InitialSeal(masterKey)
			require.NoError(t, err)
			require.NotNil(t, initial)

			assert.Equal(t, crypto.AlgorithmRSAOAEPSHA256, initial.Algorithm)
			assert.NotEmpty(t, initial.PublicKey)
			assert.NotNil(t, initial.PrivateKey)
			assert.Equal(t, crypto.KeyStateActive, initial.State)
			assert.False(t, initial.Created.IsZero())
			assert.Equal(t, masterKey.ID(), initial.EncrypterID)

			pubKey, err := crypto.ImportRSAPublicKey(initial.PublicKey)
			assert.NoError(t, err)
			assert.NotNil(t, pubKey)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("with nil encrypter", func(t *testing.T) {
			unsealed, err := kms.GenerateKeyPair()
			require.NoError(t, err)

			initial, err := unsealed.InitialSeal(nil)
			assert.ErrorIs(t, err, kms.ErrNilEncrypter)
			assert.Nil(t, initial)
		})

		t.Run("with destroyed state", func(t *testing.T) {
			unsealed, err := kms.GenerateKeyPair()
			require.NoError(t, err)

			unsealed.Clear() // Set to destroyed state

			initial, err := unsealed.InitialSeal(masterKey)
			assert.ErrorIs(t, err, kms.ErrKeyRevoked)
			assert.Nil(t, initial)
		})
	})
}

func TestKeyPair_Unseal(t *testing.T) {
	masterKey := setupMasterKey(t)

	// Setup: Generate and seal a key pair
	unsealed, err := kms.GenerateKeyPair()
	require.NoError(t, err)
	initial, err := unsealed.InitialSeal(masterKey)
	require.NoError(t, err)

	id := "123e4567-e89b-12d3-a456-426614174000" // Valid UUID
	stored := &kms.KeyPair{
		ID:          id,
		Algorithm:   initial.Algorithm,
		PublicKey:   initial.PublicKey,
		PrivateKey:  initial.PrivateKey,
		Created:     initial.Created,
		State:       initial.State,
		EncrypterID: initial.EncrypterID,
	}

	t.Run("success cases", func(t *testing.T) {
		t.Run("unseals valid key pair", func(t *testing.T) {
			unsealed, err := stored.Unseal(masterKey)
			require.NoError(t, err)
			require.NotNil(t, unsealed)
			defer unsealed.Clear()

			plaintext := []byte("test message")
			ct, err := unsealed.Encrypt(plaintext)
			require.NoError(t, err)
			require.NotNil(t, ct)

			decrypted, err := unsealed.Decrypt(ct)
			require.NoError(t, err)
			assert.Equal(t, plaintext, decrypted)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name    string
			setup   func() (*kms.KeyPair, crypto.Decrypter)
			wantErr error
		}{
			{
				name: "with nil decrypter",
				setup: func() (*kms.KeyPair, crypto.Decrypter) {
					return stored, nil
				},
				wantErr: kms.ErrNilDecrypter,
			},
			{
				name: "with destroyed state",
				setup: func() (*kms.KeyPair, crypto.Decrypter) {
					destroyed := &kms.KeyPair{
						ID:          id,
						Algorithm:   stored.Algorithm,
						PublicKey:   stored.PublicKey,
						PrivateKey:  stored.PrivateKey,
						Created:     stored.Created,
						State:       crypto.KeyStateDestroyed,
						EncrypterID: stored.EncrypterID,
					}
					return destroyed, masterKey
				},
				wantErr: kms.ErrKeyRevoked,
			},
			{
				name: "with mismatched decrypter ID",
				setup: func() (*kms.KeyPair, crypto.Decrypter) {
					mismatched := &kms.KeyPair{
						ID:          id,
						Algorithm:   stored.Algorithm,
						PublicKey:   stored.PublicKey,
						PrivateKey:  stored.PrivateKey,
						Created:     stored.Created,
						State:       stored.State,
						EncrypterID: "different-id",
					}
					return mismatched, masterKey
				},
				wantErr: kms.ErrDecrypterIDMismatch,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				kp, dec := tt.setup()
				unsealed, err := kp.Unseal(dec)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, unsealed)
			})
		}
	})
}

func TestKeyPair_Operations(t *testing.T) {
	unsealed, err := kms.GenerateKeyPair()
	require.NoError(t, err)
	defer unsealed.Clear()

	t.Run("encryption operations", func(t *testing.T) {
		t.Run("handles empty data", func(t *testing.T) {
			ct, err := unsealed.Encrypt([]byte{})
			assert.NoError(t, err)
			require.NotNil(t, ct)

			decrypted, err := unsealed.Decrypt(ct)
			assert.NoError(t, err)
			assert.Empty(t, decrypted)
		})

		t.Run("handles nil data", func(t *testing.T) {
			ct, err := unsealed.Encrypt(nil)
			assert.NoError(t, err)
			require.NotNil(t, ct)

			decrypted, err := unsealed.Decrypt(ct)
			assert.NoError(t, err)
			assert.Empty(t, decrypted)
		})

		t.Run("rejects operations after destroy", func(t *testing.T) {
			unsealed.Clear()
			ct, err := unsealed.Encrypt([]byte("test"))
			assert.ErrorIs(t, err, kms.ErrKeyRevoked)
			assert.Nil(t, ct)
		})
	})

	t.Run("decryption operations", func(t *testing.T) {
		t.Run("fails with invalid ciphertext", func(t *testing.T) {
			ct, err := crypto.NewCiphertext(
				make([]byte, crypto.NonceSize),
				[]byte("invalid data"),
				make([]byte, crypto.TagSize),
			)
			require.NoError(t, err)

			decrypted, err := unsealed.Decrypt(ct)
			assert.Error(t, err)
			assert.Nil(t, decrypted)
		})
	})
}
