package crypto_test

import (
	"testing"

	"github.com/neox5/openk/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSecureWipe(t *testing.T) {
	buf := []byte{1, 2, 3, 4, 5}
	expected := make([]byte, len(buf))

	crypto.SecureWipe(buf)
	assert.Equal(t, expected, buf, "Buffer should be wiped to all zeros")
}

