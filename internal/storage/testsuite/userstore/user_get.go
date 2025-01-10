package userstore

import (
	"context"
	"testing"
	"time"

	"github.com/neox5/openk/internal/opene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) testGetByID(t *testing.T) {
	t.Run("success - retrieves existing user", func(t *testing.T) {
		store := s.NewStore()
		input := CreateValidUserInput("getid_test")

		original, err := store.Create(context.Background(), input)
		require.NoError(t, err)

		retrieved, err := store.GetByID(context.Background(), original.ID)
		require.NoError(t, err)
		VerifyUsersEqual(t, original, retrieved)
	})

	t.Run("error - not found", func(t *testing.T) {
		store := s.NewStore()
		_, err := store.GetByID(context.Background(), "nonexistent")

		var e *opene.Error
		require.ErrorAs(t, err, &e)
		assert.Equal(t, opene.CodeNotFound, e.Code)
	})

	t.Run("error - context canceled", func(t *testing.T) {
		store := s.NewStore()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := store.GetByID(ctx, "any_id")
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("error - context timeout", func(t *testing.T) {
		store := s.NewStore()
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()
		time.Sleep(time.Millisecond)

		_, err := store.GetByID(ctx, "any_id")
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}

func (s *Suite) testGetByUsername(t *testing.T) {
	t.Run("success - retrieves existing user", func(t *testing.T) {
		store := s.NewStore()
		input := CreateValidUserInput("getname_test")

		original, err := store.Create(context.Background(), input)
		require.NoError(t, err)

		retrieved, err := store.GetByUsername(context.Background(), original.Username)
		require.NoError(t, err)
		VerifyUsersEqual(t, original, retrieved)
	})

	t.Run("error - not found", func(t *testing.T) {
		store := s.NewStore()
		_, err := store.GetByUsername(context.Background(), "nonexistent")

		var e *opene.Error
		require.ErrorAs(t, err, &e)
		assert.Equal(t, opene.CodeNotFound, e.Code)
	})

	t.Run("error - context canceled", func(t *testing.T) {
		store := s.NewStore()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := store.GetByUsername(ctx, "any_name")
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("error - context timeout", func(t *testing.T) {
		store := s.NewStore()
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()
		time.Sleep(time.Millisecond)

		_, err := store.GetByUsername(ctx, "any_name")
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
