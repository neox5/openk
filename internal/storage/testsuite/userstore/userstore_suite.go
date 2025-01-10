package userstore

import (
	"testing"

	"github.com/neox5/openk/internal/storage"
)

// TestCase represents a single test definition
type TestCase struct {
	name string
	fn   func(*Suite, *testing.T)
}

var (
	// Basic operations
	CreateTest = TestCase{
		name: "Create",
		fn:   func(s *Suite, t *testing.T) { s.testCreateUser(t) },
	}
	GetByIDTest = TestCase{
		name: "GetByID",
		fn:   func(s *Suite, t *testing.T) { s.testGetByID(t) },
	}
	GetByUsernameTest = TestCase{
		name: "GetByUsername",
		fn:   func(s *Suite, t *testing.T) { s.testGetByUsername(t) },
	}
	ConcurrentTest = TestCase{
		name: "Concurrent",
		fn:   func(s *Suite, t *testing.T) { s.testConcurrent(t) },
	}
)

// AllTests is the default set of test cases
var AllTests = []TestCase{
	CreateTest,
	GetByIDTest,
	GetByUsernameTest,
	ConcurrentTest,
}

// Suite defines a test suite for UserStore implementations
type Suite struct {
	NewStore func() storage.UserStore
}

// RunTestCases executes the specified test cases
func (s *Suite) RunTestCases(t *testing.T, cases []TestCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fn(s, t)
		})
	}
}

// RunAll executes all test cases in the suite
func (s *Suite) RunAll(t *testing.T) {
	t.Helper()
	s.RunTestCases(t, AllTests)
}
