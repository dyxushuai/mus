//wrapped error type

package manager

import (
	"time"
	"fmt"
)


type errorType struct {
	etime time.Time //error timestamp
	s string //error content
}

func (err *errorType) Error() string {
	return "[Err]" + err.etime.Format("2006-01-02 15:04:05") + ":" + err.s
}

func newError(format string, a ...interface{}) error {

	var errStr string
	errStr = fmt.Sprintf(format, a...)
	return &errorType{
		etime: time.Now(),
		s: errStr,
	}
}

