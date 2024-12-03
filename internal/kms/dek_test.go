package kms_test

import (
	"testing"
	"time"

	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockEncrypter implements crypto.Encrypter for testing
type mockEncrypter struct {
	encryptFunc func([]byte) (*crypto.Ciphertext, error)
}

func (m *mockEncrypter) Encrypt(data []byte) (*crypto.Ciphertext, error) {
	return m.encryptFunc(data)
}

// mockDecrypter implements crypto.Decrypter for testing
type mockDecrypter struct {
	decryptFunc func(*crypto.Ciphertext) ([]byte, error)
}

func (m *mockDecrypter) Decrypt(ct *crypto.Ciphertext) ([]byte, error) {
	return m.decryptFunc(ct)
}

// generateTestDEK creates a DEK for testing
func generateTestDEK(t *testing.T) *kms.UnsealedDEK {
	t.Helper()
	dek, err := kms.GenerateDEK()
	require.NoError(t, err)
	require.NotNil(t, dek)
	return dek
}

func TestDEK_Generate(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("generates valid key", func(t *testing.T) {
			dek := generateTestDEK(t)
			defer dek.Clear()

			// Verify can perform encryption
			ct, err := dek.Encrypt([]byte("test data"))
			assert.NoError(t, err)
			assert.NotNil(t, ct)

			// Verify can perform decryption
			pt, err := dek.Decrypt(ct)
			assert.NoError(t, err)
			assert.Equal(t, []byte("test data"), pt)
		})

		t.Run("generates unique keys", func(t *testing.T) {
			dek1 := generateTestDEK(t)
			defer dek1.Clear()

			dek2 := generateTestDEK(t)
			defer dek2.Clear()

			data := []byte("test data")
			ct1, err := dek1.Encrypt(data)
			require.NoError(t, err)
			ct2, err := dek2.Encrypt(data)
			require.NoError(t, err)

			// Ciphertexts should be different due to different keys
			assert.NotEqual(t, ct1.Data, ct2.Data)
		})
	})
}

func TestDEK_InitialSeal(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("seals key correctly", func(t *testing.T) {
			dek := generateTestDEK(t)
			defer dek.Clear()

			mockEnc := &mockEncrypter{
				encryptFunc: func(data []byte) (*crypto.Ciphertext, error) {
					return &crypto.Ciphertext{
						Nonce: make([]byte, crypto.NonceSize),
						Data:  []byte("encrypted"),
						Tag:   make([]byte, crypto.TagSize),
					}, nil
				},
			}

			initial, err := dek.InitialSeal(mockEnc)
			assert.NoError(t, err)
			assert.NotNil(t, initial)
			assert.Equal(t, crypto.AlgorithmAESGCM256, initial.Algorithm)
			assert.Equal(t, crypto.KeyStateActive, initial.State)
			assert.NotNil(t, initial.Key)
			assert.WithinDuration(t, time.Now(), initial.Created, time.Second)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil encrypter", func(t *testing.T) {
			dek := generateTestDEK(t)
			defer dek.Clear()

			initial, err := dek.InitialSeal(nil)
			assert.ErrorIs(t, err, kms.ErrNilEncrypter)
			assert.Nil(t, initial)
		})

		t.Run("handles encryption failure", func(t *testing.T) {
			dek := generateTestDEK(t)
			defer dek.Clear()

			mockEnc := &mockEncrypter{
				encryptFunc: func(data []byte) (*crypto.Ciphertext, error) {
					return nil, crypto.ErrAESEncryption
				},
			}

			initial, err := dek.InitialSeal(mockEnc)
			assert.Error(t, err)
			assert.Nil(t, initial)
		})
	})
}

func TestDEK_Unseal(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("unseals key correctly", func(t *testing.T) {
			storedDEK := &kms.DEK{
				ID:        "test",
				Algorithm: crypto.AlgorithmAESGCM256,
				Key: &crypto.Ciphertext{
					Nonce: make([]byte, crypto.NonceSize),
					Data:  []byte("encrypted"),
					Tag:   make([]byte, crypto.TagSize),
				},
				Created: time.Now(),
				State:   crypto.KeyStateActive,
			}

			mockDec := &mockDecrypter{
				decryptFunc: func(ct *crypto.Ciphertext) ([]byte, error) {
					key := make([]byte, crypto.AESKeySize)
					return key, nil
				},
			}

			unsealed, err := storedDEK.Unseal(mockDec)
			assert.NoError(t, err)
			assert.NotNil(t, unsealed)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil decrypter", func(t *testing.T) {
			storedDEK := &kms.DEK{
				State: crypto.KeyStateActive,
			}

			unsealed, err := storedDEK.Unseal(nil)
			assert.ErrorIs(t, err, kms.ErrNilDecrypter)
			assert.Nil(t, unsealed)
		})

		t.Run("rejects destroyed key", func(t *testing.T) {
			storedDEK := &kms.DEK{
				State: crypto.KeyStateDestroyed,
			}

			mockDec := &mockDecrypter{
				decryptFunc: func(ct *crypto.Ciphertext) ([]byte, error) {
					return nil, nil
				},
			}

			unsealed, err := storedDEK.Unseal(mockDec)
			assert.ErrorIs(t, err, kms.ErrKeyRevoked)
			assert.Nil(t, unsealed)
		})

		t.Run("handles incorrect key size", func(t *testing.T) {
			storedDEK := &kms.DEK{
				State: crypto.KeyStateActive,
				Key:   &crypto.Ciphertext{},
			}

			mockDec := &mockDecrypter{
				decryptFunc: func(ct *crypto.Ciphertext) ([]byte, error) {
					return []byte("wrong size"), nil
				},
			}

			unsealed, err := storedDEK.Unseal(mockDec)
			assert.ErrorIs(t, err, kms.ErrInvalidDEK)
			assert.Nil(t, unsealed)
		})
	})
}

func TestDEK_Encrypt(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("encrypts data correctly", func(t *testing.T) {
			dek := generateTestDEK(t)
			defer dek.Clear()

			data := []byte("test data")
			ct, err := dek.Encrypt(data)
			assert.NoError(t, err)
			assert.NotNil(t, ct)
			assert.Len(t, ct.Nonce, crypto.NonceSize)
			assert.NotEmpty(t, ct.Data)
			assert.Len(t, ct.Tag, crypto.TagSize)

			// Verify decryption works
			pt, err := dek.Decrypt(ct)
			assert.NoError(t, err)
			assert.Equal(t, data, pt)
		})

		t.Run("handles empty data", func(t *testing.T) {
			dek := generateTestDEK(t)
			defer dek.Clear()

			ct, err := dek.Encrypt([]byte{})
			assert.NoError(t, err)
			assert.NotNil(t, ct)

			pt, err := dek.Decrypt(ct)
			assert.NoError(t, err)
			assert.Equal(t, []byte{}, pt)
		})
	})
}

func TestDEK_Decrypt(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("decrypts data correctly", func(t *testing.T) {
			dek := generateTestDEK(t)
			defer dek.Clear()

			data := []byte("test data")
			ct, err := dek.Encrypt(data)
			require.NoError(t, err)

			pt, err := dek.Decrypt(ct)
			assert.NoError(t, err)
			assert.Equal(t, data, pt)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects nil ciphertext", func(t *testing.T) {
			dek := generateTestDEK(t)
			defer dek.Clear()

			pt, err := dek.Decrypt(nil)
			assert.Error(t, err)
			assert.Nil(t, pt)
		})

		t.Run("handles invalid ciphertext", func(t *testing.T) {
			dek := generateTestDEK(t)
			defer dek.Clear()

			ct := &crypto.Ciphertext{
				Nonce: make([]byte, crypto.NonceSize),
				Data:  []byte("invalid"),
				Tag:   make([]byte, crypto.TagSize),
			}

			pt, err := dek.Decrypt(ct)
			assert.Error(t, err)
			assert.Nil(t, pt)
		})
	})
}

func TestDEK_Clear(t *testing.T) {
	t.Run("clears key material", func(t *testing.T) {
		dek := generateTestDEK(t)
		
		// Verify key works before clearing
		_, err := dek.Encrypt([]byte("test"))
		require.NoError(t, err)

		dek.Clear()

		// Verify key doesn't work after clearing
		_, err = dek.Encrypt([]byte("test"))
		assert.ErrorIs(t, err, kms.ErrInvalidDEK)
	})

	t.Run("safe to call multiple times", func(t *testing.T) {
		dek := generateTestDEK(t)
		dek.Clear()
		dek.Clear() // Should not panic
	})
}
