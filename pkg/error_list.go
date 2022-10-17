package glice

import (
	"fmt"
	"log"
)

var _ error = (*ErrorList)(nil)

type ErrorList []error

func (el ErrorList) Error() (s string) {
	for _, err := range el {
		s = fmt.Sprintf("%s | %s", s, err.Error())
	}
	return fmt.Sprintf("[%s ]", s[2:])
}

func (el ErrorList) LogPrint() {
	for _, err := range el {
		log.Print(err.Error())
	}
}
