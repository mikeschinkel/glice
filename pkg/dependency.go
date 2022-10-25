package glice

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// Dependency holds information about a dependency
type Dependency struct {
	r          *Repository
	Import     string `yaml:"import" json:"name,omitempty"`
	RepoURL    string `yaml:"repo" json:"url,omitempty"`
	Host       string `yaml:"-" json:"host,omitempty"`
	Author     string `yaml:"-" json:"author,omitempty"`
	Project    string `yaml:"-" json:"project,omitempty"`
	LicenseID  string `yaml:"license" json:"license"`
	LicenseURL string `yaml:"legalese" json:"legalese"`
	Added      string `yaml:"added" json:"added"`
}

// Repository returns the associated Repository object
func (dep *Dependency) Repository() *Repository {
	return dep.r
}

func (dep *Dependency) GetLicenseFilepath(dir string) string {
	filename := fmt.Sprintf("%s-%s-license.md",
		dep.Author,
		dep.Project)
	return filepath.Join(dir, filename)
}

type SaveDependencyMsgFunc func(dep *Dependency, filepath string)

// SaveLicense saves a license to a file
func (dep *Dependency) SaveLicense(dir string, msgFunc SaveDependencyMsgFunc) (err error) {
	var dec []byte
	var f *os.File
	var fp string

	text := dep.GetLicenseText()
	if text == "" {
		goto end
	}

	dec, err = base64.StdEncoding.DecodeString(text)
	if err != nil {
		err = fmt.Errorf("unable to decode license text for '%s'; %w", dep.Import, err)
		goto end
	}

	fp = dep.GetLicenseFilepath(dir)
	msgFunc(dep, fp)
	f, err = os.Create(fp)
	if err != nil {
		err = fmt.Errorf("unable to create file %s; %w", fp, err)
		goto end
	}
	defer MustClose(f)

	_, err = f.Write(dec)
	if err != nil {
		err = fmt.Errorf("unable to write to file %s; %w", fp, err)
		goto end
	}

	err = f.Sync()
	if err != nil {
		err = fmt.Errorf("unable to synchronize file %s; %w", fp, err)
		goto end
	}
end:
	return err
}

func (dep *Dependency) GetLicenseText() (text string) {
	if dep.r == nil {
		goto end
	}
	if dep.r.license == nil {
		goto end
	}
	text = dep.r.license.GetText()
end:
	return text
}

func (dep *Dependency) GetReportRow() []string {
	return []string{dep.Import, dep.RepoURL, dep.LicenseID, dep.Added}
}

var reportLicenseCol = 2

func (dep *Dependency) GetColorizedReportRow() []string {
	row := dep.GetReportRow()
	row[reportLicenseCol] = dep.GetColorizedLicenseName()
	return row
}

// GetColorizedLicenseName reGetRepoName()turns a colorized name
func (dep *Dependency) GetColorizedLicenseName() (name string) {
	return color.New(GetLicenseColor(dep.LicenseID)).Sprintf(name)
}
