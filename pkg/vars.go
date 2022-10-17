package glice

import "regexp"

// exists is used as a no-memory value for maps used only to check
// if the key exists or not.
// See https://stackoverflow.com/questions/59089869/memory-usage-nil-interface-vs-struct
type exists struct{}

var (
	regexStripScheme = regexp.MustCompile("^https?://(.+)$")

	DefaultAllowedLicenses = []string{
		"Apache-2.0",
		"BSD-2-Clause",
		"BSD-3-Clause",
		"MIT",
		"MPL-2.0",
	}

	validFormats = map[string]exists{
		"table": {},
		"json":  {},
		"csv":   {},
	}

	// validOutputs to print to
	validOutputs = map[string]exists{
		"stdout": {},
		"file":   {},
	}
)
