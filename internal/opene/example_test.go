package opene_test

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/neox5/openk/internal/opene"
)

func ExampleNewValidationError() {
	err := opene.NewValidationError("invalid username format")
	fmt.Printf("Message: %s\nCode: %s\nStatus: %d\n", err.Message, err.Code, err.StatusCode)
	// Output:
	// Message: invalid username format
	// Code: validation
	// Status: 400
}

func ExampleNewNotFoundError() {
	err := opene.NewNotFoundError("user not found")
	fmt.Printf("Message: %s\nCode: %s\nStatus: %d\n", err.Message, err.Code, err.StatusCode)
	// Output:
	// Message: user not found
	// Code: not_found
	// Status: 404
}

func ExampleNewConflictError() {
	err := opene.NewConflictError("user already exists")
	fmt.Printf("Message: %s\nCode: %s\nStatus: %d\n", err.Message, err.Code, err.StatusCode)
	// Output:
	// Message: user already exists
	// Code: conflict
	// Status: 409
}

func ExampleNewInternalError() {
	err := opene.NewInternalError("database connection failed")
	fmt.Printf("Message: %s\nCode: %s\nStatus: %d\nSensitive: %v\n", err.Message, err.Code, err.StatusCode, err.IsSensitive)
	// Output:
	// Message: database connection failed
	// Code: internal
	// Status: 500
	// Sensitive: true
}

func ExampleError_WithMetadata() {
	err := opene.NewValidationError("validation failed").
		WithMetadata(opene.Metadata{
			"field": "age",
			"value": -5,
		})
	fmt.Printf("Field: %v\nValue: %v\n", err.Meta["field"], err.Meta["value"])
	// Output:
	// Field: age
	// Value: -5
}

func ExampleError_WithDomain() {
	err := opene.NewInternalError("key generation failed").
		WithDomain("crypto")
	fmt.Printf("Domain: %s\nCode: %s\n", err.Domain, err.Code)
	// Output:
	// Domain: crypto
	// Code: internal
}

func ExampleError_WithOperation() {
	err := opene.NewInternalError("key generation failed").
		WithOperation("generate_key")
	fmt.Printf("Operation: %s\nCode: %s\n", err.Operation, err.Code)
	// Output:
	// Operation: generate_key
	// Code: internal
}

func ExampleError_Wrap() {
	dbErr := errors.New("connection refused")
	storageErr := opene.NewInternalError("database error")
	finalErr := storageErr.Wrap(opene.AsError(dbErr, "db", opene.CodeInternal))
	fmt.Println(finalErr.Error())
	// Output:
	// database error: connection refused
}

func ExampleAsError() {
	stdErr := errors.New("file not found")
	err := opene.AsError(stdErr, "filesystem", opene.CodeNotFound)
	fmt.Printf("Domain: %s\nCode: %s\nMessage: %s\n", err.Domain, err.Code, err.Message)
	// Output:
	// Domain: filesystem
	// Code: not_found
	// Message: file not found
}

func ExampleAsProblem() {
	err := opene.NewValidationError("invalid input").
		WithDomain("auth").
		WithOperation("validate_credentials").
		WithMetadata(opene.Metadata{
			"field": "password",
		})
	prob := opene.AsProblem(err)
	fmt.Printf("Type: %s\nStatus: %d\nTitle: %s\n", prob.Type, prob.Status, prob.Title)
	// Output:
	// Type: https://openk.dev/errors/validation
	// Status: 400
	// Title: invalid input
}

func ExampleAsProblem_sensitive() {
	err := opene.NewInternalError("database error: invalid credentials").
		WithDomain("db").
		WithOperation("query").
		WithMetadata(opene.Metadata{
			"database": "users",
			"host":     "internal.db",
		})

	prob := opene.AsProblem(err)
	fmt.Printf("Type: %s\nTitle: %s\nStatus: %d\n", prob.Type, prob.Title, prob.Status)
	// Output:
	// Type: https://openk.dev/errors/internal
	// Title: Internal Server Error
	// Status: 500
}

func ExampleAsProblem_http() {
	w := &mockResponseWriter{}
	r, _ := http.NewRequest("GET", "/test?id=invalid", nil)

	// Simulated handler that checks a query parameter
	err := func(r *http.Request) error {
		if id := r.URL.Query().Get("id"); id == "invalid" {
			return opene.NewValidationError("invalid parameter").
				WithDomain("http").
				WithOperation("validate_query").
				WithMetadata(opene.Metadata{
					"param": "id",
					"value": id,
				})
		}
		return nil
	}(r)
	if err != nil {
		prob := opene.AsProblem(err)
		w.WriteHeader(prob.Status)
		fmt.Printf("Response Status: %d\nProblem Type: %s\n", prob.Status, prob.Type)
	}
	// Output:
	// Response Status: 400
	// Problem Type: https://openk.dev/errors/validation
}

// mockResponseWriter for http examples
type mockResponseWriter struct {
	status int
}

func (m *mockResponseWriter) Header() http.Header        { return make(http.Header) }
func (m *mockResponseWriter) Write([]byte) (int, error)  { return 0, nil }
func (m *mockResponseWriter) WriteHeader(statusCode int) { m.status = statusCode }
