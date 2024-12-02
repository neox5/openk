package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
)

const (
	// RSAKeySize2048 is the required RSA key size (2048 bits)
	RSAKeySize2048 = 2048

	// RSAKeySize4096 is the extended RSA key size (4096 bits)
	RSAKeySize4096 = 4096
)

var (
	// ErrInvalidKeySize indicates the RSA key size is too small
	ErrInvalidKeySize = errors.New("RSA key size must be at least 2048 bits")
	// ErrInvalidPrivateKey indicates the private key is invalid
	ErrInvalidPrivateKey = errors.New("invalid RSA private key")
	// ErrInvalidPublicKey indicates the public key is invalid
	ErrInvalidPublicKey = errors.New("invalid RSA public key")
	// ErrMessageTooLong indicates the message exceeds OAEP max length
	ErrMessageTooLong = errors.New("message too long for RSA-OAEP encryption")
)

// GenerateRSAKeyPair generates a new RSA key pair
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	if bits < RSAKeySize2048 {
		return nil, ErrInvalidKeySize
	}
	return rsa.GenerateKey(rand.Reader, bits)
}

// ExportRSAPrivateKey exports private key in PKCS#8 DER binary format
func ExportRSAPrivateKey(key *rsa.PrivateKey) (derBytes []byte, err error) {
	if key == nil {
		return nil, ErrInvalidPrivateKey
	}

	// Convert to PKCS#8 DER format
	derBytes, err = x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	return derBytes, nil
}

// ImportRSAPrivateKey imports private key from PKCS#8 DER binary format
func ImportRSAPrivateKey(derBytes []byte) (*rsa.PrivateKey, error) {
	key, err := x509.ParsePKCS8PrivateKey(derBytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidPrivateKey
	}

	if rsaKey.Size()*8 < RSAKeySize2048 {
		return nil, ErrInvalidKeySize
	}

	return rsaKey, nil
}

// ExportRSAPublicKey exports public key in SPKI/X.509 DER binary format
func ExportRSAPublicKey(key *rsa.PublicKey) (derBytes []byte, err error) {
	if key == nil {
		return nil, ErrInvalidPublicKey
	}

	derBytes, err = x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	return derBytes, nil
}

// ImportRSAPublicKey imports public key from SPKI/X.509 DER binary format
func ImportRSAPublicKey(derBytes []byte) (*rsa.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(derBytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidPublicKey
	}

	if rsaKey.Size()*8 < RSAKeySize2048 {
		return nil, ErrInvalidKeySize
	}

	return rsaKey, nil
}

// RSAEncrypt encrypts a message using RSA-OAEP-SHA256
func RSAEncrypt(key *rsa.PublicKey, message []byte) (*Ciphertext, error) {
	if key == nil {
		return nil, ErrInvalidPublicKey
	}

	encData, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		key,
		message,
		nil, // No label used
	)
	if err != nil {
		return nil, err
	}

	// RSA-OAEP doesn't use nonce/tag, but we maintain Ciphertext structure
	// with all data in the Data field for consistency
	return NewCiphertext(
		make([]byte, NonceSize), // Empty nonce
		encData,                 // Encrypted data
		make([]byte, TagSize),   // Empty tag
	)
}

// RSADecrypt decrypts a message using RSA-OAEP-SHA256
func RSADecrypt(key *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	if key == nil {
		return nil, ErrInvalidPrivateKey
	}

	return rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		key,
		ciphertext,
		nil, // No label used
	)
}
