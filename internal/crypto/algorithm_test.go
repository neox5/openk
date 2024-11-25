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

func TestAlgorithm_Valid(t *testing.T) {
	tests := []struct {
		name string
		alg  crypto.Algorithm
		want bool
	}{
		{
			name: "RSA OAEP SHA256 is valid",
			alg:  crypto.AlgorithmRSAOAEPSHA256,
			want: true,
		},
		{
			name: "AES GCM 256 is valid",
			alg:  crypto.AlgorithmAESGCM256,
			want: true,
		},
		{
			name: "Unknown algorithm is invalid",
			alg:  crypto.Algorithm(999),
			want: false,
		},
		{
			name: "Negative algorithm is invalid",
			alg:  crypto.Algorithm(-1),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.alg.Valid())
		})
	}
}
