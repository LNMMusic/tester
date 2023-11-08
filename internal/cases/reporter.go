package cases

import "net/http"

// Reporter is a reporter of test cases.
type Reporter interface {
	// Report reports the result of the test case.
	Report(c *Case, w *http.Response) (err error)
}