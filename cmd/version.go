package cmd

import (
	glice "github.com/ribice/glice/v3/pkg"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Run:   glice.RunVersion,
	Short: "Display the current version of Glice",
	Long:  `Display the current version of Glice in 'v{Major}.{Minor}.{Patch}[-{Prerelease}" format.`,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
