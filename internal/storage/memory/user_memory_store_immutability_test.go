package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/neox5/openk/internal/storage/memory"
	"github.com/neox5/openk/internal/storage/models"
	"github.com/neox5/openk/internal/storage/testsuite/userstore"
	"github.com/stretchr/testify/require"
)

func TestUserMemoryStore_DataImmutability(t *testing.T) {
	t.Run("create returned data is immutable", func(t *testing.T) {
		store := memory.NewUserMemoryStore()
		input := userstore.CreateValidUserInput("create_immutable")

		original, err := store.Create(context.Background(), input)
		require.NoError(t, err)

		// Save original ID before modifying
		originalID := original.ID

		// Modify returned data
		modifyAllFields(original)

		// Verify data is unchanged in store
		retrieved, err := store.GetByID(context.Background(), originalID)
		require.NoError(t, err)
		userstore.VerifyUserWithInput(t, retrieved, input)
	})

	t.Run("getbyid returned data is immutable", func(t *testing.T) {
		store := memory.NewUserMemoryStore()
		input := userstore.CreateValidUserInput("getid_immutable")

		original, err := store.Create(context.Background(), input)
		require.NoError(t, err)

		// Get and modify data
		toChange, err := store.GetByID(context.Background(), original.ID)
		require.NoError(t, err)
		modifyAllFields(toChange)

		// Verify data is unchanged in store
		retrieved, err := store.GetByID(context.Background(), original.ID)
		require.NoError(t, err)
		userstore.VerifyUsersEqual(t, original, retrieved)
	})

	t.Run("getByUsername returned data is immutable", func(t *testing.T) {
		store := memory.NewUserMemoryStore()
		input := userstore.CreateValidUserInput("getusername_immutable")

		original, err := store.Create(context.Background(), input)
		require.NoError(t, err)

		// Get and modify data
		toChange, err := store.GetByUsername(context.Background(), original.Username)
		require.NoError(t, err)
		modifyAllFields(toChange)

		// Verify data is unchanged in store
		retrieved, err := store.GetByID(context.Background(), original.ID)
		require.NoError(t, err)
		userstore.VerifyUsersEqual(t, original, retrieved)
	})
}

func modifyAllFields(user *models.User) {
	user.ID = "modified_id"
	user.Username = "modified_name"
	user.Iterations += 1000
	user.CreatedAt = time.Now().Add(24 * time.Hour)
	user.Salt[0] = 0xFF
	user.AuthKeyHash[0] = 0xFF
	user.PublicKey[0] = 0xFF
	user.EncryptedKeyPair.Nonce[0] = 0xFF
	user.EncryptedKeyPair.Data[0] = 0xFF
	user.EncryptedKeyPair.Tag[0] = 0xFF
}
