package crypto_test

import (
	"bytes"
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSalt(t *testing.T) {
	// Test correct length
	salt1, err := crypto.GenerateSalt()
	assert.NoError(t, err)
	assert.Equal(t, crypto.DefaultSaltSize, len(salt1))

	// Test uniqueness
	salt2, err := crypto.GenerateSalt()
	assert.NoError(t, err)
	assert.False(t, bytes.Equal(salt1, salt2))
}

func TestDeriveKey(t *testing.T) {
	salt := []byte("test-username") // Using username-like salt
	password := []byte("test-password")

	t.Run("valid input", func(t *testing.T) {
		key, err := crypto.DeriveKey(password, salt, 1000, 32)
		assert.NoError(t, err)
		assert.Equal(t, 32, len(key))
	})

	t.Run("empty password", func(t *testing.T) {
		key, err := crypto.DeriveKey([]byte{}, salt, 1000, 32)
		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Equal(t, crypto.ErrEmptyPassword, err)
	})

	t.Run("empty salt", func(t *testing.T) {
		key, err := crypto.DeriveKey(password, []byte{}, 1000, 32)
		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Equal(t, crypto.ErrInvalidSalt, err)
	})

	t.Run("different passwords produce different keys", func(t *testing.T) {
		key1, err := crypto.DeriveKey([]byte("password1"), salt, 1000, 32)
		assert.NoError(t, err)
		key2, err := crypto.DeriveKey([]byte("password2"), salt, 1000, 32)
		assert.NoError(t, err)
		assert.False(t, bytes.Equal(key1, key2))
	})

	t.Run("same input reproduces same key", func(t *testing.T) {
		key1, err := crypto.DeriveKey(password, salt, 1000, 32)
		assert.NoError(t, err)
		key2, err := crypto.DeriveKey(password, salt, 1000, 32)
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(key1, key2))
	})

	t.Run("different iterations produce different keys", func(t *testing.T) {
		key1, err := crypto.DeriveKey(password, salt, 1000, 32)
		assert.NoError(t, err)
		key2, err := crypto.DeriveKey(password, salt, 2000, 32)
		assert.NoError(t, err)
		assert.False(t, bytes.Equal(key1, key2))
	})

	t.Run("variable length salts", func(t *testing.T) {
		// Test with different username-like salts
		salts := [][]byte{
			[]byte("short"),
			[]byte("medium-length-username"),
			[]byte("very-long-username-that-exceeds-sixteen-bytes"),
		}

		for _, salt := range salts {
			key, err := crypto.DeriveKey(password, salt, 1000, 32)
			assert.NoError(t, err)
			assert.Equal(t, 32, len(key))
		}
	})
}

func TestDeriveMasterKey(t *testing.T) {
	salt := []byte("test-username")
	password := []byte("test-password")

	t.Run("produces correct key size", func(t *testing.T) {
		key, err := crypto.DeriveMasterKey(password, salt)
		assert.NoError(t, err)
		assert.Equal(t, crypto.MasterKeySize, len(key))
	})

	t.Run("validates input", func(t *testing.T) {
		key, err := crypto.DeriveMasterKey([]byte{}, salt)
		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Equal(t, crypto.ErrEmptyPassword, err)

		key, err = crypto.DeriveMasterKey(password, []byte{})
		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Equal(t, crypto.ErrInvalidSalt, err)
	})

	t.Run("reproduces same key", func(t *testing.T) {
		key1, err := crypto.DeriveMasterKey(password, salt)
		assert.NoError(t, err)
		key2, err := crypto.DeriveMasterKey(password, salt)
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(key1, key2))
	})
}

