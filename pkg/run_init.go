package glice

import (
	"context"
	"github.com/spf13/cobra"
)

func RunInit(cmd *cobra.Command, args []string) {
	var err error

	ctx := context.Background()

	options := GetOptions()

	pf := NewProjectFile(options.SourceDir)
	if pf.Exists() {
		Failf(exitFileExistsCannotOverwrite,
			"Cannot overwrite existing file %s.\nRename or delete file then re-run 'glice init'.",
			pf.Filepath)
	}
	Notef("\nCreating %s\n", pf.Filepath)

	pf.Editors = Editors{
		{Name: "Name goes here", Email: "email-alias@singlestore.com"},
	}
	pf.Overrides = Overrides{
		{
			DependencyImport: "https://github.com/example.com/sample",
			LicenseID:        "MIT",
			VerifiedBy:       "*email-alias",
			LastVerified:     Timestamp()[:10],
			Notes:            "This is a sample override added by 'glice init' command",
		},
	}

	Notef("Scanning dependencies...")
	pf.Dependencies, err = ScanDependencies(ctx, options)
	if err != nil {
		Failf(exitCannotParseDependencies,
			"Failed while scanning dependencies: %s",
			err.Error())
	}

	err = pf.Initialize()
	if err != nil {
		Failf(exitCannotCreateFile,
			"Failed to create file %s: %s",
			pf.Filepath,
			err.Error())
	}
	Notef("\n`%s` file created.\n", ProjectFilename)
}
