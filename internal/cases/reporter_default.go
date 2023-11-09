package cases

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// NewReporterDefault creates a new default reporter.
func NewReporterDefault(excludedHeaders []string) *ReporterDefault {
	// default excluded headers
	defaultExcludedHeaders := []string{"Date", "Content-Length"}
	if excludedHeaders != nil {
		if len(excludedHeaders) > 0 {
			defaultExcludedHeaders = excludedHeaders
		}
	}

	return &ReporterDefault{
		excludedHeaders: defaultExcludedHeaders,
	}
}

// ReporterDefault is the default reporter of test cases.
type ReporterDefault struct {
	// excluded headers
	excludedHeaders []string
}

func (r *ReporterDefault) Report(c *Case, w *http.Response) (err error) {
	// expectations
	expectedCode := c.Response.Code
	expectedBody := c.Response.Body
	expectedHeader := c.Response.Header
	// actual
	actualCode := w.StatusCode
	actualBody := any(nil)
	err = json.NewDecoder(w.Body).Decode(&actualBody)
	if err != nil {
		return err
	}
	actualHeader := w.Header

	// exclusions
	for _, h := range r.excludedHeaders {
		delete(expectedHeader, h)
		delete(actualHeader, h)
	}

	// verify
	validCode := expectedCode == actualCode
	validBody := reflect.DeepEqual(expectedBody, actualBody)
	validHeader := reflect.DeepEqual(expectedHeader, actualHeader)
	if !(validCode && validBody && validHeader) {
		fmt.Printf("> Case '%s': FAIL\n", c.Name)

		if !validCode {
			fmt.Printf("- expected code: %d\n", expectedCode)
			fmt.Printf("- actual code: %d\n", actualCode)
		}
		if !validBody {
			fmt.Printf("- expected body: %v\n", expectedBody)
			fmt.Printf("- actual body: %v\n", actualBody)
		}
		if !validHeader {
			fmt.Printf("- expected header: %v\n", expectedHeader)
			fmt.Printf("- actual header: %v\n", actualHeader)
		}
		fmt.Println()
		return
	}

	fmt.Printf("> Case '%s': PASS\n", c.Name)
	fmt.Println()

	return nil
}