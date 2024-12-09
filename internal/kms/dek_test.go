package kms_test

import (
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestType_GenerateDEK(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("generates valid DEK", func(t *testing.T) {
			dek, err := kms.GenerateDEK()
			require.NoError(t, err)
			require.NotNil(t, dek)
			defer dek.Clear()

			// Test the generated DEK can encrypt/decrypt
			plaintext := []byte("test message")
			ct, err := dek.Encrypt(plaintext)
			require.NoError(t, err)
			require.NotNil(t, ct)

			decrypted, err := dek.Decrypt(ct)
			require.NoError(t, err)
			assert.Equal(t, plaintext, decrypted)
		})
	})
}

func TestType_DEKSeal(t *testing.T) {
	masterKey := setupMasterKey(t)

	t.Run("success cases", func(t *testing.T) {
		t.Run("creates initial DEK with envelope", func(t *testing.T) {
			dek, err := kms.GenerateDEK()
			require.NoError(t, err)
			defer dek.Clear()

			initial, err := dek.Seal(masterKey)
			require.NoError(t, err)
			require.NotNil(t, initial)

			// Verify DEK properties
			assert.Equal(t, crypto.AlgorithmAESGCM256, initial.Algorithm)
			assert.Equal(t, crypto.KeyStateActive, initial.State)
			assert.False(t, initial.Created.IsZero())

			// Verify envelope
			require.Len(t, initial.Envelopes, 1)
			env := initial.Envelopes[0]
			assert.Equal(t, crypto.AlgorithmAESGCM256, env.Algorithm)
			assert.Equal(t, masterKey.ID(), env.EncrypterID)
			assert.Equal(t, crypto.KeyStateActive, env.State)
			assert.NotNil(t, env.Key)
			assert.False(t, env.Created.IsZero())
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil encrypter", func(t *testing.T) {
			dek, err := kms.GenerateDEK()
			require.NoError(t, err)
			defer dek.Clear()

			initial, err := dek.Seal(nil)
			assert.ErrorIs(t, err, kms.ErrNilEncrypter)
			assert.Nil(t, initial)
		})
	})
}

func TestType_CreateEnvelope(t *testing.T) {
	masterKey := setupMasterKey(t)

	t.Run("success cases", func(t *testing.T) {
		t.Run("creates additional envelope", func(t *testing.T) {
			dek, err := kms.GenerateDEK()
			require.NoError(t, err)
			defer dek.Clear()

			env, err := dek.CreateEnvelope(masterKey)
			require.NoError(t, err)
			require.NotNil(t, env)

			assert.Equal(t, crypto.AlgorithmAESGCM256, env.Algorithm)
			assert.Equal(t, masterKey.ID(), env.EncrypterID)
			assert.Equal(t, crypto.KeyStateActive, env.State)
			assert.NotNil(t, env.Key)
			assert.False(t, env.Created.IsZero())
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil encrypter", func(t *testing.T) {
			dek, err := kms.GenerateDEK()
			require.NoError(t, err)
			defer dek.Clear()

			env, err := dek.CreateEnvelope(nil)
			assert.ErrorIs(t, err, kms.ErrNilEncrypter)
			assert.Nil(t, env)
		})
	})
}

func TestType_DEKUnseal(t *testing.T) {
	masterKey := setupMasterKey(t)

	// Helper to create a stored DEK
	setupStoredDEK := func(t *testing.T) (*kms.DEK, *kms.Envelope) {
		t.Helper()
		dek, err := kms.GenerateDEK()
		require.NoError(t, err)
		defer dek.Clear()

		initial, err := dek.Seal(masterKey)
		require.NoError(t, err)

		dekID := "123e4567-e89b-12d3-a456-426614174000"
		envID := "123e4567-e89b-12d3-a456-426614174001"

		storedDEK := &kms.DEK{
			ID:        dekID,
			Algorithm: initial.Algorithm,
			Created:   initial.Created,
			State:     initial.State,
			Envelopes: make(map[string]*kms.Envelope),
		}

		env := &kms.Envelope{
			ID:          envID,
			DEKID:       storedDEK.ID,
			Algorithm:   initial.Envelopes[0].Algorithm,
			Key:         initial.Envelopes[0].Key,
			Created:     initial.Envelopes[0].Created,
			State:       initial.Envelopes[0].State,
			EncrypterID: masterKey.ID(),
		}
		storedDEK.Envelopes[masterKey.ID()] = env // Use masterKey.ID() as map key

		return storedDEK, env
	}

	t.Run("success cases", func(t *testing.T) {
		t.Run("unseals with valid decrypter", func(t *testing.T) {
			storedDEK, _ := setupStoredDEK(t)

			unsealed, err := storedDEK.Unseal(masterKey)
			require.NoError(t, err)
			require.NotNil(t, unsealed)
			defer unsealed.Clear()

			// Test the unsealed DEK can encrypt/decrypt
			plaintext := []byte("test message")
			ct, err := unsealed.Encrypt(plaintext)
			require.NoError(t, err)

			decrypted, err := unsealed.Decrypt(ct)
			require.NoError(t, err)
			assert.Equal(t, plaintext, decrypted)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil decrypter", func(t *testing.T) {
			storedDEK, _ := setupStoredDEK(t)

			unsealed, err := storedDEK.Unseal(nil)
			assert.ErrorIs(t, err, kms.ErrNilDecrypter)
			assert.Nil(t, unsealed)
		})

		t.Run("rejects when no valid envelope exists", func(t *testing.T) {
			storedDEK, _ := setupStoredDEK(t)
			// Create a master key with different username to get different ID
			otherKey := kms.NewMasterKey()
			err := otherKey.Derive([]byte("other-password"), []byte("other-user"))
			require.NoError(t, err)

			unsealed, err := storedDEK.Unseal(otherKey)
			assert.ErrorIs(t, err, kms.ErrNoValidEnvelope)
			assert.Nil(t, unsealed)
		})

		t.Run("rejects destroyed DEK", func(t *testing.T) {
			storedDEK, _ := setupStoredDEK(t)
			storedDEK.State = crypto.KeyStateDestroyed

			unsealed, err := storedDEK.Unseal(masterKey)
			assert.ErrorIs(t, err, kms.ErrKeyRevoked)
			assert.Nil(t, unsealed)
		})

		t.Run("rejects destroyed envelope", func(t *testing.T) {
			storedDEK, env := setupStoredDEK(t)
			env.State = crypto.KeyStateDestroyed

			unsealed, err := storedDEK.Unseal(masterKey)
			assert.ErrorIs(t, err, kms.ErrKeyRevoked)
			assert.Nil(t, unsealed)
		})
	})
}

func TestType_DEKOperations(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("handles empty data", func(t *testing.T) {
			dek, err := kms.GenerateDEK()
			require.NoError(t, err)
			defer dek.Clear()

			ct, err := dek.Encrypt([]byte{})
			require.NoError(t, err)
			require.NotNil(t, ct)

			decrypted, err := dek.Decrypt(ct)
			require.NoError(t, err)
			assert.Empty(t, decrypted)
		})

		t.Run("cleans up on Clear", func(t *testing.T) {
			dek, err := kms.GenerateDEK()
			require.NoError(t, err)

			// First operation works
			ct, err := dek.Encrypt([]byte("test"))
			require.NoError(t, err)
			require.NotNil(t, ct)

			// Clear the key
			dek.Clear()

			// Operation after clear should fail
			ct2, err := dek.Encrypt([]byte("test"))
			assert.ErrorIs(t, err, kms.ErrInvalidDEK)
			assert.Nil(t, ct2)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("fails with invalid ciphertext", func(t *testing.T) {
			dek, err := kms.GenerateDEK()
			require.NoError(t, err)
			defer dek.Clear()

			ct, err := crypto.NewCiphertext(
				make([]byte, crypto.NonceSize),
				[]byte("invalid data"),
				make([]byte, crypto.TagSize),
			)
			require.NoError(t, err)

			decrypted, err := dek.Decrypt(ct)
			assert.Error(t, err)
			assert.Nil(t, decrypted)
		})
	})
}
