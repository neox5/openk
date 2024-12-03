package crypto_test

import (
	"bytes"
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func TestMemory_SecureWipe(t *testing.T) {
	t.Run("basic buffer wiping", func(t *testing.T) {
		tests := []struct {
			name string
			size int
			fill byte
		}{
			{
				name: "wipes empty buffer",
				size: 0,
				fill: 0x00,
			},
			{
				name: "wipes small buffer",
				size: 16,
				fill: 0xFF,
			},
			{
				name: "wipes medium buffer",
				size: 1024,
				fill: 0xAA,
			},
			{
				name: "wipes large buffer",
				size: 1024 * 1024, // 1MB
				fill: 0x55,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Create and fill buffer with known pattern
				buf := bytes.Repeat([]byte{tt.fill}, tt.size)

				// Perform secure wipe
				crypto.SecureWipe(buf)

				// Verify buffer is zeroed
				expected := make([]byte, tt.size)
				assert.True(t, bytes.Equal(buf, expected), "buffer should be zeroed after wipe")
			})
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Run("handles nil buffer", func(t *testing.T) {
			// Should not panic
			crypto.SecureWipe(nil)
		})

		t.Run("handles slice with zero capacity", func(t *testing.T) {
			buf := make([]byte, 0)
			crypto.SecureWipe(buf)
			assert.Equal(t, 0, len(buf), "length should remain zero")
		})

		t.Run("handles slice with capacity > length", func(t *testing.T) {
			// Create slice with extra capacity
			buf := make([]byte, 10, 20)
			for i := range buf {
				buf[i] = 0xFF
			}

			// Wipe the buffer
			crypto.SecureWipe(buf)

			// Verify visible portion is zeroed
			expected := make([]byte, 10)
			assert.True(t, bytes.Equal(buf, expected), "visible portion should be zeroed")

			// Verify length hasn't changed
			assert.Equal(t, 10, len(buf), "length should remain unchanged")
		})
	})

	t.Run("pattern tests", func(t *testing.T) {
		tests := []struct {
			name    string
			pattern []byte
		}{
			{
				name:    "wipes alternating bits",
				pattern: []byte{0xAA, 0x55, 0xAA, 0x55},
			},
			{
				name:    "wipes all ones",
				pattern: []byte{0xFF, 0xFF, 0xFF, 0xFF},
			},
			{
				name:    "wipes all zeros",
				pattern: []byte{0x00, 0x00, 0x00, 0x00},
			},
			{
				name:    "wipes mixed pattern",
				pattern: []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Create buffer with repeated pattern
				buf := make([]byte, 1024)
				for i := 0; i < len(buf); i += len(tt.pattern) {
					copy(buf[i:], tt.pattern)
				}

				// Perform secure wipe
				crypto.SecureWipe(buf)

				// Verify buffer is zeroed
				expected := make([]byte, len(buf))
				assert.True(t, bytes.Equal(buf, expected), "buffer should be zeroed after wipe")
			})
		}
	})
}
