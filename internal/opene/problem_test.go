package opene_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/neox5/openk/internal/opene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetErrorBaseURI(t *testing.T) {
	// Save original for cleanup
	defer opene.ResetErrorBaseURI()

	t.Run("accepts valid URI", func(t *testing.T) {
		err := opene.SetErrorBaseURI("https://api.example.com/errors")
		assert.Nil(t, err)

		// Verify through a problem conversion
		prob := opene.AsProblem(opene.NewValidationError("test"))
		assert.Equal(t, "https://api.example.com/errors/validation", prob.Type)
	})

	t.Run("accepts http protocol", func(t *testing.T) {
		err := opene.SetErrorBaseURI("http://api.example.com/errors")
		assert.Nil(t, err)
	})

	t.Run("removes trailing slash", func(t *testing.T) {
		err := opene.SetErrorBaseURI("https://api.example.com/errors/")
		assert.Nil(t, err)

		prob := opene.AsProblem(opene.NewValidationError("test"))
		assert.Equal(t, "https://api.example.com/errors/validation", prob.Type)
	})

	t.Run("rejects empty URI", func(t *testing.T) {
		err := opene.SetErrorBaseURI("")
		require.NotNil(t, err)
		assert.Equal(t, opene.CodeValidation, err.Code)
		assert.Equal(t, "opene", err.Domain)
		assert.Equal(t, "set_base_uri", err.Operation)
		assert.Equal(t, "", err.Meta["provided_uri"])
	})

	t.Run("rejects URI without protocol", func(t *testing.T) {
		uri := "api.example.com/errors"
		err := opene.SetErrorBaseURI(uri)
		require.NotNil(t, err)
		assert.Equal(t, opene.CodeValidation, err.Code)
		assert.Equal(t, "opene", err.Domain)
		assert.Equal(t, "set_base_uri", err.Operation)
		assert.Equal(t, uri, err.Meta["provided_uri"])
	})
}

func TestAsProblem(t *testing.T) {
	t.Run("converts basic error to problem", func(t *testing.T) {
		err := &opene.Error{
			Message:    "validation failed",
			Code:       opene.CodeValidation,
			Domain:     "test",
			Operation:  "validate",
			StatusCode: http.StatusBadRequest,
			Meta: opene.Metadata{
				"field": "username",
			},
		}

		prob := opene.AsProblem(err)
		require.NotNil(t, prob)

		assert.Equal(t, "https://openk.dev/errors/validation", prob.Type)
		assert.Equal(t, "validation failed", prob.Title)
		assert.Equal(t, http.StatusBadRequest, prob.Status)
		assert.Equal(t, "validation failed", prob.Detail)

		meta, ok := prob.Extra.(opene.Metadata)
		require.True(t, ok)
		assert.Equal(t, "username", meta["field"])
	})

	t.Run("handles sensitive error", func(t *testing.T) {
		err := &opene.Error{
			Message:     "database error",
			Code:        opene.CodeInternal,
			Domain:      "db",
			Operation:   "query",
			StatusCode:  http.StatusInternalServerError,
			IsSensitive: true,
			Meta: opene.Metadata{
				"query": "SELECT * FROM users",
			},
		}

		prob := opene.AsProblem(err)
		require.NotNil(t, prob)

		assert.Equal(t, "https://openk.dev/errors/internal", prob.Type)
		assert.Equal(t, "Internal Server Error", prob.Title)
		assert.Equal(t, http.StatusInternalServerError, prob.Status)
		assert.Empty(t, prob.Detail)
		assert.Nil(t, prob.Extra)
	})

	t.Run("handles standard error", func(t *testing.T) {
		err := errors.New("standard error")
		prob := opene.AsProblem(err)
		require.NotNil(t, prob)

		assert.Equal(t, "https://openk.dev/errors/internal", prob.Type)
		assert.Equal(t, "Internal Server Error", prob.Title)
		assert.Equal(t, http.StatusInternalServerError, prob.Status)
		assert.Empty(t, prob.Detail)
		assert.Nil(t, prob.Extra)
	})
}
