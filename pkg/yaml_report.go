package glice

import (
	yaml "gopkg.in/yaml.v3"
)

var _ ReportWriterAdapter = (*YAMLReport)(nil)

type YAMLReport struct {
	*Report
}

func (jr *YAMLReport) FileExtension() FileExtension {
	return YAMLExtension
}

func (jr *YAMLReport) WriteReport() error {
	return yaml.NewEncoder(jr.Writer).Encode(jr.Dependencies)
}
