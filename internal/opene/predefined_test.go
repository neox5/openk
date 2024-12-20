package opene_test

import (
	"net/http"
	"testing"

	"github.com/neox5/openk/internal/opene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidationError(t *testing.T) {
	err := opene.NewValidationError("invalid input")
	require.NotNil(t, err)

	assert.Equal(t, "invalid input", err.Message)
	assert.Equal(t, opene.CodeValidation, err.Code)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.NotNil(t, err.Meta)
	assert.False(t, err.IsSensitive)
}

func TestNewNotFoundError(t *testing.T) {
	err := opene.NewNotFoundError("resource not found")
	require.NotNil(t, err)

	assert.Equal(t, "resource not found", err.Message)
	assert.Equal(t, opene.CodeNotFound, err.Code)
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.NotNil(t, err.Meta)
	assert.False(t, err.IsSensitive)
}

func TestNewConflictError(t *testing.T) {
	err := opene.NewConflictError("resource exists")
	require.NotNil(t, err)

	assert.Equal(t, "resource exists", err.Message)
	assert.Equal(t, opene.CodeConflict, err.Code)
	assert.Equal(t, http.StatusConflict, err.StatusCode)
	assert.NotNil(t, err.Meta)
	assert.False(t, err.IsSensitive)
}

func TestNewInternalError(t *testing.T) {
	err := opene.NewInternalError("internal error")
	require.NotNil(t, err)

	assert.Equal(t, "internal error", err.Message)
	assert.Equal(t, opene.CodeInternal, err.Code)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	assert.NotNil(t, err.Meta)
	assert.True(t, err.IsSensitive)
}
