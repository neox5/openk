package crypto_test

import (
	"bytes"
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAES_GenerateKey(t *testing.T) {
	key1, err := crypto.AESGenerateKey()
	require.NoError(t, err)
	assert.Equal(t, crypto.AESKeySize, len(key1))

	key2, err := crypto.AESGenerateKey()
	require.NoError(t, err)
	assert.Equal(t, crypto.AESKeySize, len(key2))

	// Verify keys are unique
	assert.False(t, bytes.Equal(key1, key2))
}

func TestAES_GenerateNonce(t *testing.T) {
	nonce1, err := crypto.AESGenerateNonce()
	require.NoError(t, err)
	assert.Equal(t, crypto.NonceSize, len(nonce1))

	nonce2, err := crypto.AESGenerateNonce()
	require.NoError(t, err)
	assert.Equal(t, crypto.NonceSize, len(nonce2))

	// Verify nonces are unique
	assert.False(t, bytes.Equal(nonce1, nonce2))
}

func TestAES_EncryptDecrypt(t *testing.T) {
	// Generate a valid key for tests
	validKey, err := crypto.AESGenerateKey()
	require.NoError(t, err)

	// Create a 1MB test message
	largeMessage := make([]byte, 1024*1024)
	for i := range largeMessage {
		largeMessage[i] = byte(i % 256)
	}

	t.Run("success cases", func(t *testing.T) {
		tests := []struct {
			name      string
			key       []byte
			plaintext []byte
		}{
			{
				name:      "encrypts and decrypts small message",
				key:       validKey,
				plaintext: []byte("test message"),
			},
			{
				name:      "handles empty message",
				key:       validKey,
				plaintext: []byte{},
			},
			{
				name:      "handles large message (1MB)",
				key:       validKey,
				plaintext: largeMessage,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Test encryption
				ciphertext, err := crypto.AESEncrypt(tt.key, tt.plaintext)
				require.NoError(t, err)
				require.NotNil(t, ciphertext)

				// Verify ciphertext structure
				assert.Equal(t, crypto.NonceSize, len(ciphertext.Nonce))
				assert.Equal(t, crypto.TagSize, len(ciphertext.Tag))
				if len(tt.plaintext) > 0 {
					assert.NotEmpty(t, ciphertext.Data)
				}

				// Test decryption
				decrypted, err := crypto.AESDecrypt(tt.key, ciphertext)
				require.NoError(t, err)
				assert.Equal(t, tt.plaintext, decrypted)
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name      string
			key       []byte
			plaintext []byte
			wantErr   error
		}{
			{
				name:      "fails with nil message",
				key:       validKey,
				plaintext: nil,
				wantErr:   crypto.ErrAESInvalidMessage,
			},
			{
				name:      "fails with key too short",
				key:       validKey[:crypto.AESKeySize-1],
				plaintext: []byte("test message"),
				wantErr:   crypto.ErrAESKeySize,
			},
			{
				name:      "fails with key too long",
				key:       append(validKey, 0xFF),
				plaintext: []byte("test message"),
				wantErr:   crypto.ErrAESKeySize,
			},
			{
				name:      "fails with nil key",
				key:       nil,
				plaintext: []byte("test message"),
				wantErr:   crypto.ErrAESKeySize,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ciphertext, err := crypto.AESEncrypt(tt.key, tt.plaintext)
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, ciphertext)
			})
		}
	})
}

func TestAES_Decrypt(t *testing.T) {
	key, err := crypto.AESGenerateKey()
	require.NoError(t, err)

	validCT, err := crypto.AESEncrypt(key, []byte("test message"))
	require.NoError(t, err)

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name    string
			key     []byte
			ct      *crypto.Ciphertext
			wantErr error
		}{
			{
				name:    "fails with nil ciphertext",
				key:     key,
				ct:      nil,
				wantErr: crypto.ErrAESInvalidMessage,
			},
			{
				name:    "fails with invalid key size",
				key:     make([]byte, crypto.AESKeySize-1),
				ct:      validCT,
				wantErr: crypto.ErrAESKeySize,
			},
			{
				name: "fails with modified tag",
				key:  key,
				ct: &crypto.Ciphertext{
					Nonce: validCT.Nonce,
					Data:  validCT.Data,
					Tag:   make([]byte, crypto.TagSize),
				},
				wantErr: crypto.ErrAESDecryption,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				decrypted, err := crypto.AESDecrypt(tt.key, tt.ct)
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, decrypted)
			})
		}
	})
}

func TestAES_VectorTests(t *testing.T) {
	// Test vector from NIST SP 800-38D
	key := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}
	plaintext := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	}

	ct, err := crypto.AESEncrypt(key, plaintext)
	require.NoError(t, err)

	decrypted, err := crypto.AESDecrypt(key, ct)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}
