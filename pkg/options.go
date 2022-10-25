package glice

import (
	"fmt"
)

// Options provides a place to store command line options
// IMPORTANT: If any arrays, slices or pointers are used
// here be sure to update Clone() below.
type Options struct {
	VerbosityLevel int
	DirectOnly     bool
	LogOutput      bool
	NoCache        bool
	CaptureLicense bool
	LogFilepath    string
	SourceDir      string
	CacheFilepath  string
}

var options = &Options{}

// GetOptions returns a pointer to the Options object
// with global lifetime
func GetOptions() *Options {
	return options
}

// SetOptions sets options to the passed in Options object
func SetOptions(o *Options) {
	err := o.setLogging()
	if err != nil {
		Failf(ExitCannotSetOptions,
			"Unable to set options: %s",
			err.Error())
	}
	options = o
}

// Clone returns a copy of the receiver object
func (o *Options) Clone() *Options {
	// This is a shallow clone. If any slices or other pointers
	// get added to Options it may need to be updated.
	_o := *o
	return &_o
}

// IsLogging returns true when the user has either requested
// that output be logged, or set the log filepath to a value.
func (o *Options) IsLogging() (result bool) {
	if o.LogOutput {
		result = true
		goto end
	}
	if o.LogFilepath != "" {
		result = true
		goto end
	}
end:
	return result
}

// DiscardOutput returns true when the user has set the
// verbosity level of either requested
// that output be logged, or set the log filepath to a value.
func (o *Options) DiscardOutput() (result bool) {
	if o.LogOutput {
		result = true
		goto end
	}
	if o.LogFilepath != "" {
		result = true
		goto end
	}
end:
	return result
}

func (o *Options) setLogging() (err error) {

	if !o.IsLogging() {
		goto end
	}
	o.LogOutput = true

	if o.LogFilepath == "" {
		o.LogFilepath = LogFilepath()
	}

	err = InitializeLogging(o.LogFilepath)
	if err != nil {
		err = fmt.Errorf("unable to set logging; %w", err)
		goto end
	}
end:
	return err
}
