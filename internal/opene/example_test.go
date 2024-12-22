package opene_test

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/neox5/openk/internal/opene"
)

func ExampleNewValidationError() {
	err := opene.NewValidationError("auth", "validate_username", "invalid username format")
	fmt.Printf("Message: %s\nCode: %s\nDomain: %s\nOperation: %s\nStatus: %d\n",
		err.Message, err.Code, err.Domain, err.Operation, err.StatusCode)
	// Output:
	// Message: invalid username format
	// Code: validation
	// Domain: auth
	// Operation: validate_username
	// Status: 400
}

func ExampleNewNotFoundError() {
	err := opene.NewNotFoundError("user", "fetch", "user not found")
	fmt.Printf("Message: %s\nCode: %s\nDomain: %s\nOperation: %s\nStatus: %d\n",
		err.Message, err.Code, err.Domain, err.Operation, err.StatusCode)
	// Output:
	// Message: user not found
	// Code: not_found
	// Domain: user
	// Operation: fetch
	// Status: 404
}

func ExampleNewConflictError() {
	err := opene.NewConflictError("user", "create", "user already exists")
	fmt.Printf("Message: %s\nCode: %s\nDomain: %s\nOperation: %s\nStatus: %d\n",
		err.Message, err.Code, err.Domain, err.Operation, err.StatusCode)
	// Output:
	// Message: user already exists
	// Code: conflict
	// Domain: user
	// Operation: create
	// Status: 409
}

func ExampleNewInternalError() {
	err := opene.NewInternalError("db", "connect", "database connection failed")
	fmt.Printf("Message: %s\nCode: %s\nDomain: %s\nOperation: %s\nStatus: %d\nSensitive: %v\n",
		err.Message, err.Code, err.Domain, err.Operation, err.StatusCode, err.IsSensitive)
	// Output:
	// Message: database connection failed
	// Code: internal
	// Domain: db
	// Operation: connect
	// Status: 500
	// Sensitive: true
}

func ExampleError_WithMetadata() {
	err := opene.NewValidationError("auth", "validate", "validation failed").
		WithMetadata(opene.Metadata{
			"field": "age",
			"value": -5,
		})
	fmt.Printf("Field: %v\nValue: %v\n", err.Meta["field"], err.Meta["value"])
	// Output:
	// Field: age
	// Value: -5
}

func ExampleError_Wrap() {
	dbErr := errors.New("connection refused")
	inner := opene.AsError(dbErr, "db", opene.CodeInternal)
	outer := opene.NewInternalError("storage", "connect", "database error").Wrap(inner)
	fmt.Println(outer.Error())
	// Output:
	// database error: connection refused
}

func ExampleAsError() {
	stdErr := errors.New("file not found")
	err := opene.AsError(stdErr, "filesystem", opene.CodeNotFound)
	fmt.Printf("Domain: %s\nCode: %s\nMessage: %s\n",
		err.Domain, err.Code, err.Message)
	// Output:
	// Domain: filesystem
	// Code: not_found
	// Message: file not found
}

func ExampleAsProblem() {
	err := opene.NewValidationError("auth", "validate_credentials", "invalid input").
		WithMetadata(opene.Metadata{
			"field": "password",
		})
	prob := opene.AsProblem(err)
	fmt.Printf("Type: %s\nStatus: %d\nTitle: %s\n",
		prob.Type, prob.Status, prob.Title)
	// Output:
	// Type: https://openk.dev/errors/validation
	// Status: 400
	// Title: invalid input
}

func ExampleAsProblem_http() {
	w := &mockResponseWriter{}
	r, _ := http.NewRequest("GET", "/test?id=invalid", nil)

	// Simulated handler that checks a query parameter
	err := func(r *http.Request) error {
		if id := r.URL.Query().Get("id"); id == "invalid" {
			return opene.NewValidationError("http", "validate_query", "invalid parameter").
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
