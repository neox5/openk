# OpenK Code Style Guide

## Error Handling

### Error Organization
- Define errors in `errors.go` within each package
- Keep error definitions close to their usage
- Group related errors together

### Error Definition Pattern
```go
var (
    // Username errors - subject after Err
    ErrUsernameEmpty    = errors.New("username cannot be empty")
    ErrUsernameLength   = errors.New("username exceeds maximum length")
    ErrUsernameInvalid  = errors.New("username contains invalid characters")

    // Key errors - subject after Err
    ErrKeyRevoked       = errors.New("key has been revoked")
    ErrKeyNotDerived    = errors.New("key not derived")
    ErrKeyAlreadySet    = errors.New("key already set")

    // Envelope errors - subject after Err
    ErrEnvelopeMissing  = errors.New("envelope not found")
    ErrEnvelopeInvalid  = errors.New("envelope contains invalid data")

    // Don't:
    ErrNoKey           = errors.New("key not found")          // Use ErrKeyMissing instead
    ErrInvalidUsername = errors.New("username is invalid")    // Use ErrUsernameInvalid instead
    ErrorKeyNotFound   = errors.New("Key not found.")         // Don't prefix with Error or use punctuation
    errKeyInvalid     = errors.New("Invalid key")            // Don't use lowercase for exported errors
)
```

### Error Wrapping
- Wrap errors when adding context
- Keep the error chain focused
```go
// Do:
return fmt.Errorf("failed to decrypt key: %w", err)
return fmt.Errorf("invalid key format: %w", err)

// Don't:
return fmt.Errorf("error occurred while trying to perform decryption: %w", err)  // Too verbose
return fmt.Errorf("Error: %w", err)  // Too vague
return err  // Missing context
```
