package glice

import (
	"github.com/spf13/cobra"
)

func RunInit(cmd *cobra.Command, args []string) {
	var err error

	options := GetOptions()

	yf := NewYAMLFile(options.SourceDir)
	if yf.Exists() {
		Failf(exitYAMLFileExistsCannotOverwrite,
			"Cannot overwrite existing YAML file %s.\nRename or delete file then re-run 'glice init'.",
			yf.Filepath)
	}
	Notef("\nCreating %s\n", yf.Filepath)

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
	Notef("Scanning dependencies...")

	yf.Dependencies, err = ScanDependencies(options)
	if err != nil {
		Failf(exitCannotParseDependencies,
			"Failed while parsing dependencies: %s",
			err.Error())
	}

	yf.Filepath = YAMLFilepath(options.SourceDir)
	err = yf.Init()
	if err != nil {
		Failf(exitCannotInitializeYAMLFile,
			"Failed to create YAML file %s: %s",
			options.SourceDir,
			err.Error())
	}
	Notef("\nYAML file created.\n")
}
