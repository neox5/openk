package crypto_test

import (
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func TestKeyState_String(t *testing.T) {
	tests := []struct {
		name  string
		state crypto.KeyState
		want  string
	}{
		{
			name:  "Active state",
			state: crypto.KeyStateActive,
			want:  "Active",
		},
		{
			name:  "PendingRotation state",
			state: crypto.KeyStatePendingRotation,
			want:  "PendingRotation",
		},
		{
			name:  "Inactive state",
			state: crypto.KeyStateInactive,
			want:  "Inactive",
		},
		{
			name:  "Destroyed state",
			state: crypto.KeyStateDestroyed,
			want:  "Destroyed",
		},
		{
			name:  "Unknown state",
			state: crypto.KeyState(999),
			want:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.state.String())
		})
	}
}
