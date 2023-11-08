package internal

import (
	"errors"
	"fmt"
	"net/http"

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

// NewTester creates a new tester.
func NewTester(dbExecuter cases.DbExecuter, requester cases.Requester, reporter cases.Reporter) *Tester {
	return &Tester{
		dbExecuter: dbExecuter,
		requester: requester,
		reporter: reporter,
	}
}

// Tester is a tester of cases for http servers.
type Tester struct {
	// dbexecuter is the database executer of test cases.
	dbExecuter cases.DbExecuter
	// requester is the requester of test cases.
	requester cases.Requester
	// reporter is the reporter of test cases.
	reporter cases.Reporter
}

// Test tests the server.
func (t *Tester) Test(c *cases.Case) (err error) {
	// arrange
	// - database: tear down
	defer func() {
		e := t.dbExecuter.Exec(c.Database.TearDown...)
		if e != nil {
			// case of multiple errors:
			// - wrap errors in a slice of errors
			// 	 > from `{msg string;err error}`
			//   > to `{msg string;errs []error}`
			// 
			// - notes to take on consideration:
			//   > generation process:
			//     > both unrelated multiple errors are wrapped -> []error{ErrDefer, ErrOther}
			//   > identification process:
			//     > verification order matters
			//     > rule #1: accept the idea the last error might not be identified as first with errors.Is
			//                regardless of the order of the errors, both are unrelated, so is not wrong to identify the first error or the last error
			//                the full message of the error is still intact
			//     > advantages: we can use errors.As to get more details about the inner multiple unrelated errors. More programmatic control.
			err = fmt.Errorf("%w. %v. %w", ErrTesterDatabase, e, err)
		}
	}()
	// - database: set up
	err = t.dbExecuter.Exec(c.Database.SetUp...)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrTesterDatabase, err)
		return
	}

	// act
	var resp *http.Response
	resp, err = t.requester.Do(c)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrTesterRequest, err)
		return
	}

	// assert
	err = t.reporter.Report(c, resp)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrTesterReporter, err)
		return
	}

	return
}
