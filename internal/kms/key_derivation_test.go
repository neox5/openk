package kms_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/neox5/openk/internal/kms"
)

func TestKeyDerivation_New(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("creates with minimum values", func(t *testing.T) {
			params, err := kms.NewKeyDerivation("testuser", kms.MinIterations)
			require.NoError(t, err)
			assert.Equal(t, "testuser", params.Username)
			assert.Equal(t, kms.MinIterations, params.Iterations)
		})

		t.Run("creates with higher iterations", func(t *testing.T) {
			params, err := kms.NewKeyDerivation("testuser", kms.MinIterations*2)
			require.NoError(t, err)
			assert.Equal(t, kms.MinIterations*2, params.Iterations)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("rejects empty username", func(t *testing.T) {
			params, err := kms.NewKeyDerivation("", kms.MinIterations)
			assert.ErrorIs(t, err, kms.ErrUsernameEmpty)
			assert.Nil(t, params)
		})

		t.Run("rejects long username", func(t *testing.T) {
			username := make([]byte, kms.MaxUsernameLen+1)
			for i := range username {
				username[i] = 'a'
			}
			params, err := kms.NewKeyDerivation(string(username), kms.MinIterations)
			assert.ErrorIs(t, err, kms.ErrUsernameLength)
			assert.Nil(t, params)
		})

		t.Run("rejects non-printable characters", func(t *testing.T) {
			params, err := kms.NewKeyDerivation("test\x00user", kms.MinIterations)
			assert.ErrorIs(t, err, kms.ErrUsernameInvalid)
			assert.Nil(t, params)
		})

		t.Run("rejects non-ascii characters", func(t *testing.T) {
			params, err := kms.NewKeyDerivation("test£user", kms.MinIterations)
			assert.ErrorIs(t, err, kms.ErrUsernameInvalid)
			assert.Nil(t, params)
		})

		t.Run("rejects low iteration count", func(t *testing.T) {
			params, err := kms.NewKeyDerivation("testuser", kms.MinIterations-1)
			assert.ErrorIs(t, err, kms.ErrIterationsInvalid)
			assert.Nil(t, params)
		})
	})
}
