package internal_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/LNMMusic/tester/internal"
	"github.com/LNMMusic/tester/internal/cases"

	"github.com/stretchr/testify/require"
)

// Tests for Tester Test method
func TestTester_Test(t *testing.T) {
	t.Run("case 1: success to test", func(t *testing.T) {
		// arrange
		// - dbexecuter
		db := cases.NewDbExecuterMock()
		db.On("Exec", []string{"query 1", "query 2"}).Return(nil)
		db.On("Exec", []string{"query 3", "query 4"}).Return(nil)
		// - requester
		rq := cases.NewRequesterMock()
		rq.On("Do", &cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		}).Return(&http.Response{}, nil)
		// - reporter
		rp := cases.NewReporterMock()
		rp.On("Report", &cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		}, &http.Response{}).Return(nil)
		// - tester
		ts := internal.NewTester(db, rq, rp)

		// act
		err := ts.Test(&cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		})

		// assert
		require.NoError(t, err)
		db.AssertExpectations(t)
		rq.AssertExpectations(t)
		rp.AssertExpectations(t)
	})

	t.Run("case 2: fail to test - database set-up", func(t *testing.T) {
		// arrange
		// - dbexecuter
		db := cases.NewDbExecuterMock()
		db.On("Exec", []string{"query 1", "query 2"}).Return(errors.New("dbexecuter: internal error"))
		db.On("Exec", []string{"query 3", "query 4"}).Return(nil)
		// - requester
		// ...
		// - reporter
		// ...
		// - tester
		ts := internal.NewTester(db, nil, nil)

		// act
		err := ts.Test(&cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		})

		// assert
		require.ErrorIs(t, err, internal.ErrTesterDatabase)
		require.EqualError(t, err, "tester: database error. dbexecuter: internal error")
		db.AssertExpectations(t)
	})

	t.Run("case 3: fail to test - database tear-down", func(t *testing.T) {
		// arrange
		// - dbexecuter
		db := cases.NewDbExecuterMock()
		db.On("Exec", []string{"query 1", "query 2"}).Return(nil)
		db.On("Exec", []string{"query 3", "query 4"}).Return(errors.New("dbexecuter: internal error"))
		// - requester
		rq := cases.NewRequesterMock()
		rq.On("Do", &cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		}).Return(&http.Response{}, nil)
		// - reporter
		rp := cases.NewReporterMock()
		rp.On("Report", &cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		}, &http.Response{}).Return(nil)
		// - tester
		ts := internal.NewTester(db, rq, rp)
			
		// act
		err := ts.Test(&cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		})

		// assert
		require.ErrorIs(t, err, internal.ErrTesterDatabase)
		require.EqualError(t, err, "tester: database error. dbexecuter: internal error. %!w(<nil>)")
		db.AssertExpectations(t)
		rq.AssertExpectations(t)
		rp.AssertExpectations(t)
	})

	t.Run("case 4: fail to test - requester", func(t *testing.T) {
		// arrange
		// - dbexecuter
		db := cases.NewDbExecuterMock()
		db.On("Exec", []string{"query 1", "query 2"}).Return(nil)
		db.On("Exec", []string{"query 3", "query 4"}).Return(nil)
		// - requester
		rq := cases.NewRequesterMock()
		rq.On("Do", &cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		}).Return((*http.Response)(nil), errors.New("requester: internal error"))
		// - reporter
		// ...
		// - tester
		ts := internal.NewTester(db, rq, nil)

		// act
		err := ts.Test(&cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		})

		// assert
		require.ErrorIs(t, err, internal.ErrTesterRequest)
		require.EqualError(t, err, "tester: request error. requester: internal error")
		db.AssertExpectations(t)
		rq.AssertExpectations(t)
	})

	t.Run("case 5: fail to test - reporter", func(t *testing.T) {
		// arrange
		// - dbexecuter
		db := cases.NewDbExecuterMock()
		db.On("Exec", []string{"query 1", "query 2"}).Return(nil)
		db.On("Exec", []string{"query 3", "query 4"}).Return(nil)
		// - requester
		rq := cases.NewRequesterMock()
		rq.On("Do", &cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		}).Return(&http.Response{}, nil)
		// - reporter
		rp := cases.NewReporterMock()
		rp.On("Report", &cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		}, &http.Response{}).Return(errors.New("reporter: internal error"))
		// - tester
		ts := internal.NewTester(db, rq, rp)

		// act
		err := ts.Test(&cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		})

		// assert
		require.ErrorIs(t, err, internal.ErrTesterReporter)
		require.EqualError(t, err, "tester: reporter error. reporter: internal error")
		db.AssertExpectations(t)
		rq.AssertExpectations(t)
		rp.AssertExpectations(t)
	})

	t.Run("case 6: fail to test - multiple error - reporter and database tear-down", func(t *testing.T) {
		// arrange
		// - dbexecuter
		db := cases.NewDbExecuterMock()
		db.On("Exec", []string{"query 1", "query 2"}).Return(nil)
		db.On("Exec", []string{"query 3", "query 4"}).Return(errors.New("dbexecuter: internal error"))
		// - requester
		rq := cases.NewRequesterMock()
		rq.On("Do", &cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		}).Return(&http.Response{}, nil)
		// - reporter
		rp := cases.NewReporterMock()
		rp.On("Report", &cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		}, &http.Response{}).Return(errors.New("reporter: internal error"))
		// - tester
		ts := internal.NewTester(db, rq, rp)

		// act
		err := ts.Test(&cases.Case{
			Database: cases.Database{
				SetUp: []string{"query 1", "query 2"},
				TearDown: []string{"query 3", "query 4"},
			},
		})

		// assert
		require.ErrorIs(t, err, internal.ErrTesterDatabase)
		require.ErrorIs(t, err, internal.ErrTesterReporter)
		require.EqualError(t, err, "tester: database error. dbexecuter: internal error. tester: reporter error. reporter: internal error")
		db.AssertExpectations(t)
		rq.AssertExpectations(t)
		rp.AssertExpectations(t)
	})
}