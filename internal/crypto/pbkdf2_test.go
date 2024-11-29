package crypto_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/neox5/openk/internal/crypto"
    "bytes"
)

func TestGenerateSalt(t *testing.T) {
    // Test correct length
    salt1, err := crypto.GenerateSalt()
    assert.NoError(t, err)
    assert.Equal(t, crypto.SaltSize, len(salt1))

    // Test uniqueness
    salt2, err := crypto.GenerateSalt()
    assert.NoError(t, err)
    assert.False(t, bytes.Equal(salt1, salt2))
}

func TestDeriveKey(t *testing.T) {
    salt := make([]byte, crypto.SaltSize)
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

    t.Run("invalid salt size", func(t *testing.T) {
        key, err := crypto.DeriveKey(password, []byte{1,2,3}, 1000, 32)
        assert.Error(t, err)
        assert.Nil(t, key)
        assert.Equal(t, crypto.ErrInvalidSaltSize, err)
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
}

func TestDeriveMasterKey(t *testing.T) {
    salt := make([]byte, crypto.SaltSize)
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

        key, err = crypto.DeriveMasterKey(password, []byte{1,2,3})
        assert.Error(t, err)
        assert.Nil(t, key)
        assert.Equal(t, crypto.ErrInvalidSaltSize, err)
    })

    t.Run("reproduces same key", func(t *testing.T) {
        key1, err := crypto.DeriveMasterKey(password, salt)
        assert.NoError(t, err)
        key2, err := crypto.DeriveMasterKey(password, salt)
        assert.NoError(t, err)
        assert.True(t, bytes.Equal(key1, key2))
    })
}
