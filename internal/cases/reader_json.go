package cases

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrInvalidToken is the error returned when the token is invalid.
	ErrInvalidToken = errors.New("invalid token")
	// ErrMalformedJSON is the error returned when the JSON is malformed.
	ErrMalformedJSON = errors.New("malformed json")
)

// NewReaderJSON creates a new reader of test cases in JSON format.
func NewReaderJSON(decoder *json.Decoder, ch chan CaseErr) *ReaderJSON {
	return &ReaderJSON{
		decoder: decoder,
		ch:      ch,
	}
}

// CaseErr is a test case with an error.
type CaseErr struct {
	// Case is the test case.
	Case Case
	// Err is the error.
	Err error
}

// ReaderJSON is a reader of test cases in JSON format.
type ReaderJSON struct {
	// decoder is the JSON decoder to use.
	decoder *json.Decoder
	// ch is the channel of test cases.
	ch chan CaseErr
}

// Read reads the next test case.
func (r *ReaderJSON) Read() (c Case, err error) {
	// fetch the next test case
	ce, ok := <-r.ch
	if !ok {
		err = ErrEndOfLine
		return
	}
	if ce.Err != nil {
		err = ce.Err
		return
	}

	c = ce.Case
	return
}

// Stream is a concurrent reader of test cases
func (r *ReaderJSON) Stream() {
	// close the channel at the end
	defer close(r.ch)

	// read the opening bracket of the array
	_, err := r.decoder.Token()
	if err != nil {
		r.ch <- CaseErr{Err: fmt.Errorf("%w - %s", ErrInvalidToken, err.Error())}
		return
	}

	// read the test cases
	for r.decoder.More() {
		var c Case
		err = r.decoder.Decode(&c)
		if err != nil {
			r.ch <- CaseErr{Err: fmt.Errorf("%w - %s", ErrMalformedJSON, err.Error())}
			return
		}

		r.ch <- CaseErr{Case: c}
	}
}