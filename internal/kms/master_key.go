package kms

import "github.com/neox5/openk/internal/crypto"

// MasterKey handles master key operations
type MasterKey struct {}

func NewMasterKey() *MasterKey {
    return &MasterKey{}
}

func (mk *MasterKey) Derive(v *Vault, password []byte, username string) error {
    key, err := crypto.DeriveMasterKey(password, []byte(username))
    if err != nil {
        return err
    }
    return v.Store(key)
}

func (mk *MasterKey) DeriveAuth(v *Vault) ([]byte, error) {
    var authKey []byte
    err := v.UseKey(func(key []byte) error {
        derived, err := crypto.DeriveKey(key, []byte("auth4openk"), crypto.IterationCount, crypto.MasterKeySize)
        if err != nil {
            return err
        }
        authKey = derived
        return nil
    })
    return authKey, err
}

func (mk *MasterKey) Encrypt(v *Vault, plaintext []byte) (*crypto.Ciphertext, error) {
    var ciphertext *crypto.Ciphertext
    err := v.UseKey(func(key []byte) error {
        encrypted, err := crypto.AESEncrypt(key, plaintext)
        if err != nil {
            return err
        }
        ciphertext = encrypted
        return nil
    })
    return ciphertext, err
}

func (mk *MasterKey) Decrypt(v *Vault, ciphertext *crypto.Ciphertext) ([]byte, error) {
    var plaintext []byte
    err := v.UseKey(func(key []byte) error {
        decrypted, err := crypto.AESDecrypt(key, ciphertext)
        if err != nil {
            return err
        }
        plaintext = decrypted
        return nil
    })
    return plaintext, err
}
