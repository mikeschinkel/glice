package cmd

import (
	"context"
	"fmt"
	glice "github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

// VerifyCmd represents the Verify command
var VerifyCmd = &cobra.Command{
	Use:   "verify",
	Run:   RunVerify,
	Short: fmt.Sprintf("Verify your project's `%s` file", glice.ProjectFilename),
	Long:  fmt.Sprintf("Verify your project's `%s` file by loading it", glice.ProjectFilename),
}

func init() {
	rootCmd.AddCommand(VerifyCmd)
}

//goland:noinspection GoUnusedParameter
func RunVerify(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	NoteBegin()
	LoadingProfileFile(ctx)
	NoteEnd()
}
