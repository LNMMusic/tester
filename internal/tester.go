package internal

import "github.com/LNMMusic/tester/internal/cases"

// NewTester creates a new tester.
func NewTester(rd cases.Reader, ct CaseTester) (t *Tester) {
	t = &Tester{
		rd: rd,
		ct: ct,
	}
	return
}

// Tester is an struct that test an stream of cases.
type Tester struct {
	// rd is the reader of cases.
	rd cases.Reader
	// ct is the tester of cases.
	ct CaseTester
}

// Run test a stream of cases.
func (t *Tester) Run() (err error) {
	// read cases
	for {
		// read case
		var c cases.Case
		c, err = t.rd.Read()
		if err != nil {
			if err == cases.ErrEndOfLine {
				err = nil
				break
			}
			return
		}
		// test case
		err = t.ct.Test(&c)
		if err != nil {
			return
		}
	}

	return
}