package glice

import (
	"fmt"
	"github.com/spf13/cobra"
)

func RunInit(cmd *cobra.Command, args []string) {
	var err error

	options := GetOptions()

	yf := NewYAMLFile(options.SourceDir)
	if yf.Exists() {
		LogAndExit(exitYAMLFileExistsCannotOverwrite,
			"Cannot overwrite existing YAML file %s.\nRename or delete file then re-run 'glice init'.",
			yf.Filepath)
	}
	fmt.Printf("\nCreating %s\n", yf.Filepath)

	yf.Editors = Editors{
		{Name: "Name goes here", Email: "email-alias@singlestore.com"},
	}
	yf.Overrides = Overrides{
		{
			DependencyImport: "https://github.com/example.com/sample",
			LicenseID:        "MIT",
			VerifiedBy:       "*email-alias",
			LastVerified:     Timestamp()[:10],
			Notes:            "This is a sample override added by 'glice init' command",
		},
	}
	fmt.Print("Scanning dependencies...")

	yf.Dependencies, err = ScanDependencies(options)
	if err != nil {
		LogAndExit(exitCannotParseDependencies,
			"Failed while parsing dependencies: %s",
			err.Error())
	}

	yf.Filepath = YAMLFilepath(options.SourceDir)
	err = yf.Init()
	if err != nil {
		LogAndExit(exitCannotInitializeYAMLFile,
			"Failed to create YAML file %s: %s",
			options.SourceDir,
			err.Error())
	}
	fmt.Println("\nYAML file created.")
}
