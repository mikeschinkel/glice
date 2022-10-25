package glice

import (
	"fmt"
)

type ReportWriterAdapter interface {
	DependenciesSetter
	FileExtensionGetter
	FilepathSetter
	FilepathGetter
	FormatGetter
	ReportWriter
	WriterSetter
}

// GetReportWriterAdapter returns an object that implements the ReportWriterAdapter interface
// based on the format passed in, or returns an error if the format is not valid
func GetReportWriterAdapter(format OutputFormat) (adapter ReportWriterAdapter, err error) {
	report := NewReport(format)
	switch format {
	case TableFormat:
		adapter = &TableReport{Report: report}

	case JSONFormat:
		adapter = &JSONReport{Report: report}

	case YAMLFormat:
		adapter = &YAMLReport{Report: report}

	case CSVFormat:
		adapter = &CSVReport{Report: report}

	default:
		err = fmt.Errorf("invalid report format '%s', must be one of: %s",
			format,
			ValidOutputFormatsOrString)

	}
	return adapter, err
}
