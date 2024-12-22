package opene_test

import (
	"net/http"
	"testing"

	"github.com/neox5/openk/internal/opene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidationError(t *testing.T) {
	err := opene.NewValidationError("test", "validate", "invalid input")
	require.NotNil(t, err)

	assert.Equal(t, "invalid input", err.Message)
	assert.Equal(t, opene.CodeValidation, err.Code)
	assert.Equal(t, "test", err.Domain)
	assert.Equal(t, "validate", err.Operation)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.NotNil(t, err.Meta)
	assert.False(t, err.IsSensitive)
}

func TestNewNotFoundError(t *testing.T) {
	err := opene.NewNotFoundError("test", "fetch", "resource not found")
	require.NotNil(t, err)

	assert.Equal(t, "resource not found", err.Message)
	assert.Equal(t, opene.CodeNotFound, err.Code)
	assert.Equal(t, "test", err.Domain)
	assert.Equal(t, "fetch", err.Operation)
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.NotNil(t, err.Meta)
	assert.False(t, err.IsSensitive)
}

func TestNewConflictError(t *testing.T) {
	err := opene.NewConflictError("test", "create", "resource exists")
	require.NotNil(t, err)

	assert.Equal(t, "resource exists", err.Message)
	assert.Equal(t, opene.CodeConflict, err.Code)
	assert.Equal(t, "test", err.Domain)
	assert.Equal(t, "create", err.Operation)
	assert.Equal(t, http.StatusConflict, err.StatusCode)
	assert.NotNil(t, err.Meta)
	assert.False(t, err.IsSensitive)
}

func TestNewInternalError(t *testing.T) {
	err := opene.NewInternalError("test", "process", "internal error")
	require.NotNil(t, err)

	assert.Equal(t, "internal error", err.Message)
	assert.Equal(t, opene.CodeInternal, err.Code)
	assert.Equal(t, "test", err.Domain)
	assert.Equal(t, "process", err.Operation)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	assert.NotNil(t, err.Meta)
	assert.True(t, err.IsSensitive)
}
