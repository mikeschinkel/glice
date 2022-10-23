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

//goland:noinspection GoUnusedParameter
func RunOverrides(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	NoteBegin()
	Notef("\nGenerating Overrides file")
	deps := ScanDependencies(ctx)
	pf := AuditingProjectDependencies(ctx, deps)
	NoteEnd()
	GeneratingOverrides(ctx, cmd, pf, glice.ErrorLevel)
}
