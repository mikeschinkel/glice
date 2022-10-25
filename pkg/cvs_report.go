package glice

import (
	"encoding/csv"
	"fmt"
)

var _ ReportWriterAdapter = (*CSVReport)(nil)

type CSVReport struct {
	*Report
}

func (cr *CSVReport) FileExtension() FileExtension {
	return CSVExtension
}

func (cr *CSVReport) WriteReport() error {
	writer := csv.NewWriter(cr.Writer)
	defer writer.Flush()
	err := writer.Write(cr.Dependencies.GetReportHeader())
	if err != nil {
		err = fmt.Errorf("unable to write report header; %w", err)
		goto end
	}
	for _, dep := range cr.Dependencies {
		err = writer.Write(dep.GetReportRow())
		if err != nil {
			err = fmt.Errorf("unable to write report row for %s; %w",
				dep.Import,
				err)
			break
		}
	}
end:
	return err
}
