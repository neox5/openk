package kms_test

import (
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/kms"
	"github.com/stretchr/testify/assert"
)

var (
	testPassword = []byte("test-password")
	testUsername = "test-user"
)

func TestMasterKey_Derive(t *testing.T) {
	t.Run("successful derivation", func(t *testing.T) {
		v := kms.NewVault()
		mk := kms.NewMasterKey()

		err := v.Do(func(v *kms.Vault) error {
			return mk.Derive(v, testPassword, testUsername)
		})
		assert.NoError(t, err)
	})

	t.Run("empty password", func(t *testing.T) {
		v := kms.NewVault()
		mk := kms.NewMasterKey()

		err := v.Do(func(v *kms.Vault) error {
			return mk.Derive(v, []byte{}, testUsername)
		})
		assert.Error(t, err)
	})

	t.Run("empty username", func(t *testing.T) {
		v := kms.NewVault()
		mk := kms.NewMasterKey()

		err := v.Do(func(v *kms.Vault) error {
			return mk.Derive(v, testPassword, "")
		})
		assert.Error(t, err)
	})
}

func TestMasterKey_DeriveAuth(t *testing.T) {
	t.Run("successful auth derivation", func(t *testing.T) {
		v := kms.NewVault()
		mk := kms.NewMasterKey()

		var authKey []byte
		err := v.Do(func(v *kms.Vault) error {
			if err := mk.Derive(v, testPassword, testUsername); err != nil {
				return err
			}

			derived, err := mk.DeriveAuth(v)
			if err != nil {
				return err
			}
			authKey = derived
			return nil
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, authKey)
	})

	t.Run("no master key derived", func(t *testing.T) {
		v := kms.NewVault()
		mk := kms.NewMasterKey()

		err := v.Do(func(v *kms.Vault) error {
			_, err := mk.DeriveAuth(v)
			return err
		})
		assert.ErrorIs(t, err, kms.ErrNoKey)
	})
}

func TestMasterKey_Encryption(t *testing.T) {
	plaintext := []byte("secret data")

	t.Run("encrypt and decrypt", func(t *testing.T) {
		v := kms.NewVault()
		mk := kms.NewMasterKey()

		var ciphertext *crypto.Ciphertext
		err := v.Do(func(v *kms.Vault) error {
			if err := mk.Derive(v, testPassword, testUsername); err != nil {
				return err
			}

			encrypted, err := mk.Encrypt(v, plaintext)
			if err != nil {
				return err
			}
			ciphertext = encrypted
			return nil
		})
		assert.NoError(t, err)
		assert.NotNil(t, ciphertext)

		// Decrypt in new vault operation
		var decrypted []byte
		err = v.Do(func(v *kms.Vault) error {
			if err := mk.Derive(v, testPassword, testUsername); err != nil {
				return err
			}

			d, err := mk.Decrypt(v, ciphertext)
			if err != nil {
				return err
			}
			decrypted = d
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("encrypt without master key", func(t *testing.T) {
		v := kms.NewVault()
		mk := kms.NewMasterKey()

		err := v.Do(func(v *kms.Vault) error {
			_, err := mk.Encrypt(v, plaintext)
			return err
		})
		assert.ErrorIs(t, err, kms.ErrNoKey)
	})

	t.Run("decrypt without master key", func(t *testing.T) {
		v := kms.NewVault()
		mk := kms.NewMasterKey()

		err := v.Do(func(v *kms.Vault) error {
			_, err := mk.Decrypt(v, &crypto.Ciphertext{})
			return err
		})
		assert.ErrorIs(t, err, kms.ErrNoKey)
	})
}
