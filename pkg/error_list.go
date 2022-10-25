package glice

import (
	"fmt"
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
	level := 3
	LogPrintFunc(level, func() {
		for _, err := range el {
			LogPrintf(level, "%s: - %s\n", LogLevels[level], err.Error())
		}
	})
}
