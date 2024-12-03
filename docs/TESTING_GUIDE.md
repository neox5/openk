# OpenK Testing Guide

## File Structure
- Name test files: `foo_test.go` for `foo.go`
- Use separate test package: `package foo_test`
- Group imports with stdlib first

## Naming Conventions

### Test Functions
```go
// Standard pattern 
TestType_Method(t *testing.T)

// Examples
TestRSA_Generate(t *testing.T)
TestMemory_SecureWipe(t *testing.T)
TestPBKDF2_DeriveKey(t *testing.T)
```

### Subtests
Format: Descriptive phrase in lowercase, grouped by success/error cases
```go
t.Run("success cases", func(t *testing.T) {
    t.Run("valid key is exported", func(t *testing.T) {...})
    t.Run("empty message is handled", func(t *testing.T) {...})
})

t.Run("error cases", func(t *testing.T) {
    t.Run("rejects nil key", func(t *testing.T) {...})
    t.Run("handles oversized input", func(t *testing.T) {...})
})
```

## Test Structure

### Basic Pattern
```go
func TestType_Method(t *testing.T) {
    // Optional common setup
    key := generateTestKey(t)

    t.Run("success cases", func(t *testing.T) {
        // Test happy paths
    })

    t.Run("error cases", func(t *testing.T) {
        // Test error conditions
    })
}
```

### Table-Driven Tests
Use only when testing multiple similar input/output variations:

```go
func TestRSA_Generate(t *testing.T) {
    tests := []struct {
        name    string
        bits    int
    }{
        {
            name: "generates 2048-bit key",
            bits: RSAKeySize2048,
        },
        {
            name: "generates 4096-bit key",
            bits: RSAKeySize4096,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            key, err := GenerateRSAKeyPair(tt.bits)
            assert.NoError(t, err)
            assert.Equal(t, tt.bits, key.Size()*8)
        })
    }
}
```

### Complex Test Cases
For tests with complex setup or multiple steps, prefer direct testing over tables:

```go
func TestRSA_ImportPrivateKey(t *testing.T) {
    t.Run("success cases", func(t *testing.T) {
        key := generateTestKey(t, RSAKeySize2048)
        der, err := MarshalPKCS8PrivateKey(key)
        require.NoError(t, err)

        imported, err := ImportRSAPrivateKey(der)
        assert.NoError(t, err)
        assert.Equal(t, key.D.Bytes(), imported.D.Bytes())
    })

    t.Run("error cases", func(t *testing.T) {
        t.Run("rejects invalid DER", func(t *testing.T) {
            imported, err := ImportRSAPrivateKey([]byte("invalid"))
            assert.ErrorIs(t, err, ErrInvalidPrivateKey)
            assert.Contains(t, err.Error(), "invalid format")
            assert.Nil(t, imported)
        })

        t.Run("rejects EC key", func(t *testing.T) {
            ecKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
            der, _ := MarshalPKCS8PrivateKey(ecKey)

            imported, err := ImportRSAPrivateKey(der)
            assert.ErrorIs(t, err, ErrInvalidPrivateKey)
            assert.Contains(t, err.Error(), "not an RSA key")
            assert.Nil(t, imported)
        })
    })
}
```

### Error Checking
- Use require for critical setup that must succeed
- Use assert for test validations
- Check both error type and message when relevant

```go
// Setup with require
key := generateTestKey(t)
der, err := MarshalPrivateKey(key)
require.NoError(t, err)

// Test validations with assert
result, err := ImportKey(der)
assert.ErrorIs(t, err, ErrInvalidKey)
assert.Contains(t, err.Error(), "invalid format")
assert.Nil(t, result)
```

## Best Practices (in order of priority)

1. Keep tests focused and small
   - Each test should validate one piece of functionality
   - Split complex tests into smaller, focused subtests
   - Keep setup close to where it's used

2. Group logically with t.Run()
   - Separate success and error cases
   - Use descriptive names for test groups
   - Keep nesting to maximum of 2 levels

3. Use table-driven tests appropriately
   - Only for two or more similar input/output variations
   - Avoid for complex setup or multi-step tests
   - Keep test cases clear and readable

4. Use require vs assert correctly
   - require: For critical setup that must succeed
   - assert: For actual test validations
   - Fail fast when setup fails

5. Write clear error checks
   - Test both error type and message when relevant
   - Use ErrorIs() for error types
   - Use Contains() for error messages

6. Add useful test helpers
   - Put common setup in helper functions
   - Mark helpers with t.Helper()
   - Keep helpers focused and simple

## Running Tests

```bash
go test ./...                     # All tests
go test -run TestType_Method ./...   # Single test
go test -v ./...                  # Verbose output
```
