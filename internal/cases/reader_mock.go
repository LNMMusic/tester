package cases

import "github.com/stretchr/testify/mock"

// NewReaderMock creates a new mock of Reader.
func NewReaderMock() (m *ReaderMock) {
	m = &ReaderMock{}
	return
}

// ReaderMock is a mock of Reader.
type ReaderMock struct {
	mock.Mock
}

// Read is a mock of Read.
func (m *ReaderMock) Read() (c Case, err error) {
	args := m.Called()

	c = args.Get(0).(Case)
	err = args.Error(1)

	return
}