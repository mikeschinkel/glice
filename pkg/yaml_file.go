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
	AllowedLicenses []string     `yaml:"allowed"`
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
		LogAndExit(exitCannotStatFile, "Unable to check existence for %s: %s",
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
	yf = &YAMLFile{}

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
