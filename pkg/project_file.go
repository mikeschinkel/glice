package glice

import (
	"fmt"
	"path/filepath"
)

const (
	ProjectFileSchemaVersion = "v1"
	ProjectFilename          = "glice.yaml"
)

var _ FilepathGetter = (*ProjectFile)(nil)

type ProjectFile struct {
	Filepath        string       `yaml:"-"`
	SchemaVersion   string       `yaml:"schema"`
	Editors         Editors      `yaml:"editors"`
	Generated       string       `yaml:"generated"`
	AllowedLicenses LicenseIDs   `yaml:"allowed"`
	allowedMap      LicenseIDMap `yaml:"-"`
	Overrides       Overrides    `yaml:"overrides"`
	Dependencies    Dependencies `yaml:"dependencies"`
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
	if pf.AllowedLicenses == nil {
		pf.AllowedLicenses = make([]string, 0)
	}
	if pf.Overrides == nil {
		pf.Overrides = make(Overrides, 0)
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
	}
	return fg.(*ProjectFile), err
}

func (pf *ProjectFile) Initialize() (err error) {
	return CreateYAMLFile(pf)
}

// removeOverridden accepts a DependencyMap and removes any found to be
// overridden in the `glice.yaml` file, returning the smaller map.
func (pf *ProjectFile) removeOverridden(depMap DependencyMap) DependencyMap {
	// First scan the overrides from the glice.yaml file
	for _, _or := range pf.Overrides {
		// If none found in the deps provided as overridden
		if _, ok := depMap[_or.DependencyImport]; !ok {
			// continue looking
			continue
		}
		// If any deps provides WERE found to be overridden
		// then let's remove them from the list of deps
		// TODO: Address when license changes to unacceptable AFTER it was
		//       overridden as acceptable
		delete(depMap, _or.DependencyImport)
	}
	return depMap
}

// IsLicenseAllowed inspects a single scanned dependency to ensure
// it has a proper license, returning false if not.
func (pf *ProjectFile) IsLicenseAllowed(dep *Dependency) (ok bool) {
	_, ok = pf.allowedMap[dep.LicenseID]
	return ok
}

// AuditDependencies returns any disallowed licenses found in the provided dependencies.
// Also returns changes based on the dependencies there were in the glice.yaml file.
func (pf *ProjectFile) AuditDependencies(deps Dependencies) (changes *Changes, disalloweds Dependencies) {
	var scanDeps = pf.removeOverridden(deps.ToMap())
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
	pf.allowedMap = pf.AllowedLicenses.ToMap()
	disalloweds = make(Dependencies, 0)
	for imp, dep := range scanDeps {
		if !pf.IsLicenseAllowed(dep) {
			disalloweds = append(disalloweds, dep)
		}
		if _, ok := fileDeps[imp]; ok {
			continue
		}
		changes.Additions = append(changes.Additions, dep)
	}
	return changes, disalloweds
}
