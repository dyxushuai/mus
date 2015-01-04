//wrapped error type

package manager

import (
	goLog "github.com/segmentio/go-log"
)


type Verbose bool

func (self *Verbose) withVerboseDo(fn func()) {
	if self {
		fn()
	}
}

func (self *Verbose) Debug(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = goLog.Debug(msg, args...)
	})
	return
}

// Info log.
func (self *Verbose) Info(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = goLog.Info(msg, args...)
	})
	return
}

// Notice log.
func (self *Verbose) Notice(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = goLog.Notice(msg, args...)
	})
	return
}

// Warning log.
func (self *Verbose) Warning(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = goLog.Warning(msg, args...)
	})
	return
}

// Error log.
func (self *Verbose) Error(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = goLog.Error(msg, args...)
	})
	return
}

// Critical log.
func (self *Verbose) Critical(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = goLog.Critical(msg, args...)
	})
	return
}

// Alert log.
func (self *Verbose) Alert(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = goLog.Alert(msg, args...)
	})
	return
}

// Emergency log.
func (self *Verbose) Emergency(msg string, args ...interface{}) (err error) {
	self.withVerboseDo(func() {
		err = goLog.Emergency(msg, args...)
	})
	return
}


