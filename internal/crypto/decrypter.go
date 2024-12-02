package crypto

// Decrypter represents any entity capable of decrypting data.
// This interface decouples key protection from specific decryption implementations.
type Decrypter interface {
	// Decrypt decrypts the provided ciphertext and returns the plaintext
	Decrypt(ct *Ciphertext) ([]byte, error)
}
