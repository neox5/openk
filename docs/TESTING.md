# Testing Cheat Sheet

## Example
```go
package crypto_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/neox5/openk/internal/crypto"
)

// Single method test
func TestAlgorithm_String(t *testing.T) {
    assert.Equal(t, "RSA-2048-OAEP-SHA256", crypto.AlgorithmRSAOAEPSHA256.String())
}

// Multiple test cases
func TestAlgorithm_Valid_Multiple(t *testing.T) {
    tests := []struct {
        name string
        alg  crypto.Algorithm
        want bool
    }{
        {"RSA is valid", crypto.AlgorithmRSAOAEPSHA256, true},
        {"Invalid algorithm", crypto.Algorithm(999), false},
    }

    for _, tt := range tests {
        assert.Equal(t, tt.want, tt.alg.Valid())
    }
}
```

## Naming
- Test files: `foo_test.go` for `foo.go`
- Test package: `package foo_test`
- Test functions: `TestType_Method` or `TestType_Method_Scenario`
- Benchmark functions: `BenchmarkType_Method`

## Testify Assertions
```go
assert.Equal(t, expected, actual)
assert.True(t, value)
assert.False(t, value)
assert.Error(t, err)
assert.NoError(t, err)
assert.ErrorIs(t, err, expectedErr)
assert.Nil(t, value)
assert.NotNil(t, value)
```

## Running Tests
```bash
# All tests in current project
go test ./...

# Specific package
go test ./internal/crypto/...

# Single test
go test -run TestMyFunction ./...
```

## Coverage
```bash
go test -cover ./...
go test -coverprofile=cover.out
go tool cover -html=cover.out
```