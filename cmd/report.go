package cmd

import (
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a dependency report to screen or a file using a subcommand",
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
