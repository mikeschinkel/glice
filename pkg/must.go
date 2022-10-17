package glice

import "io"

func MustClose(c io.Closer) {
	_ = c.Close()
}
