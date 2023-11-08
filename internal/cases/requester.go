package cases

import "net/http"

// Requester is the interface for cases's requests
type Requester interface {
	// Do makes the request.
	Do(c *Case) (resp *http.Response, err error)
}