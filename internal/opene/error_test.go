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
		inner := &opene.Error{
			Message: "inner error",
			Code:    opene.CodeValidation,
		}
		outer := &opene.Error{
			Message:    "outer error",
			Code:       opene.CodeInternal,
			WrappedErr: inner,
		}

		unwrapped := outer.Unwrap()
		require.NotNil(t, unwrapped)

		innerErr, ok := unwrapped.(*opene.Error)
		require.True(t, ok)
		assert.Equal(t, opene.CodeValidation, innerErr.Code)
	})

	t.Run("unwraps standard error", func(t *testing.T) {
		stdErr := errors.New("standard error")
		err := &opene.Error{
			Message:    "wrapper",
			Code:       opene.CodeValidation,
			WrappedErr: stdErr,
		}

		unwrapped := err.Unwrap()
		assert.Equal(t, stdErr, unwrapped)
	})

	t.Run("returns nil for no wrapped error", func(t *testing.T) {
		err := &opene.Error{
			Message: "no wrapped error",
			Code:    opene.CodeValidation,
		}
		assert.Nil(t, err.Unwrap())
	})
}

func TestError_UnwrapAll(t *testing.T) {
	t.Run("gets to root cause through Error chain", func(t *testing.T) {
		rootErr := errors.New("root cause")

		inner := &opene.Error{
			Message:    "inner error",
			Code:       opene.CodeValidation,
			WrappedErr: rootErr,
		}

		outer := &opene.Error{
			Message:    "outer error",
			Code:       opene.CodeInternal,
			WrappedErr: inner,
		}

		result := outer.UnwrapAll()
		assert.Equal(t, rootErr, result)
	})

	t.Run("returns nil for no wrapped error", func(t *testing.T) {
		err := &opene.Error{
			Message: "no wrapped error",
			Code:    opene.CodeValidation,
		}
		assert.Nil(t, err.UnwrapAll())
	})
}

func TestError_WithDomain(t *testing.T) {
	err := &opene.Error{Message: "test error"}
	result := err.WithDomain("crypto")

	assert.Equal(t, "crypto", result.Domain)
	assert.Same(t, err, result)
}

func TestError_WithOperation(t *testing.T) {
	err := &opene.Error{Message: "test error"}
	result := err.WithOperation("key_rotation")

	assert.Equal(t, "key_rotation", result.Operation)
	assert.Same(t, err, result)
}

func TestError_WithMetadata(t *testing.T) {
	err := &opene.Error{Message: "test error"}
	md := opene.Metadata{
		"key": "value",
		"num": 42,
	}

	result := err.WithMetadata(md)
	assert.Equal(t, md, result.Meta)
	assert.Same(t, err, result)
}

func TestError_Sensitive(t *testing.T) {
	err := &opene.Error{Message: "test error"}
	result := err.Sensitive()

	assert.True(t, result.IsSensitive)
	assert.Same(t, err, result)
}

func TestError_Wrap(t *testing.T) {
	t.Run("wraps another error", func(t *testing.T) {
		inner := &opene.Error{
			Message:    "inner error",
			Code:       opene.CodeValidation,
			Domain:     "inner",
			Operation:  "test",
			StatusCode: http.StatusBadRequest,
		}

		outer := &opene.Error{
			Message:    "outer error",
			Code:       opene.CodeInternal,
			Domain:     "outer",
			Operation:  "wrap",
			StatusCode: http.StatusInternalServerError,
		}

		wrapped := outer.Wrap(inner)

		assert.Contains(t, wrapped.Message, "outer error")
		assert.Contains(t, wrapped.Message, "inner error")
		assert.Equal(t, opene.CodeInternal, wrapped.Code)
		assert.Equal(t, "outer", wrapped.Domain)
		assert.Equal(t, "wrap", wrapped.Operation)
		assert.Equal(t, http.StatusInternalServerError, wrapped.StatusCode)
		assert.Equal(t, inner, wrapped.WrappedErr)
	})

	t.Run("propagates sensitivity flag", func(t *testing.T) {
		inner := &opene.Error{
			Message:     "sensitive inner error",
			Code:        opene.CodeInternal,
			IsSensitive: true,
		}

		outer := &opene.Error{
			Message: "outer error",
			Code:    opene.CodeValidation,
		}

		wrapped := outer.Wrap(inner)
		assert.True(t, wrapped.IsSensitive)
	})

	t.Run("handles nil error", func(t *testing.T) {
		base := &opene.Error{
			Message: "base error",
			Code:    opene.CodeValidation,
			Domain:  "test",
		}

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
		original := &opene.Error{
			Message: "original error",
			Code:    opene.CodeValidation,
			Domain:  "original",
		}

		err := opene.AsError(original, "test", opene.CodeInternal)
		assert.Same(t, original, err)
	})

	t.Run("handles nil error", func(t *testing.T) {
		err := opene.AsError(nil, "test", opene.CodeValidation)
		assert.Nil(t, err)
	})
}
