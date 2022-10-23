package cmd

import (
	"context"
	"fmt"
	glice "github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

const (
	DefaultReportFilename = "glice-dependencies-report.txt"
	DefaultReportFormat   = glice.TableFormat
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Run:   RunReportSave,
	Short: fmt.Sprintf("Save a report to file in %s, %s or %s format", glice.TableFormat, glice.JSONFormat, glice.CSVFormat),
}

func init() {
	reportCmd.AddCommand(saveCmd)
	saveCmd.Flags().String("filename", DefaultReportFilename, "File to save the license report to")
	saveCmd.Flags().String("format", string(DefaultReportFormat), fmt.Sprintf("Format in which to save report: %s, %s or %s", glice.TableFormat, glice.JSONFormat, glice.CSVFormat))
	saveCmd.MarkFlagsMutuallyExclusive("filename", "format")
}

//goland:noinspection GoUnusedParameter
func RunReportSave(cmd *cobra.Command, args []string) {
	NoteBegin()
	Notef("\nGenerating report to save")
	ctx := context.Background()
	adapter := GetReportWriterAdapter(cmd)
	adapter.SetDependencies(ScanDependencies(ctx))
	WriteReport(adapter)
	Notef("\nReport generated")
	NoteEnd()
}
