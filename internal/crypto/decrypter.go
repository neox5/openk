package crypto

// Decrypter represents any entity capable of decrypting data.
// This interface decouples key protection from specific decryption implementations.
type Decrypter interface {
	// Decrypt decrypts the provided ciphertext and returns the plaintext
	Decrypt(ct *Ciphertext) ([]byte, error)

	// ID returns a unique identifier for this encryption provider.
	// This ID is used to track which provider encrypted data in envelopes
	// and must remain stable across program restarts.
	// Format requirements:
	// - Must be non-empty
	// - Must be URL-safe (only alphanumeric, '-', '_')
	// - Maximum length of 63 characters
	ID() string
}
