package glice

import (
	"context"
	"github.com/spf13/cobra"
)

func RunOverrides(cmd *cobra.Command, args []string) {
	var err error
	var deps Dependencies

	ctx := context.Background()

	options := GetOptions()

	Notef("\nCreating %s", OverridesFilename)

	of := NewOverridesFile(options.SourceDir)
	if of.Exists() {
		Failf(exitFileExistsCannotOverwrite,
			"\nCannot overwrite existing file %s.\nRename or delete file then re-run 'glice overrides'.",
			of.Filepath)
	}

	Notef("\nLoading `%s`", ProjectFilename)
	pf, err := LoadProjectFile(options.SourceDir)
	if err != nil {
		Failf(exitFileDoesNotExist,
			"Cannot generate overrides; %s",
			err.Error())
	}
	Notef("\nLoaded `%s`", pf.Filepath)

	Notef("\nScanning dependencies...")
	deps, err = ScanDependencies(ctx, options)
	if err != nil {
		Failf(exitCannotParseDependencies,
			"Failed while scanning dependencies: %s",
			err.Error())
	}

	Notef("\nAuditing dependencies...")
	_, disalloweds := pf.AuditDependencies(deps)
	Notef("\nAudit complete.")

	Notef("\nExtractings editors and overrides")
	of.Editors, of.Overrides = disalloweds.ToEditorsAndOverrides(ctx)

	Notef("\nCreating %s", of.Filepath)
	err = of.Create()
	if err != nil {
		Failf(exitCannotCreateFile,
			"Failed to create file %s: %s",
			of.Filepath,
			err.Error())
	}

	Notef("\nOverrides files created.\n")
	Notef("\n")

}
