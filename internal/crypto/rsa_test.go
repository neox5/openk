package crypto_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRSA_Generate(t *testing.T) {
	tests := []struct {
		name    string
		bits    int
		wantErr error
	}{
		{
			name:    "valid RSAKeySize2048 key",
			bits:    crypto.RSAKeySize2048,
			wantErr: nil,
		},
		{
			name:    "valid RSAKeySize4096 key",
			bits:    crypto.RSAKeySize4096,
			wantErr: nil,
		},
		{
			name:    "invalid 1024-bit key",
			bits:    1024,
			wantErr: crypto.ErrInvalidKeySize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := crypto.GenerateRSAKeyPair(tt.bits)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, key)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, key)
			assert.Equal(t, tt.bits, key.Size()*8)
		})
	}
}

func TestRSA_Export(t *testing.T) {
	privKey, err := crypto.GenerateRSAKeyPair(crypto.RSAKeySize2048)
	require.NoError(t, err)

	t.Run("private key with nil", func(t *testing.T) {
		der, err := crypto.ExportRSAPrivateKey(nil)
		assert.Error(t, err)
		assert.Equal(t, crypto.ErrInvalidPrivateKey, err)
		assert.Nil(t, der)
	})

	t.Run("public key with nil", func(t *testing.T) {
		der, err := crypto.ExportRSAPublicKey(nil)
		assert.Error(t, err)
		assert.Equal(t, crypto.ErrInvalidPublicKey, err)
		assert.Nil(t, der)
	})

	t.Run("valid keys", func(t *testing.T) {
		privDer, err := crypto.ExportRSAPrivateKey(privKey)
		assert.NoError(t, err)
		assert.NotNil(t, privDer)

		pubDer, err := crypto.ExportRSAPublicKey(&privKey.PublicKey)
		assert.NoError(t, err)
		assert.NotNil(t, pubDer)
	})
}

func TestRSA_Import(t *testing.T) {
	privKey, err := crypto.GenerateRSAKeyPair(crypto.RSAKeySize2048)
	require.NoError(t, err)

	t.Run("private key import/export", func(t *testing.T) {
		privDer, err := crypto.ExportRSAPrivateKey(privKey)
		assert.NoError(t, err)
		assert.NotNil(t, privDer)

		imported, err := crypto.ImportRSAPrivateKey(privDer)
		assert.NoError(t, err)
		assert.NotNil(t, imported)
		assert.Equal(t, privKey.D.Bytes(), imported.D.Bytes())
		assert.Equal(t, privKey.N.Bytes(), imported.N.Bytes())
	})

	t.Run("public key import/export", func(t *testing.T) {
		pubKey := &privKey.PublicKey
		pubDer, err := crypto.ExportRSAPublicKey(pubKey)
		assert.NoError(t, err)
		assert.NotNil(t, pubDer)

		imported, err := crypto.ImportRSAPublicKey(pubDer)
		assert.NoError(t, err)
		assert.NotNil(t, imported)
		assert.Equal(t, pubKey.N.Bytes(), imported.N.Bytes())
		assert.Equal(t, pubKey.E, imported.E)
	})
}

func TestRSA_Import_InvalidInput(t *testing.T) {
	t.Run("invalid private key DER", func(t *testing.T) {
		// Invalid DER bytes
		invalidDer := []byte("not a valid DER")
		key, err := crypto.ImportRSAPrivateKey(invalidDer)
		assert.Error(t, err)
		assert.Nil(t, key)
	})

	t.Run("invalid public key DER", func(t *testing.T) {
		// Invalid DER bytes
		invalidDer := []byte("not a valid DER")
		key, err := crypto.ImportRSAPublicKey(invalidDer)
		assert.Error(t, err)
		assert.Nil(t, key)
	})

	t.Run("wrong key type", func(t *testing.T) {
		// Generate an EC key pair instead of RSA
		ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err)

		// Export EC private key
		ecDer, err := x509.MarshalPKCS8PrivateKey(ecKey)
		require.NoError(t, err)

		// Try to import as RSA key
		key, err := crypto.ImportRSAPrivateKey(ecDer)
		assert.Error(t, err)
		assert.Equal(t, crypto.ErrInvalidPrivateKey, err)
		assert.Nil(t, key)

		// Export EC public key
		ecPubDer, err := x509.MarshalPKIXPublicKey(&ecKey.PublicKey)
		require.NoError(t, err)

		// Try to import as RSA key
		pubKey, err := crypto.ImportRSAPublicKey(ecPubDer)
		assert.Error(t, err)
		assert.Equal(t, crypto.ErrInvalidPublicKey, err)
		assert.Nil(t, pubKey)
	})

	t.Run("undersized key import", func(t *testing.T) {
		// Generate a 1024-bit key directly using crypto/rsa for testing
		smallKey, err := rsa.GenerateKey(rand.Reader, 1024)
		require.NoError(t, err)

		// Try to import undersized private key
		smallPrivDer, err := x509.MarshalPKCS8PrivateKey(smallKey)
		require.NoError(t, err)

		key, err := crypto.ImportRSAPrivateKey(smallPrivDer)
		assert.Error(t, err)
		assert.Equal(t, crypto.ErrInvalidKeySize, err)
		assert.Nil(t, key)

		// Try to import undersized public key
		smallPubDer, err := x509.MarshalPKIXPublicKey(&smallKey.PublicKey)
		require.NoError(t, err)

		pubKey, err := crypto.ImportRSAPublicKey(smallPubDer)
		assert.Error(t, err)
		assert.Equal(t, crypto.ErrInvalidKeySize, err)
		assert.Nil(t, pubKey)
	})
}

func TestRSA_Encrypt(t *testing.T) {
	privKey, err := crypto.GenerateRSAKeyPair(crypto.RSAKeySize2048)
	require.NoError(t, err)
	pubKey := &privKey.PublicKey

	t.Run("valid encryption", func(t *testing.T) {
		message := []byte("test message")
		ct, err := crypto.RSAEncrypt(pubKey, message)
		assert.NoError(t, err)
		require.NotNil(t, ct)

		// Verify ciphertext structure
		assert.Equal(t, crypto.NonceSize, len(ct.Nonce))
		assert.Equal(t, crypto.TagSize, len(ct.Tag))
		assert.NotEmpty(t, ct.Data)

		// Verify decryption
		decrypted, err := crypto.RSADecrypt(privKey, ct.Data)
		assert.NoError(t, err)
		assert.Equal(t, message, decrypted)
	})

	t.Run("empty message", func(t *testing.T) {
		ct, err := crypto.RSAEncrypt(pubKey, []byte{})
		assert.NoError(t, err)
		require.NotNil(t, ct)

		decrypted, err := crypto.RSADecrypt(privKey, ct.Data)
		assert.NoError(t, err)
		assert.Empty(t, decrypted)
	})

	t.Run("nil message", func(t *testing.T) {
		ct, err := crypto.RSAEncrypt(pubKey, nil)
		assert.NoError(t, err)
		require.NotNil(t, ct)

		decrypted, err := crypto.RSADecrypt(privKey, ct.Data)
		assert.NoError(t, err)
		assert.Empty(t, decrypted)
	})

	t.Run("nil public key", func(t *testing.T) {
		ct, err := crypto.RSAEncrypt(nil, []byte("test"))
		assert.Error(t, err)
		assert.Equal(t, crypto.ErrInvalidPublicKey, err)
		assert.Nil(t, ct)
	})

	t.Run("message size limits", func(t *testing.T) {
		// Calculate maximum message size for RSA-OAEP
		maxSize := privKey.Size() - 2*crypto.AESKeySize - 2

		// Test maximum valid size
		message := bytes.Repeat([]byte("a"), maxSize)
		ct, err := crypto.RSAEncrypt(pubKey, message)
		assert.NoError(t, err)
		require.NotNil(t, ct)

		decrypted, err := crypto.RSADecrypt(privKey, ct.Data)
		assert.NoError(t, err)
		assert.Equal(t, message, decrypted)

		// Test exceeding maximum size
		message = bytes.Repeat([]byte("a"), maxSize+1)
		ct, err = crypto.RSAEncrypt(pubKey, message)
		assert.Error(t, err)
		assert.Nil(t, ct)
	})
}

func TestRSA_Decrypt(t *testing.T) {
	privKey, err := crypto.GenerateRSAKeyPair(crypto.RSAKeySize2048)
	require.NoError(t, err)
	pubKey := &privKey.PublicKey

	t.Run("valid decryption", func(t *testing.T) {
		message := []byte("test message")
		ct, err := crypto.RSAEncrypt(pubKey, message)
		require.NoError(t, err)

		decrypted, err := crypto.RSADecrypt(privKey, ct.Data)
		assert.NoError(t, err)
		assert.Equal(t, message, decrypted)
	})

	t.Run("nil private key", func(t *testing.T) {
		decrypted, err := crypto.RSADecrypt(nil, []byte("test"))
		assert.Error(t, err)
		assert.Equal(t, crypto.ErrInvalidPrivateKey, err)
		assert.Nil(t, decrypted)
	})

	t.Run("invalid ciphertext", func(t *testing.T) {
		decrypted, err := crypto.RSADecrypt(privKey, []byte("invalid"))
		assert.Error(t, err)
		assert.Nil(t, decrypted)
	})
}
