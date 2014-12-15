//wrapped error type

package manager

import (
	"time"
)



type errorType struct {
	etime time.Time //error timestamp
	s string //error content
}

func (err *errorType) Error() string {
	return err.etime.String() + ":" + err.s
}

func newError(s string) error {
	return &errorType{
		etime: time.Now(),
		s: s,
	}
}
