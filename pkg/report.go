package glice

import (
	"fmt"
	"io"
)

type OutputFormat string
type FileExtension string

const (
	TableFormat OutputFormat = "table"
	JSONFormat  OutputFormat = "json"
	YAMLFormat  OutputFormat = "yaml"
	CSVFormat   OutputFormat = "csv"

	TableExtension FileExtension = "txt"
	JSONExtension  FileExtension = "json"
	YAMLExtension  FileExtension = "yaml"
	CSVExtension   FileExtension = "csv"
)

type OutputFormats []OutputFormat

var ValidOutputFormatsFormat = fmt.Sprintf("'%s', '%s', '%s' %s '%s'", TableFormat, JSONFormat, YAMLFormat, "%s", CSVFormat)
var ValidOutputFormatsOrString = fmt.Sprintf(ValidOutputFormatsFormat, "or")

//goland:noinspection GoUnusedGlobalVariable
var ValidOutputFormatsAndString = fmt.Sprintf(ValidOutputFormatsFormat, "and")

type Report struct {
	io.Writer
	Format       OutputFormat
	Dependencies Dependencies
	Filepath     string
}

// NewReport returns a new base report containing the passed writer and dependencies
func NewReport(format OutputFormat) *Report {
	return &Report{
		Format: format,
	}
}

func (r *Report) SetWriter(w io.Writer) {
	r.Writer = w
}

func (r *Report) SetFilepath(fp string) {
	r.Filepath = fp
}

func (r *Report) GetFilepath() string {
	return r.Filepath
}

func (r *Report) GetFormat() OutputFormat {
	return r.Format
}

func (r *Report) SetDependencies(deps Dependencies) {
	r.Dependencies = deps
}
