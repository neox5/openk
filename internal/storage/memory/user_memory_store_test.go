package memory_test

import (
	"testing"

	"github.com/neox5/openk/internal/storage"

	"github.com/neox5/openk/internal/storage/memory"
	"github.com/neox5/openk/internal/storage/testsuite/userstore"
)

func TestUserMemoryStore(t *testing.T) {
	// Create test suite for memory implementation
	suite := &userstore.Suite{
		NewStore: func() storage.UserStore {
			return memory.NewUserMemoryStore()
		},
	}

	// Run all standard test cases
	suite.RunAll(t)
}
