package cmd

import (
	"fmt"
	glice "github.com/ribice/glice/v3/pkg"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Run:   RunVersion,
	Short: "Display the current version of Glice",
	Long:  `Display the current version of Glice in 'v{Major}.{Minor}.{Patch}[-{Prerelease}" format.`,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

//goland:noinspection GoUnusedParameter
func RunVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("%s v%s\n", glice.CLIName, glice.AppVersion)
}
