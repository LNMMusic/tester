package cases_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LNMMusic/tester/internal/cases"

	"github.com/stretchr/testify/require"
)

// Tests for RequesterDefault Do method
func TestRequesterDefault_Do(t *testing.T) {
	t.Run("case 1: success to make request", func(t *testing.T) {
		// arrange
		// - server: mock
		hd := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello world"))
		}
		sv := httptest.NewServer(http.HandlerFunc(hd))
		defer sv.Close()
		// - requester
		rq := cases.NewRequesterDefault(sv.URL, nil)

		// act
		c := &cases.Case{
			Request: cases.Request{
				Method: http.MethodGet,
				Path: "/",
				Query: map[string]string{
					"q1": "v1",
				},
				Body: map[string]interface{}{
					"b1": "v1",
				},
				Header: http.Header{
					"h1": []string{"v1"},
				},
			},
		}
		resp, err := rq.Do(c)

		// assert
		require.NoError(t, err)
		
		expectedCode := http.StatusOK
		expectedBody := "hello world"
		currentBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, expectedCode, resp.StatusCode)
		require.Equal(t, expectedBody, string(currentBody))
	})

	t.Run("case 2: success to make request - body nil", func(t *testing.T) {
		// arrange
		// - server: mock
		hd := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello world"))
		}
		sv := httptest.NewServer(http.HandlerFunc(hd))
		defer sv.Close()
		// - requester
		rq := cases.NewRequesterDefault(sv.URL, nil)

		// act
		c := &cases.Case{
			Request: cases.Request{
				Method: http.MethodGet,
				Path: "/",
				Query: map[string]string{
					"q1": "v1",
				},
				Body: nil,
				Header: http.Header{
					"h1": []string{"v1"},
				},
			},
		}
		resp, err := rq.Do(c)

		// assert
		require.NoError(t, err)
		
		expectedCode := http.StatusOK
		expectedBody := "hello world"
		currentBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, expectedCode, resp.StatusCode)
		require.Equal(t, expectedBody, string(currentBody))
	})

	t.Run("case 3: fail to make request - invalid url", func(t *testing.T) {
		// arrange
		// - server: mock
		// ...
		// - requester
		rq := cases.NewRequesterDefault("invalid", nil)

		// act
		c := &cases.Case{
			Request: cases.Request{
				Method: http.MethodGet,
				Path: "/",
				Query: map[string]string{
					"q1": "v1",
				},
				Body: map[string]interface{}{
					"b1": "v1",
				},
				Header: http.Header{
					"h1": []string{"v1"},
				},
			},
		}
		resp, err := rq.Do(c)

		// assert
		require.Error(t, err)
		require.EqualError(t, err, "Get \"invalid/?q1=v1\": unsupported protocol scheme \"\"")
		require.Nil(t, resp)
	})
}