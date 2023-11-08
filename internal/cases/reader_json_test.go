package cases_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/LNMMusic/tester/internal/cases"

	"github.com/stretchr/testify/require"
)

// Tests for NewReaderJSON Stream
func TestReaderJSON_Stream(t *testing.T) {
	t.Run("case 1 - success to read some cases", func(t *testing.T) {
		// arrange
		dc := json.NewDecoder(strings.NewReader(
			`[
				{"case_name":"case 1","database":{"set_up":["INSERT 1","INSERT 2"],"tear_down":["DELETE","RESET"]},"request":{"method":"GET","path":"/","query":{"key":"value"},"body":{"key":1},"header":{"Content-Type":["application/json"]}},"response":{"code":200,"body":{"key":1},"header":{"Content-Type":["application/json"]}}},
				{"case_name":"case 2","request":{"method":"POST"},"response":{"code":200}}
			]`,
		))
		ch := make(chan cases.CaseErr)
		rd := cases.NewReaderJSON(dc, ch)

		// act
		go rd.Stream()
		c1 := <-ch
		c2 := <-ch
		_, ok := <-ch

		// assert
		require.Equal(t, cases.CaseErr{
			Case: cases.Case{
				Name: "case 1",
				Database: cases.Database{
					SetUp: []string{
						"INSERT 1",
						"INSERT 2",
					},
					TearDown: []string{
						"DELETE",
						"RESET",
					},
				},
				Request: cases.Request{
					Method: "GET",
					Path:   "/",
					Query: map[string]string{
						"key": "value",
					},
					Body: map[string]any{
						"key": 1.0,
					},
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				},
				Response: cases.Response{
					Code: 200,
					Body: map[string]any{
						"key": 1.0,
					},
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				},
			},
			Err: nil,
		}, c1) 
		require.Equal(t, cases.CaseErr{
			Case: cases.Case{
				Name: "case 2",
				Request: cases.Request{
					Method: "POST",
				},
				Response: cases.Response{
					Code: 200,
				},
			},
			Err: nil,
		}, c2)
		require.False(t, ok)
	})

	t.Run("case 2 - success to read empty cases", func(t *testing.T) {
		// arrange
		dc := json.NewDecoder(strings.NewReader("[]"))
		ch := make(chan cases.CaseErr)
		rd := cases.NewReaderJSON(dc, ch)

		// act
		go rd.Stream()
		_, ok := <-ch

		// assert
		require.False(t, ok)
	})

	t.Run("case 3 - error wrong token", func(t *testing.T) {
		// arrange
		dc := json.NewDecoder(strings.NewReader("}"))
		ch := make(chan cases.CaseErr)
		rd := cases.NewReaderJSON(dc, ch)

		// act
		go rd.Stream()
		c1 := <-ch
		_, ok := <-ch

		// assert
		require.Equal(t, cases.Case{}, c1.Case)
		require.ErrorIs(t, c1.Err, cases.ErrInvalidToken)
		require.EqualError(t, c1.Err, fmt.Sprintf("%s - %s", cases.ErrInvalidToken.Error(), "invalid character '}' looking for beginning of value"))
		require.False(t, ok)
	})

	t.Run("case 4 - malformed object json", func(t *testing.T) {
		// arrange
		dc := json.NewDecoder(strings.NewReader(
			`[
				invalid json
			]`,
		))
		ch := make(chan cases.CaseErr)
		rd := cases.NewReaderJSON(dc, ch)

		// act
		go rd.Stream()
		c1 := <-ch
		_, ok := <-ch

		// assert
		require.Equal(t, cases.Case{}, c1.Case)
		require.ErrorIs(t, c1.Err, cases.ErrMalformedJSON)
		require.EqualError(t, c1.Err, fmt.Sprintf("%s - %s", cases.ErrMalformedJSON.Error(), "invalid character 'i' looking for beginning of value"))
		require.False(t, ok)
	})
}

func TestReaderJSON_Read(t *testing.T) {
	t.Run("case 1 - success to read some cases", func(t *testing.T) {
		// arrange
		ch := make(chan cases.CaseErr)
		go func () {
			ch <- cases.CaseErr{
				Case: cases.Case{
					Name: "case 1",
					Request: cases.Request{
						Method: "GET",
					},
					Response: cases.Response{
						Code: 200,
					},
				},
				Err: nil,
			}
			ch <- cases.CaseErr{
				Case: cases.Case{
					Name: "case 2",
					Request: cases.Request{
						Method: "POST",
					},
					Response: cases.Response{
						Code: 200,
					},
				},
				Err: nil,
			}
			close(ch)
		}()
		rd := cases.NewReaderJSON(nil, ch)

		// act
		c1, err1 := rd.Read()
		c2, err2 := rd.Read()
		c3, err3 := rd.Read()

		// assert
		require.NoError(t, err1)
		require.Equal(t, cases.Case{
			Name: "case 1",
			Request: cases.Request{
				Method: "GET",
			},
			Response: cases.Response{
				Code: 200,
			},
		}, c1)
		require.NoError(t, err2)
		require.Equal(t, cases.Case{
			Name: "case 2",
			Request: cases.Request{
				Method: "POST",
			},
			Response: cases.Response{
				Code: 200,
			},
		}, c2)
		require.ErrorIs(t, err3, cases.ErrEndOfLine)
		require.EqualError(t, err3, cases.ErrEndOfLine.Error())
		require.Equal(t, cases.Case{}, c3)
	})

	t.Run("case 2 - success to read empty cases", func(t *testing.T) {
		// arrange
		ch := make(chan cases.CaseErr)
		go func () {
			close(ch)
		}()
		rd := cases.NewReaderJSON(nil, ch)

		// act
		c1, err := rd.Read()

		// assert
		require.ErrorIs(t, err, cases.ErrEndOfLine)
		require.EqualError(t, err, cases.ErrEndOfLine.Error())
		require.Equal(t, cases.Case{}, c1)
	})

	t.Run("case 3 - error wrong token", func(t *testing.T) {
		// arrange
		ch := make(chan cases.CaseErr)
		go func () {
			ch <- cases.CaseErr{
				Case: cases.Case{},
				Err:  errors.New("wrong token"),
			}
			close(ch)
		}()
		rd := cases.NewReaderJSON(nil, ch)

		// act
		c1, err1 := rd.Read()
		c2, err2 := rd.Read()

		// assert
		require.EqualError(t, err1, "wrong token")
		require.Equal(t, cases.Case{}, c1)
		require.ErrorIs(t, err2, cases.ErrEndOfLine)
		require.EqualError(t, err2, cases.ErrEndOfLine.Error())
		require.Equal(t, cases.Case{}, c2)
	})
}