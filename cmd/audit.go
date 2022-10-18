package cmd

import (
	"github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Run:   glice.RunAudit,
	Short: "Audit your project's path for disallowed open-source licenses",
	Long: `Audit your project's path for Go-specific dependencies using disallowed open-source licenses ` +
		`while comparing with allowed licenses and dependency overrides in glice.yaml and only auditing ` +
		`those dependencies that have not been audited within a specifiable TTL (time-to-live) where` +
		`the default TLL is 24*60*60 seconds (1 day)`,
}

func init() {
	rootCmd.AddCommand(auditCmd)
	initCmd.Flags().Int("ttl", 24*60*60, "Time-to-Live for data in the cache file allowing recently audited dependencies to be skipped")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// auditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// auditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
