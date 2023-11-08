package cases

import (
	"errors"
	"net/http"
)

// Database is a database to run the test case against.
type Database struct {
	// SetUp is the set of set up queries to run before the test case.
	SetUp []string `json:"set_up"`
	// TearDown is the set of tear down queries to run after the test case.
	TearDown []string `json:"tear_down"`
}

// Request is a request to make for the test case.
type Request struct {
	// Method is the HTTP method to use for the request.
	Method string `json:"method"`
	// Path is the path to use for the request.
	Path string `json:"path"`
	// Query is the set of query parameters to use for the request.
	Query map[string]string `json:"query"`
	// Body is the body to use for the request.
	Body any `json:"body"`
	// Header is the set of headers to use for the request.
	Header http.Header `json:"header"`
}

// Response is the expected response of the test case.
type Response struct {
	// Code is the expected status code of the response.
	Code int `json:"code"`
	// Body is the expected body of the response.
	Body any `json:"body"`
	// Header is the expected set of headers of the response.
	Header http.Header `json:"header"`
}

// Case is a test case.
type Case struct {
	// Name is the name of the test case.
	Name string `json:"case_name"`
	// Arrange
	Database `json:"database"`
	// Input
	Request `json:"request"`
	// Output
	Response `json:"response"`
}

var (
	// ErrEndOfLine is the error returned when the end of the line is reached.
	ErrEndOfLine = errors.New("end of line")
)

// Reader is a reader of test cases.
type Reader interface {
	// Read reads the next test case.
	Read() (c Case, err error)
}