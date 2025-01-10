package userstore

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/neox5/openk/internal/opene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) testConcurrent(t *testing.T) {
	t.Run("duplicate creation race", func(t *testing.T) {
		store := s.NewStore()
		const numGoroutines = 10
		var wg sync.WaitGroup
		var successCount atomic.Int32

		input := CreateValidUserInput("concurrent_duplicate")
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := store.Create(context.Background(), input)
				if err == nil {
					successCount.Add(1)
				} else {
					var e *opene.Error
					require.ErrorAs(t, err, &e)
					assert.Equal(t, opene.CodeConflict, e.Code)
				}
			}()
		}
		wg.Wait()

		assert.Equal(t, int32(1), successCount.Load(), "exactly one creation should succeed")

		// Verify the user was created correctly
		user, err := store.GetByUsername(context.Background(), input.Username)
		require.NoError(t, err)
		VerifyUserWithInput(t, user, input)
	})

	t.Run("mixed read/write operations", func(t *testing.T) {
		store := s.NewStore()
		var wg sync.WaitGroup
		done := make(chan struct{})

		// Track operations
		var reads atomic.Int32
		var writes atomic.Int32

		// Create initial test user
		initialUser := CreateValidUserInput("mixed_test_0")
		created, err := store.Create(context.Background(), initialUser)
		require.NoError(t, err)

		// Start readers (3)
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-done:
						return
					default:
						_, err := store.GetByUsername(context.Background(), created.Username)
						if err == nil {
							reads.Add(1)
						}
					}
				}
			}()
		}

		// Start writers (2)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(writerID int) {
				defer wg.Done()
				counter := 0
				for {
					select {
					case <-done:
						return
					default:
						input := CreateValidUserInput(fmt.Sprintf("mixed_test_%d", counter%5))
						_, err := store.Create(context.Background(), input)
						if err == nil {
							writes.Add(1)
						}
						counter++
					}
				}
			}(i)
		}

		// Let operations run briefly
		time.Sleep(time.Second)
		close(done)
		wg.Wait()

		// Verify operations occurred
		assert.Greater(t, reads.Load(), int32(0), "should have successful reads")
		assert.Greater(t, writes.Load(), int32(0), "should have successful writes")
		t.Logf("Mixed operations completed - Reads: %d, Writes: %d", reads.Load(), writes.Load())
	})
}
