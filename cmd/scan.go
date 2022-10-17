/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Initialize a 'glice.yaml' command in your project's path",
	Long: `Initialize a 'glice.yaml' command in your project's path. ` +
		`'init' will scan the go.mod file for dependencies and write ` +
		`them to the YAML file which can then be hand-edited to add ` +
		`overrides. Optionally it can generate a cache-file to allow ` +
		`future invocations of the 'scan' command to assume data to ` +
		`be current if last scanned within a specifiable TTL (time-to-live.)` +
		`the default TLL is 24*60*60 seconds (1 day)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scan called")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
