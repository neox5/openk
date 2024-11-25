package crypto

// Algorithm represents supported cryptographic algorithms
type Algorithm int

const (
	// AlgorithmRSAOAEPSHA256 represents RSA-2048 with OAEP padding and SHA-256
	// Used for key wrapping operations
	AlgorithmRSAOAEPSHA256 Algorithm = iota

	// AlgorithmAESGCM256 represents AES-256 in GCM mode
	// Used for data encryption operations
	AlgorithmAESGCM256
)

// String implements the Stringer interface for Algorithm
func (a Algorithm) String() string {
	switch a {
	case AlgorithmRSAOAEPSHA256:
		return "RSA-2048-OAEP-SHA256"
	case AlgorithmAESGCM256:
		return "AES-256-GCM"
	default:
		return "UNKNOWN"
	}
}

// Valid returns true if the Algorithm is a known value
func (a Algorithm) Valid() bool {
	switch a {
	case AlgorithmRSAOAEPSHA256, AlgorithmAESGCM256:
		return true
	default:
		return false
	}
}
