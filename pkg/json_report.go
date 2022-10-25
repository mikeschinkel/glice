package glice

import (
	"encoding/json"
)

var _ ReportWriterAdapter = (*JSONReport)(nil)

type JSONReport struct {
	*Report
}

func (jr *JSONReport) FileExtension() FileExtension {
	return JSONExtension
}

func (jr *JSONReport) WriteReport() error {
	enc := json.NewEncoder(jr.Writer)
	enc.SetIndent("", "\t")
	return enc.Encode(jr.Dependencies)
}
