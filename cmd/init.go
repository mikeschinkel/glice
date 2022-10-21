package cmd

import (
	"context"
	"github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Run:   RunInit,
	Short: "Initialize a 'glice.yaml' command in your project's path",
	Long: `Initialize a 'glice.yaml' command in your project's path. ` +
		`'init' will scan the go.mod file for dependencies and write ` +
		`them to the YAML file which can then be hand-edited to add ` +
		`overrides. Optionally it can generate a cache-file to allow ` +
		`future invocations of the 'audit' command to assume data to ` +
		`be current if last audited within a specifiable TTL (time-to-live.)`,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func RunInit(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	glice.Notef("\n")
	glice.Notef("\nInitializing %s for project", glice.AppName)
	pf := CreatingProjectFile(ctx)
	pf.Dependencies = ScanningDependencies(ctx)
	SavingProjectFile(ctx, pf)
	glice.Notef("\n\n")
}
