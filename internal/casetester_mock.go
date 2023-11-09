package internal

import (
	"github.com/LNMMusic/tester/internal/cases"
	"github.com/stretchr/testify/mock"
)

// NewCaseTesterMock creates a new CaseTesterMock.
func NewCaseTesterMock() (m *CaseTesterMock) {
	m = &CaseTesterMock{}
	return
}

// CaseTesterMock is a mock of CaseTester.
type CaseTesterMock struct {
	mock.Mock
}

// Test is a mock of Test.
func (m *CaseTesterMock) Test(c *cases.Case) (err error) {
	args := m.Called(c)

	err = args.Error(0)

	return
}