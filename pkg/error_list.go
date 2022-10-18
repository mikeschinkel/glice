package glice

import (
	"fmt"
	"log"
)

var _ error = (*ErrorList)(nil)

type ErrorList []error

// HasErrors returns true if ErrorList has one or more errors
func (el ErrorList) HasErrors() bool {
	return len(el) > 0
}

// Error returns all errors as a single string
func (el ErrorList) Error() (s string) {
	for _, err := range el {
		s = fmt.Sprintf("%s | %s", s, err.Error())
	}
	return fmt.Sprintf("[%s ]", s[2:])
}

// LogPrint outputs all errors in list individually
func (el ErrorList) LogPrint() {
	el.LogPrintWithHeader("")
}

// LogPrintWithHeader outputs all errors in list individually but with header
func (el ErrorList) LogPrintWithHeader(header string) {
	LogPrintFunc(func() {
		if header != "" {
			log.Printf("\n%s\n", header)
		}
		for _, err := range el {
			log.Printf("- %s\n", err.Error())
		}
	})
}
