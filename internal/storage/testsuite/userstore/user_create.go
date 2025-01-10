package userstore

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/neox5/openk/internal/opene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) testCreateUser(t *testing.T) {
	t.Run("creates new user with valid input", func(t *testing.T) {
		store := s.NewStore()
		input := CreateValidUserInput("testuser")
		user, err := store.Create(context.Background(), input)
		require.NoError(t, err)
		VerifyUserWithInput(t, user, input)
	})

	t.Run("handles length boundaries", func(t *testing.T) {
		t.Run("single character", func(t *testing.T) {
			store := s.NewStore()
			input := CreateValidUserInput("x")
			user, err := store.Create(context.Background(), input)
			require.NoError(t, err)
			VerifyUserWithInput(t, user, input)
		})

		t.Run("maximum length", func(t *testing.T) {
			store := s.NewStore()
			username := strings.Repeat("a", 256)
			input := CreateValidUserInput(username)
			user, err := store.Create(context.Background(), input)
			require.NoError(t, err)
			VerifyUserWithInput(t, user, input)
		})

		t.Run("exceeds maximum length", func(t *testing.T) {
			store := s.NewStore()
			username := strings.Repeat("a", 257)
			input := CreateValidUserInput(username)
			user, err := store.Create(context.Background(), input)
			assert.Error(t, err)
			assert.Nil(t, user)

			var e *opene.Error
			require.ErrorAs(t, err, &e)
			assert.Equal(t, opene.CodeValidation, e.Code)
			assert.Contains(t, e.Message, "too long")
		})
	})

	t.Run("rejects creation with nil input", func(t *testing.T) {
		store := s.NewStore()
		user, err := store.Create(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, user)

		var e *opene.Error
		require.ErrorAs(t, err, &e)
		assert.Equal(t, opene.CodeValidation, e.Code)
		assert.Contains(t, e.Message, "nil")
	})

	t.Run("rejects creation with empty username", func(t *testing.T) {
		store := s.NewStore()
		input := CreateValidUserInput("")
		user, err := store.Create(context.Background(), input)
		assert.Error(t, err)
		assert.Nil(t, user)

		var e *opene.Error
		require.ErrorAs(t, err, &e)
		assert.Equal(t, opene.CodeValidation, e.Code)
		assert.Contains(t, e.Message, "empty")
	})

	t.Run("rejects creation with nil byte arrays", func(t *testing.T) {
		store := s.NewStore()
		input := CreateValidUserInput("nilbytes")
		input.Salt = nil
		input.AuthKeyHash = nil
		input.PublicKey = nil
		input.EncryptedKeyPair.Nonce = nil
		input.EncryptedKeyPair.Data = nil
		input.EncryptedKeyPair.Tag = nil

		user, err := store.Create(context.Background(), input)
		assert.Error(t, err)
		assert.Nil(t, user)

		var e *opene.Error
		require.ErrorAs(t, err, &e)
		assert.Equal(t, opene.CodeValidation, e.Code)
		assert.Contains(t, e.Message, "required")
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		store := s.NewStore()
		ctx := createCanceledContext(t)
		input := CreateValidUserInput("canceled")

		user, err := store.Create(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("handles context timeout", func(t *testing.T) {
		store := s.NewStore()
		ctx, cancel := createTimedOutContext(t, time.Microsecond)
		defer cancel()
		time.Sleep(time.Millisecond)

		input := CreateValidUserInput("timeout")
		user, err := store.Create(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
