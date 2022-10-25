package cmd

import (
	"github.com/spf13/cobra"
)

// gitLabCmd represents the gitLab command
var gitLabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "Commands with sub-commands specific to GitLab",
}

func init() {
	rootCmd.AddCommand(gitLabCmd)
}
