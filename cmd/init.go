package cmd

import (
	"github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Run:   glice.RunInit,
	Short: "Initialize a 'glice.yaml' command in your project's path",
	Long: `Initialize a 'glice.yaml' command in your project's path. ` +
		`'init' will scan the go.mod file for dependencies and write ` +
		`them to the YAML file which can then be hand-edited to add ` +
		`overrides. Optionally it can generate a cache-file to allow ` +
		`future invocations of the 'scan' command to assume data to ` +
		`be current if last scanned within a specifiable TTL (time-to-live.)`,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().String("cache-file", glice.CacheFilepath(), "Full filepath for cache file to create")
	initCmd.Flags().String("path", glice.SourceDir(""), "Directory path to your project's top-level go.mod file")
}
