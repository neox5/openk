package middleware_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierror "github.com/neox5/openk/internal/api/error"
	"github.com/neox5/openk/internal/server/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandler(t *testing.T) {
	t.Run("handles normal response", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})

		wrapped := middleware.ErrorHandler(handler)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)

		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())
	})

	t.Run("handles APIError", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := w.(interface {
				WriteError(error)
			})
			rw.WriteError(apierror.ValidationError("test", "invalid"))
		})

		wrapped := middleware.ErrorHandler(handler)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)

		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "application/problem+json")

		// Verify response structure
		var resp map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Equal(t, "https://openk.dev/errors/validation", resp["type"])
		assert.Equal(t, float64(400), resp["status"])
		assert.Equal(t, "/test", resp["instance"])
	})

	t.Run("handles standard error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := w.(interface {
				WriteError(error)
			})
			rw.WriteError(errors.New("something went wrong"))
		})

		wrapped := middleware.ErrorHandler(handler)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)

		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "application/problem+json")

		var resp map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Equal(t, "https://openk.dev/errors/internal", resp["type"])
		assert.Equal(t, "something went wrong", resp["detail"])
	})

	t.Run("handles panic", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("unexpected error")
		})

		wrapped := middleware.ErrorHandler(handler)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)

		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "application/problem+json")

		var resp map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Equal(t, "https://openk.dev/errors/internal", resp["type"])
		assert.Equal(t, "unexpected error", resp["extra"].(map[string]interface{})["panic"])
	})
}
