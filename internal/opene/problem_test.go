package opene_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/neox5/openk/internal/opene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsProblem(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("converts basic error to problem", func(t *testing.T) {
			err := &opene.Error{
				Message:    "validation failed",
				Code:       "VALIDATION_ERROR",
				Domain:     "validation",
				StatusCode: http.StatusBadRequest,
			}

			problem := opene.AsProblem(err)
			require.NotNil(t, problem)

			assert.Equal(t, "https://example.com/errors/validation", problem.Type)
			assert.Equal(t, "validation failed", problem.Title)
			assert.Equal(t, http.StatusBadRequest, problem.Status)
			assert.Equal(t, "validation failed", problem.Detail)
			assert.Empty(t, problem.Instance)
			assert.Nil(t, problem.Extra)
		})

		t.Run("includes metadata in extra field", func(t *testing.T) {
			err := &opene.Error{
				Message:    "invalid input",
				Code:       "VALIDATION_ERROR",
				Domain:     "validation",
				StatusCode: http.StatusBadRequest,
				Meta: opene.Metadata{
					"field": "email",
					"value": "invalid@example",
				},
			}

			problem := opene.AsProblem(err)
			require.NotNil(t, problem)

			extra, ok := problem.Extra.(opene.Metadata)
			require.True(t, ok)
			assert.Equal(t, "email", extra["field"])
			assert.Equal(t, "invalid@example", extra["value"])
		})

		t.Run("handles sensitive error", func(t *testing.T) {
			err := &opene.Error{
				Message:     "database connection failed: invalid credentials",
				Code:        "INTERNAL_ERROR",
				Domain:      "database",
				StatusCode:  http.StatusServiceUnavailable,
				IsSensitive: true,
				Meta: opene.Metadata{
					"database": "users",
					"credentials": map[string]string{
						"user": "admin",
						"host": "db.internal",
					},
				},
			}

			problem := opene.AsProblem(err)
			require.NotNil(t, problem)

			// Verify sensitive information is hidden
			assert.Equal(t, "https://example.com/errors/internal", problem.Type)
			assert.Equal(t, "Internal Server Error", problem.Title)
			assert.Equal(t, http.StatusInternalServerError, problem.Status)
			assert.Empty(t, problem.Detail)
			assert.Empty(t, problem.Instance)
			assert.Nil(t, problem.Extra)
		})

		t.Run("handles standard error", func(t *testing.T) {
			err := errors.New("standard error")
			problem := opene.AsProblem(err)
			require.NotNil(t, problem)

			assert.Equal(t, "https://example.com/errors/internal", problem.Type)
			assert.Equal(t, "Internal Server Error", problem.Title)
			assert.Equal(t, http.StatusInternalServerError, problem.Status)
			assert.Empty(t, problem.Detail)
			assert.Empty(t, problem.Instance)
			assert.Nil(t, problem.Extra)
		})

		t.Run("uses configured error base uri", func(t *testing.T) {
			// Set custom base URI
			opene.SetErrorBaseURI("https://api.myapp.com/errors/")

			err := &opene.Error{
				Message: "test error",
				Domain:  "test",
			}

			problem := opene.AsProblem(err)
			assert.Equal(t, "https://api.myapp.com/errors/test", problem.Type)

			// Reset to default
			opene.SetErrorBaseURI("https://example.com/errors/")
		})

		t.Run("handles empty domain", func(t *testing.T) {
			err := &opene.Error{
				Message: "test error",
				// No domain specified
			}

			problem := opene.AsProblem(err)
			assert.Equal(t, "https://example.com/errors/", problem.Type)
		})
	})
}

func TestSetErrorBaseURI(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("updates problem type uri", func(t *testing.T) {
			// Save original
			original := "https://example.com/errors/"

			// Set new URI
			opene.SetErrorBaseURI("https://errors.myapp.com/types/")

			err := &opene.Error{
				Message: "test error",
				Domain:  "test",
			}

			problem := opene.AsProblem(err)
			assert.Equal(t, "https://errors.myapp.com/types/test", problem.Type)

			// Reset to original
			opene.SetErrorBaseURI(original)
		})

		t.Run("handles uri without trailing slash", func(t *testing.T) {
			original := "https://example.com/errors/"

			opene.SetErrorBaseURI("https://api.myapp.com/errors")

			err := &opene.Error{
				Message: "test error",
				Domain:  "test",
			}

			problem := opene.AsProblem(err)
			assert.Equal(t, "https://api.myapp.com/errors/test", problem.Type)

			opene.SetErrorBaseURI(original)
		})
	})
}
