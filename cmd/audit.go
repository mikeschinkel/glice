package cmd

import (
	"context"
	"os"

	"github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Run:   RunAudit,
	Short: "Audit your project's path for disallowed open-source licenses",
	Long: `Audit your project's path for Go-specific dependencies using disallowed open-source licenses ` +
		`while comparing with allowed licenses and dependency overrides in glice.yaml and only auditing ` +
		`those dependencies that have not been audited within a specifiable TTL (time-to-live) where` +
		`the default TLL is 24*60*60 seconds (1 day)`,
}

func init() {
	rootCmd.AddCommand(auditCmd)
	addTTLFlag(auditCmd)
	auditCmd.Flags().Bool("overrides", false, "Write an `overrides.yaml` file if any disallowed licenses are found.")
}

func RunAudit(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	glice.Notef("\n")
	glice.Notef("\nBeginning License Audit")
	deps := ScanningDependencies(ctx)
	pf := AuditingProjectDependencies(ctx, "run", deps)
	glice.Notef("\n\n")
	HandleChanges(ctx, pf)
	exceptions := HasDisalloweds(ctx, pf)
	GeneratingOverrides(ctx, cmd, pf, glice.WarnLevel)
	if exceptions {
		os.Exit(glice.ExitAuditFoundDisallowedLicenses)
	}
	glice.Notef("\n\n")
}
