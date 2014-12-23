//wrapped error type

package manager

import (
	"time"
	"fmt"
)

type log interface {
	log() string
	print()
}

type logType struct {
	etime time.Time //error timestamp
	s string //error content
}

func (self *logType) log() string {
	return self.s
}

func (self *logType) print() {
	fmt.Println(fmt.Sprintf("[log] %s: %s", self.etime.Format("2006-01-02 15:04:05.999999999"), self.s))
}

func newLog(format string, a ...interface{}) log {
	var logStr string
	logStr = fmt.Sprintf(format, a...)
	return &logType{
		etime: time.Now(),
		s: logStr,
	}
}

