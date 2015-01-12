//wrapped error type

package manager

import (

	"github.com/dropbox/godropbox/errors"
)


func newError(format string, a ...interface{}) errors.DropboxError {
	return errors.Newf(format, a...)
}

