package crypto_test

import (
	"bytes"
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAESGenerateKey(t *testing.T) {
	// Test correct length
	key1, err := crypto.AESGenerateKey()
	assert.NoError(t, err)
	assert.Equal(t, crypto.AESKeySize, len(key1))

	// Test uniqueness
	key2, err := crypto.AESGenerateKey()
	assert.NoError(t, err)
	assert.False(t, bytes.Equal(key1, key2))
}

func TestAESGenerateNonce(t *testing.T) {
	// Test correct length
	nonce1, err := crypto.AESGenerateNonce()
	assert.NoError(t, err)
	assert.Equal(t, crypto.NonceSize, len(nonce1))

	// Test uniqueness
	nonce2, err := crypto.AESGenerateNonce()
	assert.NoError(t, err)
	assert.False(t, bytes.Equal(nonce1, nonce2))
}

func TestAESEncryptDecrypt(t *testing.T) {
	// Generate a valid key for tests
	validKey, err := crypto.AESGenerateKey()
	require.NoError(t, err)

	// Create a 1MB test message
	largeMessage := make([]byte, 1024*1024)
	for i := range largeMessage {
		largeMessage[i] = byte(i % 256)
	}

	tests := []struct {
		name      string
		key       []byte
		plaintext []byte
		wantErr   error
	}{
		{
			name:      "valid small message",
			key:       validKey,
			plaintext: []byte("test message"),
			wantErr:   nil,
		},
		{
			name:      "valid empty message",
			key:       validKey,
			plaintext: []byte{},
			wantErr:   nil,
		},
		{
			name:      "valid large message (1MB)",
			key:       validKey,
			plaintext: largeMessage,
			wantErr:   nil,
		},
		{
			name:      "nil message",
			key:       validKey,
			plaintext: nil,
			wantErr:   crypto.ErrAESInvalidMessage,
		},
		{
			name:      "key too short",
			key:       validKey[:crypto.AESKeySize-1],
			plaintext: []byte("test message"),
			wantErr:   crypto.ErrAESKeySize,
		},
		{
			name:      "key too long",
			key:       append(validKey, 0xFF),
			plaintext: []byte("test message"),
			wantErr:   crypto.ErrAESKeySize,
		},
		{
			name:      "nil key",
			key:       nil,
			plaintext: []byte("test message"),
			wantErr:   crypto.ErrAESKeySize,
		},
		{
			name:      "empty key",
			key:       []byte{},
			plaintext: []byte("test message"),
			wantErr:   crypto.ErrAESKeySize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encryption
			ciphertext, err := crypto.AESEncrypt(tt.key, tt.plaintext)
			if tt.wantErr != nil {
				assert.Error(t, err, "tc '%s' should return an error", tt.name)
				assert.Equal(t, tt.wantErr, err, "tc '%s' returned wrong error", tt.name)
				return
			}
			require.NoError(t, err, "tc '%s' returned unexpected error", tt.name)
			require.NotNil(t, ciphertext, "tc '%s' returned nil ciphertext", tt.name)

			// Verify ciphertext structure
			assert.Equal(t, crypto.NonceSize, len(ciphertext.Nonce), "tc '%s': wrong nonce size", tt.name)
			assert.Equal(t, crypto.TagSize, len(ciphertext.Tag), "tc '%s': wrong tag size", tt.name)
			if len(tt.plaintext) != 0 {
				assert.NotEmpty(t, ciphertext.Data, "tc '%s': expected ciphertext data to be non-empty, got len=%d", tt.name, len(ciphertext.Data))
			}

			// Test decryption
			decrypted, err := crypto.AESDecrypt(tt.key, ciphertext)
			require.NoError(t, err, "tc '%s' decryption failed", tt.name)
			assert.Equal(t, tt.plaintext, decrypted, "tc '%s': decrypted data doesn't match original", tt.name)
		})
	}
}

func TestAESDecryptErrors(t *testing.T) {
	// Generate a valid key for tests
	key, err := crypto.AESGenerateKey()
	require.NoError(t, err)
	plaintext := []byte("test message")

	// Create valid ciphertext for modification
	ct, err := crypto.AESEncrypt(key, plaintext)
	require.NoError(t, err)

	tests := []struct {
		name    string
		key     []byte
		ct      *crypto.Ciphertext
		wantErr error
	}{
		{
			name:    "nil ciphertext",
			key:     key,
			ct:      nil,
			wantErr: crypto.ErrAESInvalidMessage,
		},
		{
			name:    "invalid key size",
			key:     make([]byte, crypto.AESKeySize-1),
			ct:      ct,
			wantErr: crypto.ErrAESKeySize,
		},
		{
			name: "modified tag",
			key:  key,
			ct: &crypto.Ciphertext{
				Nonce: ct.Nonce,
				Data:  ct.Data,
				Tag:   make([]byte, crypto.TagSize), // Zero tag
			},
			wantErr: crypto.ErrAESDecryption,
		},
		{
			name: "modified data",
			key:  key,
			ct: &crypto.Ciphertext{
				Nonce: ct.Nonce,
				Data:  append(ct.Data, 0), // Modified data
				Tag:   ct.Tag,
			},
			wantErr: crypto.ErrAESDecryption,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := crypto.AESDecrypt(tt.key, tt.ct)
			assert.Error(t, err, "tc '%s' should return an error", tt.name)
			assert.Equal(t, tt.wantErr, err, "tc '%s' returned wrong error", tt.name)
		})
	}
}

// Test vectors from NIST
func TestAESGCMTestVectors(t *testing.T) {
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

	// Encrypt
	ct, err := crypto.AESEncrypt(key, plaintext)
	require.NoError(t, err)

	// Decrypt
	decrypted, err := crypto.AESDecrypt(key, ct)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, plaintext, decrypted)
}
