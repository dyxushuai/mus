//wrapped error type

package utils

import (
	"github.com/dropbox/godropbox/errors"
)

func NewError(format string, a ...interface{}) errors.DropboxError {
	return errors.Newf(format, a...)
}

