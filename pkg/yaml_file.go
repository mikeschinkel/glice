package glice

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	YAMLFileSchemaVersion = "v1"
	YAMLFilename          = "glice.yaml"
)

type YAMLFile struct {
	Filepath        string       `yaml:"-"`
	SchemaVersion   string       `yaml:"schema"`
	Editors         Editors      `yaml:"editors"`
	Generated       string       `yaml:"generated"`
	AllowedLicenses LicenseIDs   `yaml:"allowed"`
	allowedMap      LicenseIDMap `yaml:"-"`
	Overrides       Overrides    `yaml:"overrides"`
	Dependencies    Dependencies `yaml:"dependencies"`
}

func NewYAMLFile(dir string) *YAMLFile {
	yf := &YAMLFile{
		Filepath: filepath.Join(dir, YAMLFilename),
	}
	yf.ensureValidProperties()
	return yf
}

func (yf *YAMLFile) Exists() (exists bool) {
	_, err := os.Stat(yf.Filepath)
	if errors.Is(err, fs.ErrNotExist) {
		goto end
	}
	if err != nil {
		Failf(exitCannotStatFile,
			"Unable to check existence for %s: %s",
			yf.Filepath,
			err.Error())
	}
	exists = true
end:
	return exists
}

func (yf *YAMLFile) ensureValidProperties() {
	if yf.SchemaVersion == "" {
		yf.SchemaVersion = YAMLFileSchemaVersion
	}
	if yf.AllowedLicenses == nil {
		yf.AllowedLicenses = DefaultAllowedLicenses
	}
	if yf.AllowedLicenses == nil {
		yf.AllowedLicenses = make([]string, 0)
	}
	if yf.Overrides == nil {
		yf.Overrides = make(Overrides, 0)
	}
	if yf.Generated == "" {
		yf.Generated = Timestamp()
	}
}

func YAMLFilepath(dir string) string {
	return filepath.Join(dir, YAMLFilename)
}

func LoadYAMLFile(dir string) (yf *YAMLFile, err error) {
	var b []byte

	fp := YAMLFilepath(dir)
	yf = &YAMLFile{
		Filepath: fp,
	}

	if !yf.Exists() {
		err = fmt.Errorf("unable to find %s; %w",
			fp,
			ErrNoYAMLFile)
		goto end
	}

	b, err = os.ReadFile(fp)
	if err != nil {
		err = fmt.Errorf("unable to read Glice YAML file %s; %w", fp, err)
		goto end
	}
	err = yaml.Unmarshal(b, &yf)
	if err != nil {
		err = fmt.Errorf("unable to unmashal Glice YAML file %s; %w", fp, err)
		goto end
	}
	yf.ensureValidProperties()

end:
	yf.Filepath = fp
	return yf, err
}

func (yf *YAMLFile) Init() (err error) {
	var f *os.File
	var b []byte

	f, err = os.Create(yf.Filepath)
	if err != nil {
		err = fmt.Errorf("unable to open file '%s'; %w", yf.Filepath, err)
		goto end
	}

	b, err = yaml.Marshal(yf)
	if err != nil {
		err = fmt.Errorf("unable to encode to YAML; %w", err)
		goto end
	}

	_, err = f.Write(b)
	if err != nil {
		err = fmt.Errorf("unable to write to '%s'; %w", yf.Filepath, err)
		goto end
	}

	err = f.Close()
	if err != nil {
		err = fmt.Errorf("unable to close '%s'; %w", yf.Filepath, err)
		goto end
	}
end:
	return err
}

// removeOverridden accepts a DependencyMap and removes any found to be
// overridden in the YAML file, returning the smaller map.
func (yf *YAMLFile) removeOverridden(depMap DependencyMap) DependencyMap {
	// First scan the overrides from the glice.yaml file
	for _, _or := range yf.Overrides {
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

// auditDependency inspects a single scanned dependency to ensure
// it has a proper license, returning an error if not.
func (yf *YAMLFile) auditDependency(dep *Dependency) (d *Disallowed) {
	// Check to see if the license is d
	if _, ok := yf.allowedMap[dep.LicenseID]; !ok {
		d = NewDisallowed(dep)
	}
	return d
}

// AuditDependencies returns any disallowed licenses found in the provided dependencies.
// Also returns changes based on the dependencies there were in the glice.yaml file.
func (yf *YAMLFile) AuditDependencies(deps Dependencies) (changes *Changes, ds Disalloweds) {
	var scanDeps = yf.removeOverridden(deps.ToMap())
	var fileDeps = yf.Dependencies.ToMap()

	// Review the file dependencies to see if there are any dependencies not found
	// when scanning the go.mod file but that were previously in glice.yaml.
	changes = NewChanges()
	for _, fd := range yf.Dependencies {
		if _, ok := scanDeps[fd.Import]; ok {
			continue
		}
		changes.Deletions = append(changes.Deletions, fd.Import)
	}

	// Review the if there are any with disallowed licenses and
	// also to see if we found new dependencies when scanning.
	yf.allowedMap = yf.AllowedLicenses.ToMap()
	ds = make(Disalloweds, 0)
	for imp, dep := range scanDeps {
		disallowed := yf.auditDependency(dep)
		if disallowed != nil {
			ds = append(ds, disallowed)
		}
		if _, ok := fileDeps[imp]; ok {
			continue
		}
		changes.Additions = append(changes.Additions, imp)
	}
	return changes, ds
}
