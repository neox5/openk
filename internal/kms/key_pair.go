package kms

import (
    "crypto/rsa"
    "fmt"
    "time"
    "github.com/google/uuid"

    "github.com/neox5/openk/internal/crypto"
)

// InitialKeyPair represents a newly generated key pair before server storage
type InitialKeyPair struct {
    Algorithm  crypto.Algorithm
    PublicKey  []byte             // X.509/SPKI format
    PrivateKey *crypto.Ciphertext // Encrypted with protection key
    Created    time.Time
    State      crypto.KeyState
}

// KeyPair represents a key pair as stored in the backend
type KeyPair struct {
    ID         string
    Algorithm  crypto.Algorithm
    PublicKey  []byte
    PrivateKey *crypto.Ciphertext
    Created    time.Time
    State      crypto.KeyState
}

// UnsealedKeyPair represents an active key pair with access to private key operations
type UnsealedKeyPair struct {
    id         *uuid.UUID    // Optional, set when received from server
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
}

// GenerateKeyPair creates a new RSA key pair
func GenerateKeyPair() (*UnsealedKeyPair, error) {
    // Generate RSA key pair
    privKey, err := crypto.GenerateRSAKeyPair(crypto.RSAKeySize2048)
    if err != nil {
        return nil, err
    }

    return &UnsealedKeyPair{
        privateKey: privKey,
        publicKey:  &privKey.PublicKey,
    }, nil
}

// InitialSeal creates an InitialKeyPair by encrypting the private key using the provided encrypter
func (kp *UnsealedKeyPair) InitialSeal(enc crypto.Encrypter) (*InitialKeyPair, error) {
    if enc == nil {
        return nil, ErrNilEncrypter
    }

    // Export public key to SPKI format
    pubKeyDer, err := crypto.ExportRSAPublicKey(kp.publicKey)
    if err != nil {
        return nil, err
    }

    // Export private key to PKCS#8
    privKeyDer, err := crypto.ExportRSAPrivateKey(kp.privateKey)
    if err != nil {
        return nil, err
    }

    // Encrypt private key with provided encrypter
    encPrivKey, err := enc.Encrypt(privKeyDer)
    if err != nil {
        return nil, err
    }

    return &InitialKeyPair{
        Algorithm:  crypto.AlgorithmRSAOAEPSHA256,
        PublicKey:  pubKeyDer,
        PrivateKey: encPrivKey,
        Created:    time.Now(),
        State:      crypto.KeyStateActive,
    }, nil
}

// Encrypt implements the Encrypter interface
func (kp *UnsealedKeyPair) Encrypt(data []byte) (*crypto.Ciphertext, error) {
    return crypto.RSAEncrypt(kp.publicKey, data)
}

// ID implements the Encrypter interface by returning a stable provider ID
func (kp *UnsealedKeyPair) ID() string {
    if kp.id == nil {
        return ""
    }
    return fmt.Sprintf("keypair-%s", kp.id.String())
}

// Decrypt performs RSA-OAEP decryption of data using the private key
func (kp *UnsealedKeyPair) Decrypt(ct *crypto.Ciphertext) ([]byte, error) {
    return crypto.RSADecrypt(kp.privateKey, ct.Data)
}

// Clear wipes the private key material from memory
func (kp *UnsealedKeyPair) Clear() {
    if kp.privateKey != nil {
        crypto.SecureWipe(kp.privateKey.D.Bytes())
        kp.privateKey = nil
    }
}

// Unseal decrypts the private key using the provided decrypter and returns an UnsealedKeyPair
func (kp *KeyPair) Unseal(dec crypto.Decrypter) (*UnsealedKeyPair, error) {
    if dec == nil {
        return nil, ErrNilDecrypter
    }

    if kp.State == crypto.KeyStateDestroyed {
        return nil, ErrKeyRevoked
    }

    // Decrypt private key
    privKeyDer, err := dec.Decrypt(kp.PrivateKey)
    if err != nil {
        return nil, err
    }

    // Parse private key from PKCS#8
    privKey, err := crypto.ImportRSAPrivateKey(privKeyDer)
    if err != nil {
        return nil, err
    }

    // Parse UUID from string
    id, err := uuid.Parse(kp.ID)
    if err != nil {
        return nil, fmt.Errorf("invalid key pair ID: %w", err)
    }

    return &UnsealedKeyPair{
        id:         &id,
        privateKey: privKey,
        publicKey:  &privKey.PublicKey,
    }, nil
}
