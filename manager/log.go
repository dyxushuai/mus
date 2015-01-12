//wrapped error type

package manager

import (
	goLog "github.com/segmentio/go-log"
	"os"
)


type Verbose bool


var logger *goLog.Logger

func init() {
	logger = goLog.New(os.Stderr, goLog.DEBUG, "")
}

func (self Verbose) withVerboseDo(fn func()) {
	if self == true {
		fn()
	}
}

func (self *Verbose) Debug(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = logger.Debug(msg, args...)
	})
	return
}

// Info log.
func (self *Verbose) Info(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = logger.Info(msg, args...)
	})
	return
}

// Notice log.
func (self *Verbose) Notice(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = logger.Notice(msg, args...)
	})
	return
}

// Warning log.
func (self *Verbose) Warning(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = logger.Warning(msg, args...)
	})
	return
}

// Error log.
func (self *Verbose) Error(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = logger.Error(msg, args...)
	})
	return
}

// Critical log.
func (self *Verbose) Critical(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = logger.Critical(msg, args...)
	})
	return
}

// Alert log.
func (self *Verbose) Alert(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = logger.Alert(msg, args...)
	})
	return
}

// Emergency log.
func (self *Verbose) Emergency(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = logger.Emergency(msg, args...)
	})
	return
}


