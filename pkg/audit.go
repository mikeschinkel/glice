package glice

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func RunAudit(cmd *cobra.Command, args []string) {
	var err error
	var deps Dependencies

	options := GetOptions()

	fmt.Println("\nAuditing...")
	yf, err := LoadYAMLFile(options.SourceDir)
	if err != nil {
		LogAndExit(exitYAMLFileDoesNotExist,
			"Cannot run scan; %s",
			err.Error())

	}
	fmt.Printf("YAML file %s loaded\n", yf.Filepath)

	fmt.Print("Scanning dependencies...")
	deps, err = ScanDependencies(options)
	if err != nil {
		LogAndExit(exitCannotParseDependencies,
			"Failed while parsing dependencies: %s",
			err.Error())
	}

	changes, el := yf.AuditDependencies(deps)
	if !changes.HasChanges() {
		fmt.Println("\nNo chances detected")
	} else {
		fmt.Println()
		changes.Print()
	}

	if !el.HasErrors() {
		fmt.Println("\nNo disallowed licenses detected")
	} else {
		el.LogPrintWithHeader("ERROR! Disallowed Licenses Detected:")
		os.Exit(exitAuditFoundDisallowedLicenses)
	}

	fmt.Println("\nAudit completed successfully")
}
