package cmd

import (
	"context"
	"github.com/spf13/cobra"
)

// VerifyCmd represents the Verify command
var VerifyCmd = &cobra.Command{
	Use:   "verify",
	Run:   RunVerify,
	Short: "Verify your project's `glice.yaml` file",
	Long:  "Verify your project's `glice.yaml` file by loading it",
}

func init() {
	rootCmd.AddCommand(VerifyCmd)
}

func RunVerify(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	Notef("\n")
	LoadingProfileFile(ctx)
	Notef("\n\n")
}
