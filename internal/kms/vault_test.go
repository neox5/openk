package kms_test

import (
	"testing"

	"github.com/neox5/openk/internal/kms"
	"github.com/stretchr/testify/assert"
)

var testKey = []byte("test key")

func TestVault_Store(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr error
	}{
		{
			name: "valid key",
			key:  testKey,
		},
		{
			name:    "empty key",
			key:     []byte{},
			wantErr: kms.ErrInvalidKey,
		},
		{
			name:    "nil key",
			key:     nil,
			wantErr: kms.ErrInvalidKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := kms.NewVault()
			err := v.Store(tt.key)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.False(t, v.HasKey())
			} else {
				assert.NoError(t, err)
				assert.True(t, v.HasKey())
			}
		})
	}

	// Test storing when key already present
	t.Run("already has key", func(t *testing.T) {
		v := kms.NewVault()
		assert.NoError(t, v.Store(testKey))
		assert.ErrorIs(t, v.Store(testKey), kms.ErrKeyPresent)
	})
}

func TestVault_UseKey(t *testing.T) {
	t.Run("use key successfully", func(t *testing.T) {
		v := kms.NewVault()
		assert.NoError(t, v.Store(testKey))

		var keyUsed []byte
		err := v.UseKey(func(key []byte) error {
			keyUsed = make([]byte, len(key))
			copy(keyUsed, key)
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, testKey, keyUsed)
	})

	t.Run("empty vault", func(t *testing.T) {
		v := kms.NewVault()
		err := v.UseKey(func(key []byte) error { return nil })
		assert.ErrorIs(t, err, kms.ErrNoKey)
	})

	t.Run("operation error", func(t *testing.T) {
		v := kms.NewVault()
		assert.NoError(t, v.Store(testKey))

		err := v.UseKey(func(key []byte) error {
			return assert.AnError
		})
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestVault_Cleanup(t *testing.T) {
	t.Run("cleanup with key", func(t *testing.T) {
		v := kms.NewVault()
		assert.NoError(t, v.Store(testKey))
		assert.True(t, v.HasKey())

		assert.NoError(t, v.Cleanup())
		assert.False(t, v.HasKey())

		// Can store new key
		assert.NoError(t, v.Store(testKey))
	})

	t.Run("cleanup empty vault", func(t *testing.T) {
		v := kms.NewVault()
		assert.False(t, v.HasKey())
		assert.NoError(t, v.Cleanup())
		assert.False(t, v.HasKey())
	})
}

func TestVault_Lifecycle(t *testing.T) {
	v := kms.NewVault()
	assert.False(t, v.HasKey())

	// Store and verify key
	assert.NoError(t, v.Store(testKey))
	assert.True(t, v.HasKey())

	// Use key
	err := v.UseKey(func(key []byte) error {
		assert.Equal(t, testKey, key)
		return nil
	})
	assert.NoError(t, err)

	// Cleanup
	assert.NoError(t, v.Cleanup())
	assert.False(t, v.HasKey())
}

