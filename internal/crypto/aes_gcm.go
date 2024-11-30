package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

var (
	ErrAESKeySize        = errors.New("AES key must be 32 bytes (256 bits)")
	ErrAESEncryption     = errors.New("AES encryption failed")
	ErrAESDecryption     = errors.New("AES decryption failed")
	ErrAESInvalidMessage = errors.New("message cannot be nil")
)

const (
	AESKeySize = 32 // AES-256 key size in bytes
)

// AESGenerateKey generates a new random AES-256 key
func AESGenerateKey() ([]byte, error) {
	key := make([]byte, AESKeySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// AESGenerateNonce generates a new random nonce for AES-GCM
func AESGenerateNonce() ([]byte, error) {
	nonce := make([]byte, NonceSize)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

// AESEncrypt encrypts data using AES-256-GCM
func AESEncrypt(key, plaintext []byte) (*Ciphertext, error) {
	if len(key) != AESKeySize {
		return nil, ErrAESKeySize
	}
	if plaintext == nil {
		return nil, ErrAESInvalidMessage
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate nonce
	nonce, err := AESGenerateNonce()
	if err != nil {
		return nil, err
	}

	// Encrypt and seal
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	if len(ciphertext) < TagSize {
		return nil, ErrAESEncryption
	}

	// Split ciphertext and tag
	tag := ciphertext[len(ciphertext)-TagSize:]
	data := ciphertext[:len(ciphertext)-TagSize]
	
	// Ensure data is an empty slice rather than nil for zero-length inputs
	if data == nil {
		data = []byte{}
	}

	return NewCiphertext(nonce, data, tag)
}

// AESDecrypt decrypts data using AES-256-GCM
func AESDecrypt(key []byte, ct *Ciphertext) ([]byte, error) {
	if len(key) != AESKeySize {
		return nil, ErrAESKeySize
	}
	if ct == nil {
		return nil, ErrAESInvalidMessage
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Combine ciphertext and tag for Open
	ciphertext := make([]byte, len(ct.Data)+TagSize)
	copy(ciphertext, ct.Data)
	copy(ciphertext[len(ct.Data):], ct.Tag)

	// Decrypt and verify
	plaintext, err := gcm.Open(nil, ct.Nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrAESDecryption
	}

	// Ensure we return an empty slice rather than nil for zero-length plaintexts
	if plaintext == nil {
		return []byte{}, nil
	}
	return plaintext, nil
}
