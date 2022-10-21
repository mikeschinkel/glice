package cmd

import (
	"context"
	"github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

// overridesCmd represents the overrides command
var overridesCmd = &cobra.Command{
	Use:   "overrides",
	Run:   RunOverrides,
	Short: "Generate overrides.yaml for editing",
	Long:  `Generate overrides.yaml for manual copying into glice.yaml and then manual editing to by the user`,
}

func init() {
	rootCmd.AddCommand(overridesCmd)
	addTTLFlag(overridesCmd)
}

func RunOverrides(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	glice.Notef("\n")
	glice.Notef("\nGenerating Overrides file")
	deps := ScanningDependencies(ctx)
	pf := AuditingProjectDependencies(ctx, "generate", deps)
	glice.Notef("\n\n")
	GeneratingOverrides(ctx, cmd, pf, glice.ErrorLevel)
}
