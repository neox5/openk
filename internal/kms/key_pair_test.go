package kms_test

import (
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMasterKey(t *testing.T) *kms.MasterKey {
	mk := kms.NewMasterKey()
	err := mk.Derive([]byte("test-password"), []byte("test-user"))
	require.NoError(t, err)
	return mk
}

func TestGenerateKeyPair(t *testing.T) {
	unsealed, err := kms.GenerateKeyPair()
	require.NoError(t, err)
	require.NotNil(t, unsealed)
	defer unsealed.Clear()

	// Test basic functionality
	plaintext := []byte("test message")
	ct, err := unsealed.Encrypt(plaintext)
	require.NoError(t, err)
	require.NotNil(t, ct)

	decrypted, err := unsealed.Decrypt(ct)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestInitialSeal(t *testing.T) {
	masterKey := setupMasterKey(t)

	t.Run("successful seal", func(t *testing.T) {
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

		// Verify public key format
		pubKey, err := crypto.ImportRSAPublicKey(initial.PublicKey)
		assert.NoError(t, err)
		assert.NotNil(t, pubKey)
	})

	t.Run("with nil encrypter", func(t *testing.T) {
		unsealed, err := kms.GenerateKeyPair()
		require.NoError(t, err)

		initial, err := unsealed.InitialSeal(nil)
		assert.ErrorIs(t, err, kms.ErrNilEncrypter)
		assert.Nil(t, initial)
	})
}

func TestUnseal(t *testing.T) {
	masterKey := setupMasterKey(t)

	// Setup: Generate and seal a key pair
	unsealed, err := kms.GenerateKeyPair()
	require.NoError(t, err)
	initial, err := unsealed.InitialSeal(masterKey)
	require.NoError(t, err)

	stored := &kms.KeyPair{
		ID:         "test-id",
		Algorithm:  initial.Algorithm,
		PublicKey:  initial.PublicKey,
		PrivateKey: initial.PrivateKey,
		Created:    initial.Created,
		State:      initial.State,
	}

	t.Run("successful unseal", func(t *testing.T) {
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

	t.Run("with nil decrypter", func(t *testing.T) {
		unsealed, err := stored.Unseal(nil)
		assert.ErrorIs(t, err, kms.ErrNilDecrypter)
		assert.Nil(t, unsealed)
	})

	t.Run("with destroyed state", func(t *testing.T) {
		destroyed := &kms.KeyPair{
			ID:         stored.ID,
			Algorithm:  stored.Algorithm,
			PublicKey:  stored.PublicKey,
			PrivateKey: stored.PrivateKey,
			Created:    stored.Created,
			State:      crypto.KeyStateDestroyed,
		}

		unsealed, err := destroyed.Unseal(masterKey)
		assert.ErrorIs(t, err, kms.ErrKeyRevoked)
		assert.Nil(t, unsealed)
	})
}

func TestUnsealedKeyPairOperations(t *testing.T) {
	unsealed, err := kms.GenerateKeyPair()
	require.NoError(t, err)
	defer unsealed.Clear()

	t.Run("encryption with nil/empty data", func(t *testing.T) {
		ct, err := unsealed.Encrypt(nil)
		assert.NoError(t, err)
		require.NotNil(t, ct)

		decrypted, err := unsealed.Decrypt(ct)
		assert.NoError(t, err)
		assert.Empty(t, decrypted)

		ct, err = unsealed.Encrypt([]byte{})
		assert.NoError(t, err)
		require.NotNil(t, ct)

		decrypted, err = unsealed.Decrypt(ct)
		assert.NoError(t, err)
		assert.Empty(t, decrypted)
	})

	t.Run("decryption with invalid ciphertext", func(t *testing.T) {
		invalidCT, err := crypto.NewCiphertext(
			make([]byte, crypto.NonceSize),
			[]byte("invalid data"),
			make([]byte, crypto.TagSize),
		)
		require.NoError(t, err)

		decrypted, err := unsealed.Decrypt(invalidCT)
		assert.Error(t, err)
		assert.Nil(t, decrypted)
	})
}

func TestUnsealedKeyPairClear(t *testing.T) {
	// we leave this test out as privateKey is only an internal 
	// property and cannot accessed from outside. Therefore it 
	// is not possible to test it after we cleared it.
}
