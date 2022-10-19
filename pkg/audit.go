package glice

import (
	"github.com/spf13/cobra"
	"os"
)

func RunAudit(cmd *cobra.Command, args []string) {
	var err error
	var deps Dependencies

	options := GetOptions()

	Notef("\nBeginning License Audit")
	yf, err := LoadYAMLFile(options.SourceDir)
	if err != nil {
		Failf(exitYAMLFileDoesNotExist,
			"Cannot run scan; %s",
			err.Error())

	}
	Notef("\nYAML file %s loaded", yf.Filepath)

	Notef("\nScanning dependencies...")
	deps, err = ScanDependencies(options)
	if err != nil {
		Failf(exitCannotParseDependencies,
			"Failed while parsing dependencies: %s",
			err.Error())
	}

	Notef("\nAuditing dependencies...")
	changes, ds := yf.AuditDependencies(deps)
	Notef("\nAudit complete.\n")

	if !changes.HasChanges() {
		Notef("\nNo chances detected")
	} else {
		Notef("\n")
		changes.Print()
	}

	Errorf("\n")
	if !ds.HasDisalloweds() {
		Notef("\nOnly allowed licenses detected")
		Errorf("\n")
	} else {
		Errorf("\nDisallowed licenses detected:")
		Errorf("\n")
		ds.LogPrint()
		Errorf("\n")
		Errorf("\nAudit FAILED!")
		Errorf("\n\n")
		os.Exit(exitAuditFoundDisallowedLicenses)
	}

	Notef("\nAudit completed successfully")
	Notef("\n\n")
}
