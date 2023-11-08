package cases

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// NewRequesterDefault creates a new default requester.
func NewRequesterDefault(serverAddr string, client *http.Client) *RequesterDefault {
	// default config
	defaultServerAddr := "http://localhost:8080"
	defaultClient := &http.Client{}
	if serverAddr != "" {
		defaultServerAddr = serverAddr
	}
	if client != nil {
		defaultClient = client
	}

	return &RequesterDefault{
		serverAddr: defaultServerAddr,
		client: defaultClient,
	}
}

// RequesterDefault is the default requester.
type RequesterDefault struct {
	// serverAddr is the address of the server to test.
	serverAddr string
	// client is the client to use to make requests.
	client *http.Client
}

// Do makes the request.
func (r *RequesterDefault) Do(c *Case) (resp *http.Response, err error) {
	// request elements
	// - method
	method := c.Request.Method
	// - url
	url := r.serverAddr + c.Request.Path
	// - body
	var body io.Reader
	if c.Request.Body != nil {
		var b []byte
		b, err = json.Marshal(c.Request.Body)
		if err != nil {
			return
		}

		body = bytes.NewReader(b)
	}

	// request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	// - query
	q := req.URL.Query()
	for k, v := range c.Request.Query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	// - headers
	for k, v := range c.Request.Header {
		req.Header.Set(k, v[0])
	}

	// send
	resp, err = r.client.Do(req)
	if err != nil {
		return
	}

	return
}