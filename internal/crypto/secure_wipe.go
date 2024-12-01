package crypto

import "crypto/subtle"

// SecureWipe overwrites the buffer with zeros using constant time operations.
// This helps ensure sensitive data is properly cleared from memory.
func SecureWipe(buf []byte) {
	if len(buf) > 0 {
		zeros := make([]byte, len(buf))
		subtle.ConstantTimeCopy(1, buf, zeros)
	}
}
