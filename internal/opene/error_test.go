package opene_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/neox5/openk/internal/opene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError_Error(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("returns message for simple error", func(t *testing.T) {
			err := &opene.Error{
				Message: "test error",
				Code:    opene.CodeValidationError,
			}
			assert.Equal(t, "test error", err.Error())
		})

		t.Run("returns message with metadata present", func(t *testing.T) {
			err := &opene.Error{
				Message: "error with context",
				Code:    opene.CodeValidationError,
				Meta: opene.Metadata{
					"key": "value",
				},
			}
			assert.Equal(t, "error with context", err.Error())
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("handles empty message", func(t *testing.T) {
			err := &opene.Error{
				Message: "",
				Code:    opene.CodeValidationError,
			}
			assert.Equal(t, "", err.Error())
		})
	})
}

func TestError_Unwrap(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("unwraps Error preserving type", func(t *testing.T) {
			inner := &opene.Error{
				Message: "inner error",
				Code:    opene.CodeValidationError,
			}
			outer := &opene.Error{
				Message:    "outer error",
				Code:       opene.CodeInternalError,
				WrappedErr: inner,
			}

			unwrapped := outer.Unwrap()
			require.NotNil(t, unwrapped)

			innerErr, ok := unwrapped.(*opene.Error)
			require.True(t, ok)
			assert.Equal(t, opene.CodeValidationError, innerErr.Code)
		})

		t.Run("unwraps standard error", func(t *testing.T) {
			stdErr := errors.New("standard error")
			err := &opene.Error{
				Message:    "wrapper",
				Code:       opene.CodeValidationError,
				WrappedErr: stdErr,
			}

			unwrapped := err.Unwrap()
			assert.Equal(t, stdErr, unwrapped)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("returns nil for no wrapped error", func(t *testing.T) {
			err := &opene.Error{
				Message: "no wrapped error",
				Code:    opene.CodeValidationError,
			}
			assert.Nil(t, err.Unwrap())
		})
	})
}

func TestError_UnwrapAll(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("gets to root cause through Error chain", func(t *testing.T) {
			rootErr := errors.New("root cause")

			inner := &opene.Error{
				Message:    "inner error",
				Code:       opene.CodeValidationError,
				WrappedErr: rootErr,
			}

			outer := &opene.Error{
				Message:    "outer error",
				Code:       opene.CodeInternalError,
				WrappedErr: inner,
			}

			result := outer.UnwrapAll()
			assert.Equal(t, rootErr, result)
		})

		t.Run("handles single Error wrapping standard error", func(t *testing.T) {
			stdErr := errors.New("standard error")
			err := &opene.Error{
				Message:    "wrapper",
				Code:       opene.CodeValidationError,
				WrappedErr: stdErr,
			}

			result := err.UnwrapAll()
			assert.Equal(t, stdErr, result)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("returns nil for no wrapped error", func(t *testing.T) {
			err := &opene.Error{
				Message: "no wrapped error",
				Code:    opene.CodeValidationError,
			}
			assert.Nil(t, err.UnwrapAll())
		})

		t.Run("returns nil for Error chain with no root cause", func(t *testing.T) {
			inner := &opene.Error{
				Message: "inner error",
				Code:    opene.CodeValidationError,
			}
			outer := &opene.Error{
				Message:    "outer error",
				Code:       opene.CodeInternalError,
				WrappedErr: inner,
			}
			assert.Nil(t, outer.UnwrapAll())
		})
	})
}

func TestError_WithMetadata(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("adds new metadata", func(t *testing.T) {
			err := &opene.Error{
				Message: "test error",
				Code:    opene.CodeValidationError,
			}
			md := opene.Metadata{
				"key": "value",
				"num": 42,
			}

			result := err.WithMetadata(md)
			assert.Equal(t, md, result.Meta)
		})

		t.Run("replaces existing metadata", func(t *testing.T) {
			err := &opene.Error{
				Message: "test error",
				Code:    opene.CodeValidationError,
				Meta: opene.Metadata{
					"old": "value",
				},
			}
			newMd := opene.Metadata{
				"new": "value",
			}

			result := err.WithMetadata(newMd)
			assert.Equal(t, newMd, result.Meta)
			assert.NotContains(t, result.Meta, "old")
		})
	})
}

func TestError_Wrap(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("wraps another error", func(t *testing.T) {
			inner := &opene.Error{
				Message:    "inner error",
				Code:       opene.CodeValidationError,
				Domain:     "inner",
				StatusCode: http.StatusBadRequest,
			}

			outer := &opene.Error{
				Message:    "outer error",
				Code:       opene.CodeInternalError,
				Domain:     "outer",
				StatusCode: http.StatusInternalServerError,
			}

			wrapped := outer.Wrap(inner)

			assert.Contains(t, wrapped.Message, "outer error")
			assert.Contains(t, wrapped.Message, "inner error")
			assert.Equal(t, opene.CodeInternalError, wrapped.Code)
			assert.Equal(t, "outer", wrapped.Domain)
			assert.Equal(t, http.StatusInternalServerError, wrapped.StatusCode)
			assert.Equal(t, inner, wrapped.WrappedErr)
		})

		t.Run("propagates sensitivity flag", func(t *testing.T) {
			inner := &opene.Error{
				Message:     "sensitive inner error",
				Code:        opene.CodeInternalError,
				IsSensitive: true,
			}

			outer := &opene.Error{
				Message: "outer error",
				Code:    opene.CodeValidationError,
			}

			wrapped := outer.Wrap(inner)
			assert.True(t, wrapped.IsSensitive)
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("handles nil error", func(t *testing.T) {
			base := &opene.Error{
				Message: "base error",
				Code:    opene.CodeValidationError,
				Domain:  "test",
			}

			wrapped := base.Wrap(nil)
			assert.Same(t, base, wrapped)
		})
	})
}

func TestError_Sensitive(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("marks error as sensitive", func(t *testing.T) {
			err := &opene.Error{
				Message: "test error",
				Code:    opene.CodeValidationError,
			}
			result := err.Sensitive()
			assert.True(t, result.IsSensitive)
			assert.Same(t, err, result)
		})
	})
}

func TestAsError(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("converts standard error", func(t *testing.T) {
			stdErr := errors.New("standard error")
			err := opene.AsError(stdErr, "test", opene.CodeValidationError)
			require.NotNil(t, err)

			assert.Equal(t, "standard error", err.Message)
			assert.Equal(t, opene.CodeValidationError, err.Code)
			assert.Equal(t, "test", err.Domain)
			assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
			assert.False(t, err.IsSensitive)
			assert.Equal(t, stdErr, err.WrappedErr)
		})

		t.Run("preserves existing error", func(t *testing.T) {
			original := &opene.Error{
				Message: "original error",
				Code:    opene.CodeValidationError,
				Domain:  "original",
			}

			err := opene.AsError(original, "test", opene.CodeInternalError)
			assert.Same(t, original, err)
		})

		t.Run("preserves error chain", func(t *testing.T) {
			inner := errors.New("inner")
			middle := fmt.Errorf("middle: %w", inner)
			outer := fmt.Errorf("outer: %w", middle)

			err := opene.AsError(outer, "test", opene.CodeValidationError)
			require.NotNil(t, err)

			assert.Equal(t, outer, err.WrappedErr)
			assert.Equal(t, middle, errors.Unwrap(err.WrappedErr))
			assert.Equal(t, inner, errors.Unwrap(errors.Unwrap(err.WrappedErr)))
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("handles nil error", func(t *testing.T) {
			err := opene.AsError(nil, "test", opene.CodeValidationError)
			assert.Nil(t, err)
		})
	})
}
