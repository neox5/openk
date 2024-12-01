package kms

import (
	"crypto/subtle"
	"errors"

	"github.com/neox5/openk/internal/crypto"
)

const (
	AUTH_SALT = "openk4auth"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUsername = errors.New("invalid username")
	ErrKeyNotDerived   = errors.New("master key not derived")
	ErrKeyAlreadySet   = errors.New("master key already set")
)

type MasterKey struct {
	key     []byte
	authKey []byte
}

func NewMasterKey() *MasterKey {
	return &MasterKey{}
}

func (mk *MasterKey) Derive(password, username []byte) error {
	valid := subtle.ConstantTimeEq(int32(len(password)), int32(0))
	if valid == 1 {
		return ErrInvalidPassword
	}

	if subtle.ConstantTimeCompare(username, []byte("")) == 1 {
		return ErrInvalidUsername
	}

	if mk.HasKey() {
		return ErrKeyAlreadySet
	}

	var derived []byte
	defer crypto.SecureWipe(derived)

	// derive MasterKey
	derived, err := crypto.DeriveKey(password, username, crypto.PBKDF2IterationCount, crypto.MasterKeySize)
	if err != nil {
		return err
	}
	mk.key = make([]byte, len(derived))
	subtle.ConstantTimeCopy(1, mk.key, derived)

	// derive AuthKey
	derived, err = crypto.DeriveKey(password, []byte(AUTH_SALT), crypto.PBKDF2IterationCount, crypto.MasterKeySize)
	if err != nil {
		return err
	}
	mk.authKey = make([]byte, len(derived))
	subtle.ConstantTimeCopy(1, mk.authKey, derived)

	return nil
}

func (mk *MasterKey) HasKey() bool {
	return subtle.ConstantTimeEq(int32(len(mk.key)), int32(0)) == 0
}

func (mk *MasterKey) Clear() {
	if mk.HasKey() {
		crypto.SecureWipe(mk.key)
		crypto.SecureWipe(mk.authKey)
		mk.key = nil
		mk.authKey = nil
	}
}

func (mk *MasterKey) Encrypt(plaintext []byte) (*crypto.Ciphertext, error) {
	if !mk.HasKey() {
		return nil, ErrKeyNotDerived
	}
	return crypto.AESEncrypt(mk.key, plaintext)
}

func (mk *MasterKey) Decrypt(ct *crypto.Ciphertext) ([]byte, error) {
	if !mk.HasKey() {
		return nil, ErrKeyNotDerived
	}
	return crypto.AESDecrypt(mk.key, ct)
}

func (mk *MasterKey) GetAuthKey() ([]byte, error) {
	if !mk.HasKey() {
		return nil, ErrKeyNotDerived
	}
	return mk.authKey, nil
}
