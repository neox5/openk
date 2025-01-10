package userstore

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/neox5/openk/internal/crypto"
	"github.com/neox5/openk/internal/storage/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: move this to the models package
// CreateValidUserInput creates a models.UserCreate with valid test data
func CreateValidUserInput(username string) *models.UserCreate {
	// Generate a valid RSA key pair using our package
	rsaKey, err := crypto.GenerateRSAKeyPair(crypto.RSAKeySize2048)
	if err != nil {
		panic("Failed to generate test key pair: " + err.Error())
	}

	// Export public key using our package
	publicKey, err := crypto.ExportRSAPublicKey(&rsaKey.PublicKey)
	if err != nil {
		panic("Failed to export public key: " + err.Error())
	}

	// Generate AES key for key pair protection
	aesKey, err := crypto.AESGenerateKey()
	if err != nil {
		panic("Failed to generate AES key: " + err.Error())
	}

	// Export private key for encryption
	privateKeyDER, err := crypto.ExportRSAPrivateKey(rsaKey)
	if err != nil {
		panic("Failed to export private key: " + err.Error())
	}

	// Encrypt private key using AES-GCM
	encryptedKeyPair, err := crypto.AESEncrypt(aesKey, privateKeyDER)
	if err != nil {
		panic("Failed to encrypt private key: " + err.Error())
	}

	// Generate salt
	salt, err := crypto.GenerateSalt()
	if err != nil {
		panic("Failed to generate salt: " + err.Error())
	}

	return &models.UserCreate{
		Username:         username,
		Iterations:       crypto.PBKDF2IterationCount,
		Salt:             salt,
		AuthKeyHash:      bytes.Repeat([]byte{2}, AuthKeyHashLength), // Simulated hash
		PublicKey:        publicKey,
		EncryptedKeyPair: *encryptedKeyPair,
	}
}

// TODO: move this to the models package
// VerifyUserWithInput verfies a user object with a UserCreate object
func VerifyUserWithInput(t *testing.T, actual *models.User, expected *models.UserCreate) {
	t.Helper()

	require.NotEmpty(t, actual.ID, "User ID should not be empty")
	require.Equal(t, expected.Username, actual.Username, "Username mismatch")
	require.Equal(t, expected.Iterations, actual.Iterations, "Iterations mismatch")

	// Verify cryptographic parameters
	require.Equal(t, expected.Salt, actual.Salt, "Salt mismatch")
	require.Equal(t, expected.AuthKeyHash, actual.AuthKeyHash, "AuthKeyHash mismatch")

	// Verify public key is valid RSA key
	rsaPubKey, err := crypto.ImportRSAPublicKey(actual.PublicKey)
	require.NoError(t, err, "Invalid public key format")
	require.Equal(t, expected.PublicKey, actual.PublicKey, "PublicKey mismatch")
	require.Equal(t, crypto.RSAKeySize2048, rsaPubKey.Size()*8, "Unexpected RSA key size")

	// Verify encrypted key pair structure
	require.Equal(t, crypto.NonceSize, len(actual.EncryptedKeyPair.Nonce), "Invalid nonce length")
	require.Equal(t, crypto.TagSize, len(actual.EncryptedKeyPair.Tag), "Invalid tag length")
	require.True(t, len(actual.EncryptedKeyPair.Data) > 0, "Encrypted key pair data should not be empty")
	require.Equal(t, expected.EncryptedKeyPair, actual.EncryptedKeyPair, "EncryptedKeyPair mismatch")

	verifyTimestamp(t, actual.CreatedAt)
}

// TODO: move this to the models package
// VerifyUsersEqual verifies if two user objects are euqal
func VerifyUsersEqual(t *testing.T, expected, actual *models.User) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.Equal(t, expected.Iterations, actual.Iterations)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.Salt, actual.Salt)
	assert.Equal(t, expected.AuthKeyHash, actual.AuthKeyHash)
	assert.Equal(t, expected.PublicKey, actual.PublicKey)
	assert.Equal(t, expected.EncryptedKeyPair.Nonce, actual.EncryptedKeyPair.Nonce)
	assert.Equal(t, expected.EncryptedKeyPair.Data, actual.EncryptedKeyPair.Data)
	assert.Equal(t, expected.EncryptedKeyPair.Tag, actual.EncryptedKeyPair.Tag)
}

// verifyTimestamp verifies a timestamp is within an acceptable range
func verifyTimestamp(t *testing.T, timestamp time.Time) {
	t.Helper()

	now := time.Now()
	allowedDelta := time.Second

	if timestamp.Before(now.Add(-allowedDelta)) || timestamp.After(now.Add(allowedDelta)) {
		t.Errorf("Timestamp %v not within ±%v of current time %v",
			timestamp, allowedDelta, now)
	}
}

// createCanceledContext creates a context that is already canceled
func createCanceledContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

// createTimedOutContext creates a context that will time out after duration
func createTimedOutContext(t *testing.T, timeout time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), timeout)
}
