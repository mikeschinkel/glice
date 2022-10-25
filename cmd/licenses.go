package cmd

import (
	"context"
	"fmt"

	"github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

const DefaultLicensesPath = "licenses"

// licensesCmd represents the licenses command
var licensesCmd = &cobra.Command{
	Use:   "licenses",
	Run:   RunLicenses,
	Short: "Write licenses for each dependency to a file",
}

func init() {
	rootCmd.AddCommand(licensesCmd)
	saveCmd.Flags().String("path",
		DefaultLicensesPath,
		fmt.Sprintf("Directory path in which to write licenses. Can be relative to %s, or absolute.",
			glice.ProjectFilename))
}

//goland:noinspection GoUnusedParameter
func RunLicenses(cmd *cobra.Command, args []string) {
	dir := Flag(cmd, "path")
	ctx := context.Background()

	NoteBegin()
	Notef("\nWriting Licenses to %s", dir)
	deps := ScanDependencies(ctx)
	SavingLicenses(deps, dir)
	NoteEnd()
}

func SavingLicenses(deps glice.Dependencies, dir string) {
	Notef("\nSaving Licenses")
	err := deps.SaveLicenses(dir, func(dep *glice.Dependency, fp string) {
		Infof("\nSaving license for %s to %s", dep.Import, fp)
	})
	if err != nil {
		Failf(glice.ExitCannotSaveFile,
			"\nUnable to write licenses for individual files; %w",
			err)
	}
	Notef("\nLicenses saved")
}
