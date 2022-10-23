package glice

import (
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var _ ReportWriterAdapter = (*TableReport)(nil)

type TableReport struct {
	*Report
}

func (tr *TableReport) FileExtension() FileExtension {
	return TableExtension
}

func (tr *TableReport) WriteReport() error {
	tw := tablewriter.NewWriter(tr.Writer)
	tw.SetHeader(tr.Dependencies.GetReportHeader())
	for _, dep := range tr.Dependencies {
		ru := dep.RepoURL
		if tr.Filepath == "" {
			ru = color.BlueString(dep.RepoURL)
		}
		tw.Append([]string{dep.Import, ru, dep.LicenseID, dep.Added})
	}
	tw.Render()
	return nil
}
