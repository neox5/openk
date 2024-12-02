package crypto

// Encrypter represents any entity capable of encrypting data.
// This interface decouples key protection from specific encryption implementations.
type Encrypter interface {
    // Encrypt encrypts the provided data and returns a Ciphertext
    Encrypt(data []byte) (*Ciphertext, error)
}
