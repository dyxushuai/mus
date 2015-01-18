package utils

import (
	goLog "github.com/segmentio/go-log"
)


func Debug(err error) {
	if err != nil {
		goLog.Debug(err.Error())
	}
}

func Info(msg string, args ...interface {}) {
	if msg != nil {
		goLog.Info(msg, args...)
	}
}

