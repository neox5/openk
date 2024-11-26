package crypto_test

import (
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAlgorithm_String(t *testing.T) {
	tests := []struct {
		name string
		alg  crypto.Algorithm
		want string
	}{
		{
			name: "RSA OAEP SHA256",
			alg:  crypto.AlgorithmRSAOAEPSHA256,
			want: "RSA-2048-OAEP-SHA256",
		},
		{
			name: "AES GCM 256",
			alg:  crypto.AlgorithmAESGCM256,
			want: "AES-256-GCM",
		},
		{
			name: "Unknown algorithm",
			alg:  crypto.Algorithm(999),
			want: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.alg.String())
		})
	}
}
