package internal

import (
	"errors"

	"github.com/LNMMusic/tester/internal/cases"
)


var (
	// ErrTesterDatabase is the error of the database.
	ErrTesterDatabase = errors.New("tester: database error")
	// ErrTesterRequest is the error of the request.
	ErrTesterRequest = errors.New("tester: request error")
	// ErrTesterReporter is the error of the reporter.
	ErrTesterReporter = errors.New("tester: reporter error")
)

// CaseTester is an interface that test a case.
type CaseTester interface {
	// Test tests a case.
	Test(c *cases.Case) (err error)
}