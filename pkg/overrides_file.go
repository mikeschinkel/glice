package glice

import (
	"path/filepath"
	"strings"
)

const (
	OverridesFilename = "overrides.yaml"

	OverridesNotes = "This note is auto-generated by 'glice audit --overrides'.\n" +
		"\n" +
		"The contents of this YAML file should be manually edited and verified\n" +
		"before being copied into 'glice.yaml', this file deleted, then your\n" +
		"'glice.yaml' should be committed to your project's Git repo.\n" +
		"\n" +
		"Also be sure to change any license from 'NOASSERTION' to the correct\n" +
		"license upon discovering what this dependency's license actually is.\n" +
		"\n" +
		"And for other licenses you choose to override that are not in the list\n" +
		"of allowed licenses se sure to include a note for the dependency about\n" +
		"why you are overriding license compliance, and also if you plan for\n" +
		"the override to just be TEMPORARY, and if so be sure to specify the\n" +
		"criteria to be satisfied before the override can be removed.\n"
)

var _ FilepathGetter = (*OverridesFile)(nil)

type OverridesFile struct {
	Filepath      string    `yaml:"-"`
	SchemaVersion string    `yaml:"schema"`
	Notes         string    `yaml:"notes"`
	Editors       Editors   `yaml:"editors"`
	Overrides     Overrides `yaml:"overrides"`
}

func NewOverridesFile(dir string) *OverridesFile {
	pf := &OverridesFile{
		Filepath: filepath.Join(dir, OverridesFilename),
		Notes:    strings.TrimSpace(OverridesNotes),
	}
	pf.ensureValidProperties()
	return pf
}

func (of *OverridesFile) GetFilepath() string {
	return of.Filepath
}

func (of *OverridesFile) Exists() (exists bool) {
	return FileExists(of.Filepath)
}

func (of *OverridesFile) ensureValidProperties() {
	if of.SchemaVersion == "" {
		of.SchemaVersion = ProjectFileSchemaVersion
	}
	if of.Overrides == nil {
		of.Overrides = make(Overrides, 0)
	}
}

func (of *OverridesFile) Save() error {
	return SaveYAMLFile(of)
}
