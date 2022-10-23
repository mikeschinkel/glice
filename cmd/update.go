package cmd

import (
	"context"
	"fmt"
	glice "github.com/ribice/glice/v3/pkg"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Run:   RunUpdate,
	Short: fmt.Sprintf("Update your project's `%s` file to include newly added dependencies", glice.ProjectFilename),
	Long:  fmt.Sprintf("Update your project's `%s` file to include newly added dependencies found when scanning go.mod.", glice.ProjectFilename),
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

//goland:noinspection GoUnusedParameter
func RunUpdate(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	NoteBegin()
	Notef("\nUpdating project file '%s'", glice.ProjectFilename)
	deps := ScanDependencies(ctx)
	pf := AuditingProjectDependencies(ctx, deps)
	SavingProjectFile(ctx, pf)
	NoteEnd()
}
