package glice

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var _ FilepathGetter = (*ProjectFile)(nil)

type ProjectFile struct {
	Filepath        string       `yaml:"-"`
	SchemaVersion   string       `yaml:"schema"`
	Editors         Editors      `yaml:"editors"`
	Generated       string       `yaml:"generated"`
	AllowedLicenses LicenseIDs   `yaml:"allowed"`
	Overrides       Overrides    `yaml:"overrides"`
	Dependencies    Dependencies `yaml:"dependencies"`
	allowedMap      LicenseIDMap `yaml:"-"`
	overrideMap     OverrideMap  `yaml:"-"`
	Disalloweds     Dependencies `yaml:"-"`
	Changes         *Changes     `yaml:"-"`
}

func NewProjectFile(dir string) *ProjectFile {
	pf := &ProjectFile{
		Filepath: filepath.Join(dir, ProjectFilename),
	}
	pf.ensureValidProperties()
	return pf
}

func (pf *ProjectFile) GetFilepath() string {
	return pf.Filepath
}

func (pf *ProjectFile) Exists() (exists bool) {
	return FileExists(pf.Filepath)
}

func (pf *ProjectFile) ensureValidProperties() {
	if pf.SchemaVersion == "" {
		pf.SchemaVersion = ProjectFileSchemaVersion
	}
	if pf.AllowedLicenses == nil {
		pf.AllowedLicenses = DefaultAllowedLicenses
	}
	if pf.Overrides == nil {
		pf.Overrides = make(Overrides, 0)
	}
	if pf.Dependencies == nil {
		pf.Dependencies = make(Dependencies, 0)
	}
	if pf.Generated == "" {
		pf.Generated = Timestamp()
	}
}

func GetProjectFilepath(dir string) string {
	return filepath.Join(dir, ProjectFilename)
}

func LoadProjectFile(dir string) (pf *ProjectFile, err error) {
	var fg FilepathGetter

	pf = NewProjectFile(dir)
	fg, err = LoadYAMLFile(pf.Filepath, pf)
	if err != nil {
		err = fmt.Errorf("unable to load %s; %w",
			fg.GetFilepath(), err)
		goto end
	}
	err = pf.ValidateProperties()
end:
	return fg.(*ProjectFile), err
}

func (pf *ProjectFile) ValidateProperties() (err error) {
	var msg []string

	if pf.SchemaVersion == "" {
		pf.SchemaVersion = ProjectFileSchemaVersion
	}
	if pf.AllowedLicenses == nil {
		msg = append(msg, "no allowed licenses are set")
	}
	if pf.Overrides == nil {
		pf.Overrides = make(Overrides, 0)
	}
	if pf.Dependencies == nil {
		msg = append(msg, "no dependencies are set")
	}
	if len(msg) > 0 {
		err = errors.New(strings.Join(msg, ", "))
	}
	return err
}

// Backup creates a backup `glice.yaml` project file with a ".bak"
// or ".<n>.bak" extension while maintaining all prior backups.
func (pf *ProjectFile) Backup() ([]string, error) {
	return BackupFile(pf.GetFilepath(), DeleteOriginalFile)
}

func (pf *ProjectFile) Save() (backups []string, err error) {
	Notef("\nBacking up existing project file(s)")
	backups, err = pf.Backup()
	if err != nil {
		err = fmt.Errorf("unable to backup %s file; %w", ProjectFilename, err)
		goto end
	}
	err = SaveYAMLFile(pf)
	if err != nil {
		err = fmt.Errorf("unable to save %s file; %w", ProjectFilename, err)
		goto end
	}
end:
	return backups, err
}

// IsLicenseAllowed inspects a single scanned dependency to ensure
// it has a proper license, returning false if not.
func (pf *ProjectFile) IsLicenseAllowed(dep *Dependency) (ok bool) {
	if pf.overrideMap == nil {
		pf.overrideMap = pf.Overrides.ToMap()
	}
	if _, ok = pf.overrideMap[dep.Import]; ok {
		goto end
	}
	if pf.allowedMap == nil {
		pf.allowedMap = pf.AllowedLicenses.ToMap()
	}
	_, ok = pf.allowedMap[dep.LicenseID]
end:
	return ok
}

// AuditDependencies returns any disallowed licenses found in the provided dependencies.
// Also returns changes based on the dependencies there were in the glice.yaml file.
func (pf *ProjectFile) AuditDependencies(deps Dependencies) (changes *Changes, disalloweds Dependencies) {
	var scanDeps = deps.ToMap()
	var fileDeps = pf.Dependencies.ToMap()

	// Review the file dependencies to see if there are any dependencies not found
	// when scanning the go.mod file but that were previously in glice.yaml.
	changes = NewChanges()
	for _, fd := range pf.Dependencies {
		if _, ok := scanDeps[fd.Import]; ok {
			continue
		}
		changes.Deletions = append(changes.Deletions, fd)
	}

	// Review the if there are any with disallowed licenses and
	// also to see if we found new dependencies when scanning.
	disalloweds = make(Dependencies, 0)
	for imp, dep := range scanDeps {
		if !pf.IsLicenseAllowed(dep) {
			disalloweds = append(disalloweds, dep)
		}
		if _, ok := fileDeps[imp]; ok {
			continue
		}
		changes.Additions = append(changes.Additions, dep)
		pf.Dependencies = append(pf.Dependencies, dep)
	}
	return changes, disalloweds
}
