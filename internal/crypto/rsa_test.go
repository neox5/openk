package crypto_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateTestKey(t *testing.T, bits int) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, bits)
	require.NoError(t, err)
	return key
}

func TestRSA_Generate(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		tests := []struct {
			name string
			bits int
		}{
			{
				name: "generates 2048-bit key",
				bits: crypto.RSAKeySize2048,
			},
			{
				name: "generates 4096-bit key",
				bits: crypto.RSAKeySize4096,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				key, err := crypto.GenerateRSAKeyPair(tt.bits)
				assert.NoError(t, err)
				assert.Equal(t, tt.bits, key.Size()*8)
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects undersized key", func(t *testing.T) {
			key, err := crypto.GenerateRSAKeyPair(1024)
			assert.ErrorIs(t, err, crypto.ErrInvalidKeySize)
			assert.Nil(t, key)
		})
	})
}

func TestRSA_ExportPrivateKey(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		key := generateTestKey(t, crypto.RSAKeySize2048)

		der, err := crypto.ExportRSAPrivateKey(key)
		assert.NoError(t, err)
		assert.NotEmpty(t, der)

		parsed, err := x509.ParsePKCS8PrivateKey(der)
		assert.NoError(t, err)
		assert.IsType(t, &rsa.PrivateKey{}, parsed)
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil key", func(t *testing.T) {
			der, err := crypto.ExportRSAPrivateKey(nil)
			assert.ErrorIs(t, err, crypto.ErrInvalidPrivateKey)
			assert.Nil(t, der)
		})
	})
}

func TestRSA_ExportPublicKey(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		key := &generateTestKey(t, crypto.RSAKeySize2048).PublicKey

		der, err := crypto.ExportRSAPublicKey(key)
		assert.NoError(t, err)
		assert.NotEmpty(t, der)

		parsed, err := x509.ParsePKIXPublicKey(der)
		assert.NoError(t, err)
		assert.IsType(t, &rsa.PublicKey{}, parsed)
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil key", func(t *testing.T) {
			der, err := crypto.ExportRSAPublicKey(nil)
			assert.ErrorIs(t, err, crypto.ErrInvalidPublicKey)
			assert.Nil(t, der)
		})
	})
}

func TestRSA_ImportPrivateKey(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		original := generateTestKey(t, crypto.RSAKeySize2048)
		der, err := x509.MarshalPKCS8PrivateKey(original)
		require.NoError(t, err)

		imported, err := crypto.ImportRSAPrivateKey(der)
		assert.NoError(t, err)
		assert.Equal(t, original.D.Bytes(), imported.D.Bytes())
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects invalid DER", func(t *testing.T) {
			imported, err := crypto.ImportRSAPrivateKey([]byte("invalid"))
			assert.ErrorIs(t, err, crypto.ErrInvalidPrivateKey)
			assert.Contains(t, err.Error(), "invalid RSA private key")
			assert.Nil(t, imported)
		})

		t.Run("rejects EC key", func(t *testing.T) {
			ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			require.NoError(t, err)
			der, err := x509.MarshalPKCS8PrivateKey(ecKey)
			require.NoError(t, err)

			imported, err := crypto.ImportRSAPrivateKey(der)
			assert.ErrorIs(t, err, crypto.ErrInvalidPrivateKey)
			assert.Contains(t, err.Error(), "not an RSA private key")
			assert.Nil(t, imported)
		})

		t.Run("rejects undersized key", func(t *testing.T) {
			key, err := rsa.GenerateKey(rand.Reader, 1024)
			require.NoError(t, err)
			der, err := x509.MarshalPKCS8PrivateKey(key)
			require.NoError(t, err)

			imported, err := crypto.ImportRSAPrivateKey(der)
			assert.ErrorIs(t, err, crypto.ErrInvalidKeySize)
			assert.Contains(t, err.Error(), "got 1024 bits")
			assert.Nil(t, imported)
		})
	})
}

func TestRSA_ImportPublicKey(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		original := &generateTestKey(t, crypto.RSAKeySize2048).PublicKey
		der, err := x509.MarshalPKIXPublicKey(original)
		require.NoError(t, err)

		imported, err := crypto.ImportRSAPublicKey(der)
		assert.NoError(t, err)
		assert.Equal(t, original.N.Bytes(), imported.N.Bytes())
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects invalid DER", func(t *testing.T) {
			imported, err := crypto.ImportRSAPublicKey([]byte("invalid"))
			assert.ErrorIs(t, err, crypto.ErrInvalidPublicKey)
			assert.Contains(t, err.Error(), "invalid RSA public key")
			assert.Nil(t, imported)
		})

		t.Run("rejects EC key", func(t *testing.T) {
			ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			require.NoError(t, err)
			der, err := x509.MarshalPKIXPublicKey(&ecKey.PublicKey)
			require.NoError(t, err)

			imported, err := crypto.ImportRSAPublicKey(der)
			assert.ErrorIs(t, err, crypto.ErrInvalidPublicKey)
			assert.Contains(t, err.Error(), "not an RSA public key")
			assert.Nil(t, imported)
		})

		t.Run("rejects undersized key", func(t *testing.T) {
			key, err := rsa.GenerateKey(rand.Reader, 1024)
			require.NoError(t, err)
			der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
			require.NoError(t, err)

			imported, err := crypto.ImportRSAPublicKey(der)
			assert.ErrorIs(t, err, crypto.ErrInvalidKeySize)
			assert.Contains(t, err.Error(), "got 1024 bits")
			assert.Nil(t, imported)
		})
	})
}

func TestRSA_Encrypt(t *testing.T) {
	key := generateTestKey(t, crypto.RSAKeySize2048)
	maxSize := key.Size() - 2*sha256.Size - 2
	message := []byte("test message")

	t.Run("success cases", func(t *testing.T) {
		tests := []struct {
			name    string
			message []byte
		}{
			{
				name:    "encrypts short message",
				message: message,
			},
			{
				name:    "encrypts empty message",
				message: []byte{},
			},
			{
				name:    "encrypts binary data",
				message: []byte{0xFF, 0x00, 0xFF, 0x00},
			},
			{
				name:    "encrypts maximum size message",
				message: bytes.Repeat([]byte("a"), maxSize),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ct, err := crypto.RSAEncrypt(&key.PublicKey, tt.message)
				assert.NoError(t, err)
				assert.NotNil(t, ct)

				plaintext, err := crypto.RSADecrypt(key, ct.Data)
				assert.NoError(t, err)
				assert.Equal(t, tt.message, plaintext)
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil key", func(t *testing.T) {
			ct, err := crypto.RSAEncrypt(nil, message)
			assert.ErrorIs(t, err, crypto.ErrInvalidPublicKey)
			assert.Nil(t, ct)
		})

		t.Run("rejects oversized message", func(t *testing.T) {
			ct, err := crypto.RSAEncrypt(&key.PublicKey, bytes.Repeat([]byte("a"), maxSize+1))
			assert.Error(t, err)
			assert.ErrorIs(t, err, rsa.ErrMessageTooLong)
			assert.Contains(t, err.Error(), "RSA encryption failed")
			assert.Nil(t, ct)
		})
	})
}

func TestRSA_Decrypt(t *testing.T) {
	key := generateTestKey(t, crypto.RSAKeySize2048)
	message := []byte("test message")

	t.Run("success cases", func(t *testing.T) {
		ct, err := crypto.RSAEncrypt(&key.PublicKey, message)
		require.NoError(t, err)

		plaintext, err := crypto.RSADecrypt(key, ct.Data)
		assert.NoError(t, err)
		assert.Equal(t, message, plaintext)
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil key", func(t *testing.T) {
			plaintext, err := crypto.RSADecrypt(nil, []byte("data"))
			assert.ErrorIs(t, err, crypto.ErrInvalidPrivateKey)
			assert.Nil(t, plaintext)
		})

		t.Run("handles invalid ciphertext", func(t *testing.T) {
			plaintext, err := crypto.RSADecrypt(key, []byte("invalid"))
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "RSA decryption failed")
			assert.Nil(t, plaintext)
		})
	})
}
