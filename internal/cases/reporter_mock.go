package cases

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// NewReporterMock creates a new reporter mock.
func NewReporterMock() *ReporterMock {
	return &ReporterMock{}
}

// ReporterMock is a mock of reporter.
type ReporterMock struct {
	mock.Mock
}

// Report mocks base method.
func (m *ReporterMock) Report(c *Case, resp *http.Response) (err error) {
	args := m.Called(c, resp)

	err = args.Error(0)

	return
}