package cmd

import (
	glice "github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

// overridesCmd represents the overrides command
var overridesCmd = &cobra.Command{
	Use:   "overrides",
	Run:   glice.RunOverrides,
	Short: "Generate overrides.yaml for editing",
	Long:  `Generate overrides.yaml for manual copying into glice.yaml and then manual editing to by the user`,
}

func init() {
	rootCmd.AddCommand(overridesCmd)
	addTTLFlag(overridesCmd)
}
