package crypto

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "errors"
    "fmt"
)

const (
    RSAKeySize2048 = 2048
    RSAKeySize4096 = 4096
)

var (
    ErrInvalidKeySize = errors.New("RSA key size must be at least 2048 bits")
    ErrInvalidPrivateKey = errors.New("invalid RSA private key")
    ErrInvalidPublicKey = errors.New("invalid RSA public key")
)

func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
    if bits < RSAKeySize2048 {
        return nil, ErrInvalidKeySize
    }
    key, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
    }
    return key, nil
}

func ExportRSAPrivateKey(key *rsa.PrivateKey) ([]byte, error) {
    if key == nil {
        return nil, ErrInvalidPrivateKey
    }
    der, err := x509.MarshalPKCS8PrivateKey(key)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal private key: %w", err)
    }
    return der, nil
}

func ImportRSAPrivateKey(derBytes []byte) (*rsa.PrivateKey, error) {
    key, err := x509.ParsePKCS8PrivateKey(derBytes)
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrInvalidPrivateKey, err)
    }

    rsaKey, ok := key.(*rsa.PrivateKey)
    if !ok {
        return nil, fmt.Errorf("%w: not an RSA private key", ErrInvalidPrivateKey)
    }

    if rsaKey.Size()*8 < RSAKeySize2048 {
        return nil, fmt.Errorf("%w: got %d bits", ErrInvalidKeySize, rsaKey.Size()*8)
    }

    return rsaKey, nil
}

func ExportRSAPublicKey(key *rsa.PublicKey) ([]byte, error) {
    if key == nil {
        return nil, ErrInvalidPublicKey
    }
    der, err := x509.MarshalPKIXPublicKey(key)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal public key: %w", err)
    }
    return der, nil
}

func ImportRSAPublicKey(derBytes []byte) (*rsa.PublicKey, error) {
    key, err := x509.ParsePKIXPublicKey(derBytes)
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrInvalidPublicKey, err)
    }

    rsaKey, ok := key.(*rsa.PublicKey)
    if !ok {
        return nil, fmt.Errorf("%w: not an RSA public key", ErrInvalidPublicKey)
    }

    if rsaKey.Size()*8 < RSAKeySize2048 {
        return nil, fmt.Errorf("%w: got %d bits", ErrInvalidKeySize, rsaKey.Size()*8)
    }

    return rsaKey, nil
}

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
        return nil, fmt.Errorf("RSA encryption failed: %w", err)
    }

    return NewCiphertext(
        make([]byte, NonceSize),
        encData,
        make([]byte, TagSize),
    )
}

func RSADecrypt(key *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
    if key == nil {
        return nil, ErrInvalidPrivateKey
    }

    plaintext, err := rsa.DecryptOAEP(
        sha256.New(),
        rand.Reader,
        key,
        ciphertext,
        nil, // No label used
    )
    if err != nil {
        return nil, fmt.Errorf("RSA decryption failed: %w", err)
    }
    return plaintext, nil
}
