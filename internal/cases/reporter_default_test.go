package cases_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/LNMMusic/tester/internal/cases"

	"github.com/stretchr/testify/require"
)

// Tests for ReporterDefault Report
func TestReporterDefault_Report(t *testing.T) {
	t.Run("case 1 - success report", func(t *testing.T) {
		// arrange
		rp := cases.NewReporterDefault()

		// act
		w := &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(strings.NewReader(
				`{"message":"success","data":{"string":"string","int":1,"float_zero":1,"float":1.1,"bool":true}}`,
			)),
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
		}
		c := &cases.Case{
			Name: "case 1",
			Response: cases.Response{
				Code: 200,
				Body: map[string]any{
					"message": "success",
					"data": map[string]any{
						"string": "string",
						"int": 1.0,
						"float_zero": 1.0,
						"float": 1.1,
						"bool": true,
					},
				},
				Header: http.Header{
					"Content-Type": {"application/json"},
				},
			},
		}
		err := rp.Report(c, w)

		// assert
		require.NoError(t, err)
	})

	t.Run("case 2 - failed report - code", func(t *testing.T) {
		// arrange
		rp := cases.NewReporterDefault()

		// act
		w := &http.Response{
			StatusCode: 400,
			Body: io.NopCloser(strings.NewReader(
				`{"message":"success","data":{"string":"string","int":1,"float_zero":1,"float":1.1,"bool":true}}`,
			)),
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
		}
		c := &cases.Case{
			Name: "case 2",
			Response: cases.Response{
				Code: 200,
				Body: map[string]any{
					"message": "success",
					"data": map[string]any{
						"string": "string",
						"int": 1.0,
						"float_zero": 1.0,
						"float": 1.1,
						"bool": true,
					},
				},
				Header: http.Header{
					"Content-Type": {"application/json"},
				},
			},
		}
		err := rp.Report(c, w)

		// assert
		require.NoError(t, err)
	})

	t.Run("case 3 - failed report - body", func(t *testing.T) {
		// arrange
		rp := cases.NewReporterDefault()

		// act
		w := &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(strings.NewReader(
				`{"message":"success","data":{"string":"string","int":1,"float_zero":1,"float":1.1,"bool":true}}`,
			)),
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
		}
		c := &cases.Case{
			Name: "case 3",
			Response: cases.Response{
				Code: 200,
				Body: map[string]any{
					"message": "success",
					"data": map[string]any{
						"string": "string",
						"int": 1.0,
						"float_zero": 1.0,
						"float": 1.1,
						"bool": false,
					},
				},
				Header: http.Header{
					"Content-Type": {"application/json"},
				},
			},
		}
		err := rp.Report(c, w)

		// assert
		require.NoError(t, err)
	})

	t.Run("case 4 - failed report - header", func(t *testing.T) {
		// arrange
		rp := cases.NewReporterDefault()

		// act
		w := &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(strings.NewReader(
				`{"message":"success","data":{"string":"string","int":1,"float_zero":1,"float":1.1,"bool":true}}`,
			)),
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
		}
		c := &cases.Case{
			Name: "case 4",
			Response: cases.Response{
				Code: 200,
				Body: map[string]any{
					"message": "success",
					"data": map[string]any{
						"string": "string",
						"int": 1.0,
						"float_zero": 1.0,
						"float": 1.1,
						"bool": true,
					},
				},
				Header: http.Header{
					"Content-Type": {"application/json"},
					"Content-Length": {"0"},
				},
			},
		}
		err := rp.Report(c, w)

		// assert
		require.NoError(t, err)
	})

	t.Run("case 5 - error decode body", func(t *testing.T) {
		// arrange
		rp := cases.NewReporterDefault()

		// act
		w := &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(strings.NewReader(
				`malformed-{"message":"success","data":{"string":"string","int":1,"float_zero":1,"float":1.1,"bool":true}}`,
			)),
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
		}
		c := &cases.Case{
			Name: "case 5",
			Response: cases.Response{
				Code: 200,
				Body: map[string]any{
					"message": "success",
					"data": map[string]any{
						"string": "string",
						"int": 1.0,
						"float_zero": 1.0,
						"float": 1.1,
						"bool": true,
					},
				},
				Header: http.Header{
					"Content-Type": {"application/json"},
				},
			},
		}
		err := rp.Report(c, w)

		// assert
		require.Error(t, err)
		require.EqualError(t, err, "invalid character 'm' looking for beginning of value")
	})
}