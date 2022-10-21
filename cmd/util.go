package cmd

import (
	"context"
	"github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

func ExtractingEditorsAndOOverrides(ctx context.Context, pf *glice.ProjectFile, of *glice.OverridesFile) {
	glice.Notef("\nExtracting editors and overrides")
	of.Editors, of.Overrides = pf.Disalloweds.ToEditorsAndOverrides(ctx)
	glice.Notef("\nEditors and overrides extracted")
}

func SavingOverridesFile(ctx context.Context, of *glice.OverridesFile) {
	glice.Notef("\nSaving %s", of.Filepath)
	err := of.Save()
	if err != nil {
		Failf(glice.ExitCannotSaveFile,
			"Failed to create file %s: %s",
			of.Filepath,
			err.Error())
	}
	glice.Notef("\nOverrides files saved.")
}

func LoadingProfileFile(ctx context.Context) *glice.ProjectFile {
	glice.Notef("\nLoading `%s`", glice.ProjectFilename)
	pf, err := glice.LoadProjectFile(glice.GetOptions().SourceDir)
	if err != nil {
		Failf(glice.ExitFileDoesNotExist,
			"Cannot run %s; %s",
			glice.CallerName(),
			err.Error())
	}
	Notef("\nLoaded `%s`", pf.Filepath)
	return pf
}

func ScanningDependencies(ctx context.Context) (deps glice.Dependencies) {
	var err error

	glice.Notef("\nScanning dependencies...")
	deps, err = glice.ScanDependencies(ctx, glice.GetOptions())
	if err != nil {
		Failf(glice.ExitCannotScanDependencies,
			"Failed while scanning dependencies: %s",
			err.Error())
	}
	return deps
}

func SavingProjectFile(ctx context.Context, pf *glice.ProjectFile) {
	glice.Notef("\nSaving %s", glice.ProjectFilename)
	err := pf.Save()
	if err != nil {
		Failf(glice.ExitCannotSaveFile,
			"Failed to save %s: %s",
			pf.Filepath,
			err.Error())
	}
	glice.Notef("\nProject file saved")
}

func CreatingProjectFile(ctx context.Context) (pf *glice.ProjectFile) {
	glice.Notef("\nCreating %s", glice.ProjectFilename)
	pf = glice.NewProjectFile(glice.GetOptions().SourceDir)
	if pf.Exists() {
		Failf(glice.ExitFileExistsCannotOverwrite,
			"Cannot overwrite existing file %s.\nRename or delete file then re-run 'glice init'.",
			pf.Filepath)
	}
	pf.Editors = glice.Editors{
		{Name: "Name goes here", Email: "email-alias@singlestore.com"},
	}
	pf.Overrides = glice.Overrides{
		{
			Import:       "https://github.com/example.com/sample",
			LicenseID:    "MIT",
			VerifiedBy:   "*email-alias",
			LastVerified: glice.Timestamp()[:10],
			Notes:        "This is a sample override added by 'glice init' command",
		},
	}
	glice.Notef("\nFile %s created", pf.Filepath)
	return pf
}

func CreatingOverridesFile(ctx context.Context, onExists int) *glice.OverridesFile {
	glice.Notef("\nCreating %s", glice.OverridesFilename)

	of := glice.NewOverridesFile(glice.GetOptions().SourceDir)
	if !of.Exists() {
		goto end
	}
	switch onExists {
	case glice.WarnLevel:
		Warnf(
			"\nCannot overwrite existing file %s.\nRename or delete file then re-run 'glice overrides'.",
			of.Filepath)
	default:
		Failf(glice.ExitFileExistsCannotOverwrite,
			"\nCannot overwrite existing file %s.\nRename or delete file then re-run 'glice overrides'.",
			of.Filepath)
	}
end:
	return of
}

func AuditingProjectDependencies(ctx context.Context, deps glice.Dependencies) (pf *glice.ProjectFile) {
	pf = LoadingProfileFile(ctx)
	glice.Notef("\nAuditing dependencies...")
	pf.Changes, pf.Disalloweds = pf.AuditDependencies(deps)
	glice.Notef("\nAudit complete.")
	return pf
}

func HasDisalloweds(ctx context.Context, pf *glice.ProjectFile) (has bool) {
	if len(pf.Disalloweds) == 0 {
		glice.Notef("\n")
		glice.Notef("\nOnly allowed licenses detected")
		glice.Notef("\nAudit completed successfully")
		goto end
	}
	has = true
	glice.Errorf("\n")
	glice.Errorf("\nDisallowed licenses detected:")
	glice.Errorf("\n")
	pf.Disalloweds.LogPrint()
	glice.Errorf("\n")
	glice.Errorf("\nAudit FAILED!")
	glice.Errorf("\n\n")
end:
	return has
}

func HandleChanges(ctx context.Context, pf *glice.ProjectFile) {
	changes := pf.Changes
	if !changes.HasChanges() {
		glice.Notef("\n")
		glice.Notef("\nNo chances detected")
		glice.Notef("\n\n")
	} else {
		glice.Notef("\n")
		changes.Print()
	}
}

func ShouldGenerateOverrides(cmd *cobra.Command) bool {
	return cmd.Name() == "overrides" || glice.Flag(cmd, "overrides") == "true"
}

func GeneratingOverrides(ctx context.Context, cmd *cobra.Command, pf *glice.ProjectFile, onExists int) {
	if ShouldGenerateOverrides(cmd) {
		glice.Notef("\n")
		of := CreatingOverridesFile(ctx, onExists)
		ExtractingEditorsAndOOverrides(ctx, pf, of)
		SavingOverridesFile(ctx, of)
		glice.Notef("\n\n")
	}
}

func Warnf(format string, args ...interface{}) {
	glice.Warnf(format, args...)

}
func Notef(format string, args ...interface{}) {
	glice.Notef(format, args...)
}
func Failf(level int, format string, args ...interface{}) {
	glice.Failf(level, format, args...)
}
