package crypto_test

import (
	"bytes"
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPBKDF2_GenerateSalt(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		// Generate first salt
		salt1, err := crypto.GenerateSalt()
		require.NoError(t, err)
		assert.Equal(t, crypto.DefaultSaltSize, len(salt1))

		// Generate second salt to test uniqueness
		salt2, err := crypto.GenerateSalt()
		require.NoError(t, err)
		assert.Equal(t, crypto.DefaultSaltSize, len(salt2))

		// Verify salts are unique
		assert.False(t, bytes.Equal(salt1, salt2))
	})
}

func TestPBKDF2_DeriveKey(t *testing.T) {
	// Common test values
	validPassword := []byte("test-password")
	validSalt := []byte("test-username")
	validIterations := 1000
	validKeyLen := 32

	t.Run("success cases", func(t *testing.T) {
		tests := []struct {
			name       string
			password   []byte
			salt       []byte
			iterations int
			keyLen     int
		}{
			{
				name:       "derives key with valid input",
				password:   validPassword,
				salt:       validSalt,
				iterations: validIterations,
				keyLen:     validKeyLen,
			},
			{
				name:       "works with minimum length salt",
				password:   validPassword,
				salt:       []byte("x"),
				iterations: validIterations,
				keyLen:     validKeyLen,
			},
			{
				name:       "works with long username-like salt",
				password:   validPassword,
				salt:       []byte("very-long-username-that-exceeds-sixteen-bytes"),
				iterations: validIterations,
				keyLen:     validKeyLen,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				key, err := crypto.DeriveKey(tt.password, tt.salt, tt.iterations, tt.keyLen)
				require.NoError(t, err)
				assert.Equal(t, tt.keyLen, len(key))

				// Verify deterministic output
				key2, err := crypto.DeriveKey(tt.password, tt.salt, tt.iterations, tt.keyLen)
				require.NoError(t, err)
				assert.True(t, bytes.Equal(key, key2))
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name       string
			password   []byte
			salt       []byte
			iterations int
			keyLen     int
			wantErr    error
		}{
			{
				name:       "fails with empty password",
				password:   []byte{},
				salt:       validSalt,
				iterations: validIterations,
				keyLen:     validKeyLen,
				wantErr:    crypto.ErrEmptyPassword,
			},
			{
				name:       "fails with nil password",
				password:   nil,
				salt:       validSalt,
				iterations: validIterations,
				keyLen:     validKeyLen,
				wantErr:    crypto.ErrEmptyPassword,
			},
			{
				name:       "fails with empty salt",
				password:   validPassword,
				salt:       []byte{},
				iterations: validIterations,
				keyLen:     validKeyLen,
				wantErr:    crypto.ErrInvalidSalt,
			},
			{
				name:       "fails with nil salt",
				password:   validPassword,
				salt:       nil,
				iterations: validIterations,
				keyLen:     validKeyLen,
				wantErr:    crypto.ErrInvalidSalt,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				key, err := crypto.DeriveKey(tt.password, tt.salt, tt.iterations, tt.keyLen)
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, key)
			})
		}
	})

	t.Run("key derivation properties", func(t *testing.T) {
		t.Run("different passwords produce different keys", func(t *testing.T) {
			key1, err := crypto.DeriveKey([]byte("password1"), validSalt, validIterations, validKeyLen)
			require.NoError(t, err)

			key2, err := crypto.DeriveKey([]byte("password2"), validSalt, validIterations, validKeyLen)
			require.NoError(t, err)

			assert.False(t, bytes.Equal(key1, key2))
		})

		t.Run("different iterations produce different keys", func(t *testing.T) {
			key1, err := crypto.DeriveKey(validPassword, validSalt, 1000, validKeyLen)
			require.NoError(t, err)

			key2, err := crypto.DeriveKey(validPassword, validSalt, 2000, validKeyLen)
			require.NoError(t, err)

			assert.False(t, bytes.Equal(key1, key2))
		})
	})
}
