package cmd

import (
	"context"
	"fmt"

	"github.com/ribice/glice/v3/pkg"
	"github.com/ribice/glice/v3/pkg/gllicscan"
	"github.com/spf13/cobra"
)

// amendCmd represents the amend command
var amendCmd = &cobra.Command{
	Use:   "amend",
	Run:   RunReportAmend,
	Short: fmt.Sprintf("Amend a report to file in %s, %s or %s format", glice.TableFormat, glice.JSONFormat, glice.CSVFormat),
}

func init() {
	gitLabCmd.AddCommand(amendCmd)
	amendCmd.Flags().String("path", glice.SourceDir(), "Directory path were license scanning report can be found")
	amendCmd.Flags().String("filename", gllicscan.ReportFilename, "JSON file to amend containing results of license scan on GitLab")
}

//goland:noinspection GoUnusedParameter
func RunReportAmend(cmd *cobra.Command, args []string) {
	NoteBegin()
	Notef("\nAmending license scanning report")
	ctx := context.Background()
	toAmend := LoadJSONReportFromGitLab(cmd)
	deps := ScanDependencies(ctx)
	pf := LoadProjectFile(ctx)
	gr := glice.NewGitLabReport(deps, pf.Overrides)
	toAmend.Amend(gr)
	SaveJSONReportForGitLab(ctx, toAmend)
	Notef("\nReport amended")
	NoteEnd()
}
