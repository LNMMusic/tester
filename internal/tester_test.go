package internal_test

import (
	"testing"

	"github.com/LNMMusic/tester/internal"
	"github.com/LNMMusic/tester/internal/cases"
	"github.com/stretchr/testify/require"
)

// Tests for Tester Run method.
func TestTester_Run(t *testing.T) {
	t.Run("case 1: success - all cases processed", func(t *testing.T) {
		// arrange
		// - reader: mock
		rd := cases.NewReaderMock()
		rd.On("Read").Return(cases.Case{}, cases.ErrEndOfLine)
		// - casetester: mock
		ct := internal.NewCaseTesterMock()
		ct.On("Test", &cases.Case{}).Return(nil)
		// - tester
		ts := internal.NewTester(rd, ct)

		// act
		err := ts.Run()

		// assert
		require.NoError(t, err)
	})

	t.Run("case 2: error - error reading a case", func(t *testing.T) {
		// arrange
		// - reader: mock
		rd := cases.NewReaderMock()
		rd.On("Read").Return(cases.Case{}, cases.ErrMalformedJSON)
		// - casetester: mock
		ct := internal.NewCaseTesterMock()
		// - tester
		ts := internal.NewTester(rd, ct)

		// act
		err := ts.Run()

		// assert
		require.Error(t, err)
		require.EqualError(t, err, cases.ErrMalformedJSON.Error())
	})
	
	t.Run("case 3: error - error testing a case", func(t *testing.T) {
		// arrange
		// - reader: mock
		rd := cases.NewReaderMock()
		rd.On("Read").Return(cases.Case{}, nil)
		// - casetester: mock
		ct := internal.NewCaseTesterMock()
		ct.On("Test", &cases.Case{}).Return(internal.ErrTesterReporter)
		// - tester
		ts := internal.NewTester(rd, ct)

		// act
		err := ts.Run()

		// assert
		require.Error(t, err)
		require.EqualError(t, err, internal.ErrTesterReporter.Error())
	})
}