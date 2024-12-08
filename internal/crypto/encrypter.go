package crypto

// Encrypter represents any entity capable of encrypting data.
// This interface decouples key protection from specific encryption implementations.
type Encrypter interface {
	// Encrypt encrypts the provided data and returns a Ciphertext
	Encrypt(data []byte) (*Ciphertext, error)

	// ID returns a unique identifier for this encryption provider.
	// This ID is used to track which provider encrypted data in envelopes
	// and must remain stable across program restarts.
	// Format requirements:
	// - Must be non-empty
	// - Must be URL-safe (only alphanumeric, '-', '_')
	// - Maximum length of 63 characters
	ID() string
}
