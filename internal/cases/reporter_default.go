package cases

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// NewReporterDefault creates a new default reporter.
func NewReporterDefault() *ReporterDefault {
	return &ReporterDefault{}
}

// ReporterDefault is the default reporter of test cases.
type ReporterDefault struct {}

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

	// verify
	fmt.Println("> case:", c.Name)
	if expectedCode != actualCode {
		fmt.Println("- expected code: ", expectedCode)
		fmt.Println("- actual code: ", actualCode)
	}
	if !reflect.DeepEqual(expectedBody, actualBody) {
		fmt.Println("- expected body: ", expectedBody)
		fmt.Println("- actual body: ", actualBody)
	}
	if !reflect.DeepEqual(expectedHeader, actualHeader) {
		fmt.Println("- expected header: ", expectedHeader)
		fmt.Println("- actual header: ", actualHeader)
	}
	
	return nil
}