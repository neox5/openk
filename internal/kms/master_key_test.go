package kms_test

import (
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/kms"
	"github.com/stretchr/testify/assert"
)

var (
	testPassword = []byte("test-password")
	testUsername = []byte("test-user")
)

func TestMasterKey_Derive(t *testing.T) {
	t.Run("successful derivation", func(t *testing.T) {
		mk := kms.NewMasterKey()
		assert.NoError(t, mk.Derive(testPassword, testUsername))
	})

	t.Run("empty password", func(t *testing.T) {
		mk := kms.NewMasterKey()
		assert.ErrorIs(t, mk.Derive([]byte{}, testUsername), kms.ErrInvalidPassword)
	})

	t.Run("empty username", func(t *testing.T) {
		mk := kms.NewMasterKey()
		assert.ErrorIs(t, mk.Derive(testPassword, []byte{}), kms.ErrInvalidUsername)
	})

	t.Run("key already set", func(t *testing.T) {
		mk := kms.NewMasterKey()
		assert.NoError(t, mk.Derive(testPassword, testUsername))
		assert.ErrorIs(t, mk.Derive(testPassword, testUsername), kms.ErrKeyAlreadySet)
	})
}

func TestMasterKey_Encrypt(t *testing.T) {
	plaintext := []byte("secret data")

	t.Run("encrypt and decrypt", func(t *testing.T) {
		mk := kms.NewMasterKey()
		assert.NoError(t, mk.Derive(testPassword, testUsername))

		// encrypt plaintext
		encrypted, err := mk.Encrypt(plaintext)
		assert.NoError(t, err)
		assert.NotNil(t, encrypted)

		// Decrypt ciphertext
		decrypted, err := mk.Decrypt(encrypted)
		assert.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("encrypt without master key", func(t *testing.T) {
		mk := kms.NewMasterKey()
		_, err := mk.Encrypt(plaintext)
		assert.ErrorIs(t, err, kms.ErrKeyNotDerived)
	})

	t.Run("decrypt without master key", func(t *testing.T) {
		mk := kms.NewMasterKey()
		_, err := mk.Decrypt(&crypto.Ciphertext{})
		assert.ErrorIs(t, err, kms.ErrKeyNotDerived)
	})
}

func TestMasterKey_Clear(t *testing.T) {
	mk := kms.NewMasterKey()
	assert.NoError(t, mk.Derive(testPassword, testUsername))
	assert.True(t, mk.HasKey())
	mk.Clear()
	assert.False(t, mk.HasKey())
}

func TestMasterKey_GetAuthKey(t *testing.T) {
	t.Run("get authKey", func(t *testing.T) {
		mk := kms.NewMasterKey()
		assert.NoError(t, mk.Derive(testPassword, testUsername))
		authKey, err := mk.GetAuthKey()
		assert.NoError(t, err)
		assert.NotNil(t, authKey)
	})

	t.Run("get authKey without derive", func(t *testing.T) {
		mk := kms.NewMasterKey()
		_, err := mk.GetAuthKey()
		assert.ErrorIs(t, err, kms.ErrKeyNotDerived)
	})
}
