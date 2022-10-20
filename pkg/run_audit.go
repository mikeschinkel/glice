package glice

import (
	"context"
	"github.com/spf13/cobra"
	"os"
)

func RunAudit(cmd *cobra.Command, args []string) {
	var err error
	var deps Dependencies

	ctx := context.Background()
	options := GetOptions()

	Notef("\n")
	Notef("\nBeginning License Audit")
	pf, err := LoadProjectFile(options.SourceDir)
	if err != nil {
		Failf(exitFileDoesNotExist,
			"Cannot run audit; %s",
			err.Error())

	}
	Notef("\n%s loaded", pf.Filepath)

	Notef("\nScanning dependencies...")
	deps, err = ScanDependencies(ctx, options)
	if err != nil {
		Failf(exitCannotParseDependencies,
			"Failed while scanning dependencies: %s",
			err.Error())
	}

	Notef("\nAuditing dependencies...")
	changes, disalloweds := pf.AuditDependencies(deps)
	Notef("\nAudit complete.")
	Notef("\n\n")

	if !changes.HasChanges() {
		Notef("\nNo chances detected")
	} else {
		Notef("\n")
		changes.Print()
	}

	if len(disalloweds) == 0 {
		Notef("\n")
		Notef("\nOnly allowed licenses detected")
		Notef("\nAudit completed successfully\n")
		Notef("\n")
		goto end
	}

	Errorf("\n")
	Errorf("\nDisallowed licenses detected:")
	Errorf("\n")
	disalloweds.LogPrint()
	Errorf("\n")
	Errorf("\nAudit FAILED!")
	Errorf("\n\n")

	if ShouldGenerateOverrides(cmd) {
		GenerateOverrides(ctx, disalloweds)
	}
	os.Exit(exitAuditFoundDisallowedLicenses)

end:
}

func ShouldGenerateOverrides(cmd *cobra.Command) bool {
	return Flag(cmd, "overrides") == "true"
}

func GenerateOverrides(ctx context.Context, disalloweds Dependencies) {
	Notef("\n")
	Notef("\nCreating %s", OverridesFilename)

	of := NewOverridesFile(options.SourceDir)
	if of.Exists() {
		Warnf("\nCannot overwrite existing file %s.\nRename or delete file then re-run 'glice overrides'.",
			of.Filepath)
	}

	Notef("\nExtracting editors and overrides")
	of.Editors, of.Overrides = disalloweds.ToEditorsAndOverrides(ctx)

	Notef("\nCreating %s", of.Filepath)
	err := of.Create()
	if err != nil {
		Warnf("Failed to create file %s: %s",
			of.Filepath,
			err.Error())
	}
	Notef("\nOverrides files created.")
	Notef("\n\n")

}
