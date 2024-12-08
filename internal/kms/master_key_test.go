package kms_test

import (
	"bytes"
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Common test values
var (
    validPassword = []byte("test-password")
    validUsername = []byte("test-user")
)
// Helper function for test setup
func setupMasterKey(t *testing.T) *kms.MasterKey {
    t.Helper()
    mk := kms.NewMasterKey()
    err := mk.Derive(validPassword, validUsername)
    require.NoError(t, err)
    return mk
}

func TestMasterKey_New(t *testing.T) {
	mk := kms.NewMasterKey()
	assert.NotNil(t, mk)
	assert.False(t, mk.HasKey())
}

func TestMasterKey_Derive(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		tests := []struct {
			name     string
			password []byte
			username []byte
		}{
			{
				name:     "derives with valid credentials",
				password: validPassword,
				username: validUsername,
			},
			{
				name:     "derives with minimum length values",
				password: []byte("x"),
				username: []byte("y"),
			},
			{
				name:     "derives with long values",
				password: bytes.Repeat([]byte("long"), 100),
				username: bytes.Repeat([]byte("user"), 100),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mk := kms.NewMasterKey()
				err := mk.Derive(tt.password, tt.username)
				assert.NoError(t, err)
				assert.True(t, mk.HasKey())

				authKey, err := mk.GetAuthKey()
				assert.NoError(t, err)
				assert.NotNil(t, authKey)
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name     string
			password []byte
			username []byte
			wantErr  error
		}{
			{
				name:     "fails with empty password",
				password: []byte{},
				username: validUsername,
				wantErr:  kms.ErrInvalidPassword,
			},
			{
				name:     "fails with nil password",
				password: nil,
				username: validUsername,
				wantErr:  kms.ErrInvalidPassword,
			},
			{
				name:     "fails with empty username",
				password: validPassword,
				username: []byte{},
				wantErr:  kms.ErrInvalidUsername,
			},
			{
				name:     "fails with nil username",
				password: validPassword,
				username: nil,
				wantErr:  kms.ErrInvalidUsername,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mk := kms.NewMasterKey()
				err := mk.Derive(tt.password, tt.username)
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.False(t, mk.HasKey())

				authKey, err := mk.GetAuthKey()
				assert.Error(t, err)
				assert.Nil(t, authKey)
			})
		}
	})

	t.Run("prevents repeated derivation", func(t *testing.T) {
		mk := kms.NewMasterKey()
		err := mk.Derive(validPassword, validUsername)
		require.NoError(t, err)

		err = mk.Derive(validPassword, validUsername)
		assert.ErrorIs(t, err, kms.ErrKeyAlreadySet)
	})
}

func TestMasterKey_Clear(t *testing.T) {
	t.Run("clears derived key", func(t *testing.T) {
		mk := setupMasterKey(t)
		assert.True(t, mk.HasKey())

		mk.Clear()
		assert.False(t, mk.HasKey())

		authKey, err := mk.GetAuthKey()
		assert.Error(t, err)
		assert.Nil(t, authKey)
	})

	t.Run("safe to clear multiple times", func(t *testing.T) {
		mk := setupMasterKey(t)
		mk.Clear()
		mk.Clear() // Should not panic
		assert.False(t, mk.HasKey())
	})
}

func TestMasterKey_Encrypt(t *testing.T) {
	plaintext := []byte("secret data")

	t.Run("success cases", func(t *testing.T) {
		tests := []struct {
			name      string
			plaintext []byte
		}{
			{
				name:      "encrypts regular message",
				plaintext: plaintext,
			},
			{
				name:      "encrypts empty message",
				plaintext: []byte{},
			},
			{
				name:      "encrypts large message",
				plaintext: bytes.Repeat([]byte("large"), 1000),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mk := setupMasterKey(t)
				defer mk.Clear()

				ct, err := mk.Encrypt(tt.plaintext)
				require.NoError(t, err)
				require.NotNil(t, ct)

				assert.Equal(t, crypto.NonceSize, len(ct.Nonce))
				assert.Equal(t, crypto.TagSize, len(ct.Tag))

				decrypted, err := mk.Decrypt(ct)
				require.NoError(t, err)
				assert.Equal(t, tt.plaintext, decrypted)
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name    string
			setup   func() *kms.MasterKey
			wantErr error
		}{
			{
				name: "fails without derived key",
				setup: func() *kms.MasterKey {
					return kms.NewMasterKey()
				},
				wantErr: kms.ErrKeyNotDerived,
			},
			{
				name: "fails after clear",
				setup: func() *kms.MasterKey {
					mk := setupMasterKey(t)
					mk.Clear()
					return mk
				},
				wantErr: kms.ErrKeyNotDerived,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mk := tt.setup()
				ct, err := mk.Encrypt(plaintext)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, ct)
			})
		}
	})
}

func TestMasterKey_Decrypt(t *testing.T) {
	t.Run("error cases", func(t *testing.T) {
		validMK := setupMasterKey(t)
		validCT, err := validMK.Encrypt([]byte("test"))
		require.NoError(t, err)

		tests := []struct {
			name    string
			setup   func() (*kms.MasterKey, *crypto.Ciphertext)
			wantErr error
		}{
			{
				name: "fails without derived key",
				setup: func() (*kms.MasterKey, *crypto.Ciphertext) {
					return kms.NewMasterKey(), validCT
				},
				wantErr: kms.ErrKeyNotDerived,
			},
			{
				name: "fails with nil ciphertext",
				setup: func() (*kms.MasterKey, *crypto.Ciphertext) {
					return setupMasterKey(t), nil
				},
				wantErr: crypto.ErrAESInvalidMessage,
			},
			{
				name: "fails after clear",
				setup: func() (*kms.MasterKey, *crypto.Ciphertext) {
					mk := setupMasterKey(t)
					mk.Clear()
					return mk, validCT
				},
				wantErr: kms.ErrKeyNotDerived,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mk, ct := tt.setup()
				decrypted, err := mk.Decrypt(ct)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, decrypted)
			})
		}
	})
}

func TestMasterKey_GetAuthKey(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		mk := setupMasterKey(t)

		authKey1, err := mk.GetAuthKey()
		require.NoError(t, err)
		require.NotNil(t, authKey1)

		authKey2, err := mk.GetAuthKey()
		require.NoError(t, err)
		assert.True(t, bytes.Equal(authKey1, authKey2))
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name    string
			setup   func() *kms.MasterKey
			wantErr error
		}{
			{
				name: "fails without derivation",
				setup: func() *kms.MasterKey {
					return kms.NewMasterKey()
				},
				wantErr: kms.ErrKeyNotDerived,
			},
			{
				name: "fails after clear",
				setup: func() *kms.MasterKey {
					mk := setupMasterKey(t)
					mk.Clear()
					return mk
				},
				wantErr: kms.ErrKeyNotDerived,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mk := tt.setup()
				authKey, err := mk.GetAuthKey()
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, authKey)
			})
		}
	})
}
