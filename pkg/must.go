package glice

import (
	"io"
)

func MustClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		Warnf("Unable to close file: %s", err.Error())
	}
}
