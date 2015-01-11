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

func (self *errorType) Error() string {
	return self.s
}

func (self *errorType) Print() {
	fmt.Println(fmt.Sprintf("[err] %s: %s", self.etime.Format("2006-01-02 15:04:05.999999999"), self.s))
}

func newError(format string, a ...interface{}) error {

	var errStr string
	errStr = fmt.Sprintf(format, a...)
	return &errorType{
		etime: time.Now(),
		s: errStr,
	}
}

