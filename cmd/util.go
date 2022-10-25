package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ribice/glice/v3/pkg"
	"github.com/ribice/glice/v3/pkg/gllicscan"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// WARNING: These functions are designed to support commands and are future mutable.
//          They ARE NOT part of the published API and are subject to change.

func Infof(format string, args ...interface{}) {
	glice.Infof(format, args...)
}

func Notef(format string, args ...interface{}) {
	glice.Notef(format, args...)
}
func NoteBegin() {
	Notef("\n")
}
func NoteEnd() {
	Notef("\n\n")
}

func Errorf(format string, args ...interface{}) {
	Errorf(format, args...)
}
func ErrorBegin() {
	Errorf("\n")
}
func ErrorSeparator() {
	Errorf("\n")
}
func ErrorEnd() {
	Errorf("\n\n")
}

func Warnf(format string, args ...interface{}) {
	glice.Warnf(format, args...)
}

func Failf(level int, format string, args ...interface{}) {
	glice.Failf(level, format, args...)
}

func ExtractingEditorsAndOOverrides(ctx context.Context, pf *glice.ProjectFile, of *glice.OverridesFile) {
	Notef("\nExtracting editors and overrides")
	of.Editors, of.Overrides = pf.Disalloweds.ToEditorsAndOverrides(ctx)
	Notef("\nEditors and overrides extracted")
}

//goland:noinspection GoUnusedParameter
func SavingOverridesFile(ctx context.Context, of *glice.OverridesFile) {
	Notef("\nSaving %s", of.Filepath)
	err := of.Save()
	if err != nil {
		Failf(glice.ExitCannotSaveFile,
			"\nFailed to create file %s: %s",
			of.Filepath,
			err.Error())
	}
	Notef("\nOverrides files saved.")
}

//goland:noinspection GoUnusedParameter
func LoadProjectFile(ctx context.Context) *glice.ProjectFile {
	Notef("\nLoading `%s`", glice.ProjectFilename)
	pf, err := glice.LoadProjectFile(glice.GetOptions().SourceDir)
	if err != nil {
		Failf(glice.ExitFileDoesNotExist,
			"\nCannot run %s; %s",
			glice.CallerName(),
			err.Error())
	}
	Notef("\nFile `%s` loaded", pf.Filepath)
	return pf
}

func ScanDependencies(ctx context.Context) (deps glice.Dependencies) {
	var err error

	Notef("\nScanning dependencies...")
	deps, err = glice.ScanDependencies(ctx, glice.GetOptions())
	if err != nil {
		Failf(glice.ExitCannotScanDependencies,
			"\nFailed while scanning dependencies: %s",
			err.Error())
	}
	return deps
}

// SaveProjectFile saves the passed project file to disk
//goland:noinspection GoUnusedParameter
func SaveProjectFile(ctx context.Context, pf *glice.ProjectFile) {
	Notef("\nSaving %s", glice.ProjectFilename)
	backups, err := pf.Save()
	if err != nil {
		Failf(glice.ExitCannotSaveFile,
			"\nFailed to save %s: %s",
			pf.Filepath,
			err.Error())
	}
	if backups != nil {
		for _, bu := range backups[1:] {
			Notef("\nBackup file %s created", bu)
		}
	}
	Notef("\nProject file saved")
}

// CreateProjectFile creates a new object representing the project file
//goland:noinspection GoUnusedParameter
func CreateProjectFile(ctx context.Context) (pf *glice.ProjectFile) {
	Notef("\nCreating %s", glice.ProjectFilename)
	pf = glice.NewProjectFile(glice.GetOptions().SourceDir)
	if pf.Exists() {
		Failf(glice.ExitFileExistsCannotOverwrite,
			"\nCannot overwrite existing file %s.\nRename or delete file then re-run 'glice init'.",
			pf.Filepath)
	}
	pf.Editors = glice.Editors{
		{Name: "Name goes here", Email: "email-alias@singlestore.com"},
	}
	pf.Overrides = glice.Overrides{
		{
			Import:       "https://github.com/example.com/sample",
			LicenseIDs:   []string{"MIT"},
			VerifiedBy:   "email-alias",
			LastVerified: glice.Timestamp()[:10],
			Notes:        "This is a sample override added by 'glice init' command",
		},
	}
	Notef("\nFile %s created", pf.Filepath)
	return pf
}

// CreatingOverridesFile creates a new object representing the overrides file
//goland:noinspection GoUnusedParameter
func CreatingOverridesFile(ctx context.Context, onExists int) *glice.OverridesFile {
	Notef("\nCreating %s", glice.OverridesFilename)

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
	pf = LoadProjectFile(ctx)
	Notef("\nAuditing dependencies...")
	pf.Changes, pf.Disalloweds = pf.AuditDependencies(deps)
	Notef("\nAudit complete.")
	return pf
}

// HandleDisalloweds processes any dependencies disallowed by license
// and returns true if there were any disalloweds.
//goland:noinspection GoUnusedParameter
func HandleDisalloweds(ctx context.Context, pf *glice.ProjectFile) (has bool) {
	if len(pf.Disalloweds) == 0 {
		NoteBegin()
		Notef("\nNo disallowed licenses detected")
		Notef("\nAudit completed successfully")
		NoteEnd()
		goto end
	}
	has = true
	ErrorBegin()
	Errorf("\nDisallowed licenses detected:")
	ErrorSeparator()
	pf.Disalloweds.LogPrint()
	ErrorSeparator()
	Errorf("\nAudit FAILED!")
	ErrorEnd()
end:
	return has
}

// HandleChanges processes any additions and/or deletions when comparing dependencies
// in the project file with the dependences found on disk during scanning.
//goland:noinspection GoUnusedParameter
func HandleChanges(ctx context.Context, pf *glice.ProjectFile) {
	changes := pf.Changes
	if !changes.HasChanges() {
		NoteBegin()
		Notef("\nNo changes detected")
		NoteEnd()
	} else {
		NoteBegin()
		changes.Print()
	}
}

func ShouldGenerateOverrides(cmd *cobra.Command) bool {
	return cmd.Name() == "overrides" || Flag(cmd, "overrides") == "true"
}

func GeneratingOverrides(ctx context.Context, cmd *cobra.Command, pf *glice.ProjectFile, onExists int) {
	if ShouldGenerateOverrides(cmd) {
		NoteBegin()
		of := CreatingOverridesFile(ctx, onExists)
		ExtractingEditorsAndOOverrides(ctx, pf, of)
		SavingOverridesFile(ctx, of)
		NoteEnd()
	}
}

func GetReportWriterAdapter(cmd *cobra.Command) (adapter glice.ReportWriterAdapter) {
	var err error
	var f *os.File

	Notef("\nAcquiring report writer adapter")

	var filename = Flag(cmd, "filename")
	Notef("\nReport filename is '%s'", filename)

	var format = glice.OutputFormat(Flag(cmd, "format"))
	Notef("\nReport format is '%s'", format)

	adapter, err = glice.GetReportWriterAdapter(format)
	if err != nil {
		Failf(glice.ExitCannotGetReportWriterAdapter,
			"\nCannot get report writer adapter for '%s' format; %w",
			glice.ProjectFilename,
			err)
		err = fmt.Errorf("unable to get a report writer adapter for outputting '%s'; %w",
			format,
			err)
	}

	if format != DefaultReportFormat {
		// If user specified a different format, change the extension.
		filename = glice.ReplaceFileExtension(DefaultReportFilename, adapter.FileExtension())
	}

	fp := glice.GetSourceDir(filename)
	if glice.FileExists(fp) {
		Failf(glice.ExitFileExistsCannotOverwrite,
			"\nCannot generate report. The file %s already exists. Rename or delete and then rerun the `%s report save` command.",
			fp,
			glice.CLIName)
	}
	Notef("\nReport to be written to %s", fp)
	adapter.SetFilepath(fp)
	f, err = os.Create(fp)
	if err != nil {
		Failf(glice.ExitCannotCreateFile,
			"\nCannot create file %s for report; %w",
			fp,
			err)
	}
	adapter.SetWriter(f)
	Notef("\nReport writer adapter acquired")
	return adapter
}

func WriteReport(adapter glice.ReportWriterAdapter) {
	Notef("\nWriting report")
	err := adapter.WriteReport()
	if err != nil {
		Failf(glice.ExitCannotWriteReport,
			"\nUnable to get a report writer for outputting dependency report in '%s' format; %w",
			adapter.GetFormat(),
			err)
	}
	Notef("\nReport written")
}

// Flag returns the string value of a cobra.Command pFlag.
func Flag(cmd *cobra.Command, name string) (strVal string) {
	var value pflag.Value
	flag := cmd.Flags().Lookup(name)
	if flag == nil {
		glice.Warnf("Flag '%s' not found for the `%s `%s` command",
			name,
			glice.CLIName,
			cmd.Name())
		goto end
	}
	value = flag.Value
	if value == nil {
		glice.Warnf("The value of flag '%s' for the `%s %s` command is unexpectedly nil",
			name,
			glice.CLIName,
			cmd.Name())
		goto end
	}
	strVal = value.String()
end:
	return strVal
}

func LoadJSONReportFromGitLab(cmd *cobra.Command) (rpt *gllicscan.Report) {
	var err error

	path := Flag(cmd, "path")
	filename := Flag(cmd, "filename")
	fp := filepath.Join(path, filename)
	rpt, err = gllicscan.LoadReport(fp)

	switch {
	case errors.Is(err, gllicscan.ErrFileDoesNotExist):
		Warnf("\nCannot load %s as it does not exist.", fp)
		Warnf("\nCreating anew instead.")
	case errors.Is(err, gllicscan.ErrCannotReadFile):
		Failf(glice.ExitCannotReadFile, "\nCannot load GitLab report; %w", err)
	case errors.Is(err, gllicscan.ErrCannotUnmarshalJSON):
		Failf(glice.ExitCannotUnmarshalJSON, "\nCannot load GitLab report; %w", err)
	case err != nil:
		Failf(glice.ExitUnexpectedError, "\nCannot load GitLab report %s; %w", fp, err)
	}

	return rpt
}

//goland:noinspection GoUnusedParameter
func SaveJSONReportForGitLab(ctx context.Context, rpt *gllicscan.Report) {
	Notef("\nSaving %s", gllicscan.ReportFilename)
	err := rpt.Save()
	if err != nil {
		Failf(glice.ExitCannotSaveFile,
			"\nFailed to save %s: %s",
			rpt.Filepath,
			err.Error())
	}
	Notef("\nReport file saved")
}
