package crypto_test

import (
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func TestNewCiphertext(t *testing.T) {
	tests := []struct {
		name    string
		nonce   []byte
		data    []byte
		tag     []byte
		want    *crypto.Ciphertext
		wantErr error
	}{
		{
			name:    "valid ciphertext",
			nonce:   make([]byte, crypto.NonceSize),
			data:    []byte("test data"),
			tag:     make([]byte, crypto.TagSize),
			want:    &crypto.Ciphertext{Nonce: make([]byte, crypto.NonceSize), Data: []byte("test data"), Tag: make([]byte, crypto.TagSize)},
			wantErr: nil,
		},
		{
			name:    "nil nonce",
			nonce:   nil,
			data:    []byte("test data"),
			tag:     make([]byte, crypto.TagSize),
			want:    nil,
			wantErr: crypto.ErrNilParameter,
		},
		{
			name:    "nil data",
			nonce:   make([]byte, crypto.NonceSize),
			data:    nil,
			tag:     make([]byte, crypto.TagSize),
			want:    nil,
			wantErr: crypto.ErrNilParameter,
		},
		{
			name:    "nil tag",
			nonce:   make([]byte, crypto.NonceSize),
			data:    []byte("test data"),
			tag:     nil,
			want:    nil,
			wantErr: crypto.ErrNilParameter,
		},
		{
			name:    "invalid nonce size",
			nonce:   make([]byte, crypto.NonceSize-1),
			data:    []byte("test data"),
			tag:     make([]byte, crypto.TagSize),
			want:    nil,
			wantErr: crypto.ErrInvalidNonce,
		},
		{
			name:    "invalid tag size",
			nonce:   make([]byte, crypto.NonceSize),
			data:    []byte("test data"),
			tag:     make([]byte, crypto.TagSize-1),
			want:    nil,
			wantErr: crypto.ErrInvalidTag,
		},
		{
			name:    "empty data is valid",
			nonce:   make([]byte, crypto.NonceSize),
			data:    []byte{},
			tag:     make([]byte, crypto.TagSize),
			want:    &crypto.Ciphertext{Nonce: make([]byte, crypto.NonceSize), Data: []byte{}, Tag: make([]byte, crypto.TagSize)},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := crypto.NewCiphertext(tt.nonce, tt.data, tt.tag)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

