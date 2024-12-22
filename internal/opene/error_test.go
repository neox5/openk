package opene_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/neox5/openk/internal/opene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError_Error(t *testing.T) {
	err := &opene.Error{
		Message: "test error",
		Code:    opene.CodeValidation,
	}
	assert.Equal(t, "test error", err.Error())
}

func TestError_Unwrap(t *testing.T) {
	t.Run("unwraps Error preserving type", func(t *testing.T) {
		inner := opene.NewValidationError("test", "validate", "inner error")
		outer := opene.NewInternalError("test", "process", "outer error").Wrap(inner)

		unwrapped := outer.Unwrap()
		require.NotNil(t, unwrapped)

		innerErr, ok := unwrapped.(*opene.Error)
		require.True(t, ok)
		assert.Equal(t, opene.CodeValidation, innerErr.Code)
	})

	t.Run("unwraps standard error", func(t *testing.T) {
		stdErr := errors.New("standard error")
		err := opene.NewValidationError("test", "validate", "wrapper").
			WithMetadata(opene.Metadata{"wrapped": stdErr})

		assert.Equal(t, err.Meta["wrapped"], stdErr)
	})

	t.Run("returns nil for no wrapped error", func(t *testing.T) {
		err := opene.NewValidationError("test", "validate", "no wrapped error")
		assert.Nil(t, err.Unwrap())
	})
}

func TestError_UnwrapAll(t *testing.T) {
	t.Run("gets to root cause through Error chain", func(t *testing.T) {
		rootErr := errors.New("root cause")
		// Convert standard error to openE error
		inner := opene.AsError(rootErr, "test", opene.CodeInternal)
		// Wrap with validation error
		wrapper := opene.NewValidationError("test", "validate", "inner error").Wrap(inner)
		// Wrap again with internal error
		outer := opene.NewInternalError("test", "process", "outer error").Wrap(wrapper)

		result := outer.UnwrapAll()
		assert.Equal(t, rootErr, result)
	})

	t.Run("returns nil for no wrapped error", func(t *testing.T) {
		err := opene.NewValidationError("test", "validate", "no wrapped error")
		assert.Nil(t, err.UnwrapAll())
	})
}

func TestError_WithMetadata(t *testing.T) {
	err := opene.NewValidationError("test", "validate", "test error")
	md := opene.Metadata{
		"key": "value",
		"num": 42,
	}

	result := err.WithMetadata(md)
	assert.Equal(t, md, result.Meta)
	assert.Same(t, err, result)
}

func TestError_Sensitive(t *testing.T) {
	err := opene.NewValidationError("test", "validate", "test error")
	result := err.Sensitive()

	assert.True(t, result.IsSensitive)
	assert.Same(t, err, result)
}

func TestError_Wrap(t *testing.T) {
	t.Run("wraps another error", func(t *testing.T) {
		inner := opene.NewValidationError("inner", "validate", "validation failed")
		outer := opene.NewInternalError("outer", "process", "processing failed")

		wrapped := outer.Wrap(inner)

		assert.Contains(t, wrapped.Message, "processing failed")
		assert.Contains(t, wrapped.Message, "validation failed")
		assert.Equal(t, opene.CodeInternal, wrapped.Code)
		assert.Equal(t, "outer", wrapped.Domain)
		assert.Equal(t, "process", wrapped.Operation)
		assert.Equal(t, http.StatusInternalServerError, wrapped.StatusCode)
		assert.Equal(t, inner, wrapped.WrappedErr)
	})

	t.Run("propagates sensitivity flag", func(t *testing.T) {
		inner := opene.NewInternalError("test", "process", "sensitive error")
		outer := opene.NewValidationError("test", "validate", "wrapper")

		wrapped := outer.Wrap(inner)
		assert.True(t, wrapped.IsSensitive)
	})

	t.Run("handles nil error", func(t *testing.T) {
		base := opene.NewValidationError("test", "validate", "base error")
		wrapped := base.Wrap(nil)
		assert.Same(t, base, wrapped)
	})
}

func TestAsError(t *testing.T) {
	t.Run("converts standard error", func(t *testing.T) {
		stdErr := errors.New("standard error")
		err := opene.AsError(stdErr, "test", opene.CodeValidation)
		require.NotNil(t, err)

		assert.Equal(t, "standard error", err.Message)
		assert.Equal(t, opene.CodeValidation, err.Code)
		assert.Equal(t, "test", err.Domain)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.False(t, err.IsSensitive)
		assert.Equal(t, stdErr, err.WrappedErr)
	})

	t.Run("preserves existing error", func(t *testing.T) {
		original := opene.NewValidationError("test", "validate", "original error")
		err := opene.AsError(original, "new", opene.CodeInternal)
		assert.Same(t, original, err)
	})

	t.Run("handles nil error", func(t *testing.T) {
		err := opene.AsError(nil, "test", opene.CodeValidation)
		assert.Nil(t, err)
	})
}
