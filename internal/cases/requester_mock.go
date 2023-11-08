package cases

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// NewRequesterMock creates a new requester mock.
func NewRequesterMock() *RequesterMock {
	return &RequesterMock{}
}

// RequesterMock is a mock of requester.
type RequesterMock struct {
	mock.Mock
}

// Do mocks base method.
func (m *RequesterMock) Do(c *Case) (resp *http.Response, err error) {
	args := m.Called(c)

	resp = args.Get(0).(*http.Response)
	err = args.Error(1)

	return
}