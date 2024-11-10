# Encryption Service

```go
type EncryptionService interface {
    Encrypt(plaintext string) (string, error)
    Decrypt(ciphertext string) (string, error)
}
```
