package utils

import (
	goLog "github.com/segmentio/go-log"
	"os"
)


var log = goLog.New(os.Stderr, goLog.DEBUG, "")

func Debug(err error) {
	if err != nil {
		log.Debug(err.Error())
	}
}

func Info(msg string, args ...interface {}) {
	if msg != "" {
		log.Info(msg, args...)
	}
}

