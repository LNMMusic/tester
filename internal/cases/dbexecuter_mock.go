package cases

import "github.com/stretchr/testify/mock"

// NewDbExecuterMock creates a new dbexecuter mock.
func NewDbExecuterMock() *DbExecuterMock {
	return &DbExecuterMock{}
}

// DbExecuterMock is a mock of dbexecuter.
type DbExecuterMock struct {
	mock.Mock
}

// Exec mocks base method.
func (m *DbExecuterMock) Exec(queries ...string) (err error) {
	args := m.Called(queries)

	err = args.Error(0)

	return
}