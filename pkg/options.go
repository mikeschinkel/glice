package glice

import (
	"io"
	"log"
)

// Options provides a place to store command line options
// IMPORTANT: If any arrays, slices or pointers are used
// here be sure to update Clone() below.
type Options struct {
	LogVerbosely    bool
	IncludeIndirect bool
	LogOuput        bool
	NoCache         bool
	LogFilepath     string
	SourceDir       string
	CacheFilepath   string

	WriteFile            bool
	OutputFormat         string
	NoCaptureLicenseText bool
	OutputDestination    string
}

var options = &Options{}

// GetOptions returns a pointer to the Options object
// with global lifetime
func GetOptions() *Options {
	return options
}

// SetOptions sets options to the passed in Options object
func SetOptions(o *Options) {
	o.setLogging()
	options = o
}

// Clone returns a copy of the receiver object
func (o *Options) Clone() *Options {
	// This is a shallow clone. If any slices or other pointers
	// get added to Options it may need to be updated.
	_o := *o
	return &_o
}

func (o *Options) setLogging() {

	if !o.LogVerbosely && !o.LogOuput {
		log.SetOutput(io.Discard)
		log.SetFlags(0)
		goto end
	}

	o.LogOuput = true

	if o.LogFilepath != "" {
		goto end
	}

	o.LogFilepath = LogFilepath()
end:
}
