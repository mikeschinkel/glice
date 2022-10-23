package glice

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Dependencies []*Dependency
type DependencyMap map[string]*Dependency

// ToMap creates a map indexed by Dependency Import of Dependencies
func (deps Dependencies) ToMap() DependencyMap {
	newDeps := make(DependencyMap, len(deps))
	for _, dep := range deps {
		newDeps[dep.Import] = dep
	}
	return newDeps
}

// ToEditorsAndOverrides returns a slice of *Dependency and a slice of unique *Editor
//goland:noinspection GoUnusedParameter
func (deps Dependencies) ToEditorsAndOverrides(ctx context.Context) (editors Editors, overrides Overrides) {
	overrides = make(Overrides, len(deps))
	edMap := make(EditorMap, 0)
	for index, dep := range deps {
		eg, err := GetEditorGetter(dep)
		if err != nil {
			Warnf("Unable to add dependency '%s'; %w",
				dep.Import,
				err)
			continue
		}
		ed := eg.GetEditor()
		overrides[index] = NewOverride(dep, ed)

		id := ed.GetID()
		if _, ok := edMap[id]; !ok {
			edMap[id] = ed
		}
	}
	editors = edMap.ToEditors()
	return editors, overrides
}

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

func GetDependencyFromRepository(r *Repository) *Dependency {
	return &Dependency{
		r:          r,
		Import:     r.Import,
		RepoURL:    r.GetURL(),
		Host:       r.GetHost(),
		Author:     r.GetOrgName(),
		Project:    r.GetRepoName(),
		LicenseID:  r.GetLicenseID(),
		LicenseURL: r.GetLicenseURL(),
		Added:      Timestamp(),
	}
}

func ScanDependencies(ctx context.Context, options *Options) (ds Dependencies, err error) {
	var repos Repositories
	var deps Dependencies

	//TODO Handle this concern somewhere
	//if thanks && githubAPIKey == "" {
	//	return ErrNoAPIKey
	//}

	repos, err = ScanRepositories(ctx, options)
	if err != nil {
		goto end
	}

	Notef("\nFound %d dependencies", len(repos))
	Notef("\nResolving licenses...")

	deps = make(Dependencies, len(repos))
	for i, r := range repos {
		Infof("\nFetching license for: %s", r.Import)
		err = r.ResolveLicense(ctx, GetOptions())
		if err != nil {
			err = fmt.Errorf("failed to resolve license; %w", err)
			goto end
		}
		deps[i] = GetDependencyFromRepository(r)
	}
end:
	return deps, err
}

// Repository returns the associated Repository object
func (dep *Dependency) Repository() *Repository {
	return dep.r
}

// ImportWidth returns the length of the longest Import
func (deps Dependencies) ImportWidth() (width int) {
	for _, d := range deps {
		n := len(d.Import)
		if n <= width {
			continue
		}
		width = n
	}
	return width
}

// LogPrint outputs all rejections in list individually
func (deps Dependencies) LogPrint() {
	level := ErrorLevel
	LogPrintFunc(level, func() {
		width := strconv.Itoa(deps.ImportWidth() + 2)
		format := "\n%s: - %-" + width + "s %s"
		sort.Slice(deps, func(i, j int) bool {
			return deps[i].Import < deps[j].Import
		})
		for _, d := range deps {
			LogPrintf(level, format, LogLevels[level], d.Import+":", d.LicenseID)
		}
	})
}

// SaveLicenses writes all the dependency licenses each to their own file
func (deps Dependencies) SaveLicenses(dir string, msgFunc SaveDependencyMsgFunc) (err error) {
	var dp string

	if deps == nil {
		goto end
	}

	if len(deps) < 1 {
		goto end
	}

	if filepath.IsAbs(dir) {
		dp = dir
	} else {
		dp = SourceDir(dir)
	}

	err = os.Mkdir(dp, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("unable to create directory %s in which to save licenses; %w",
			dp,
			err)
		goto end
	}

	for _, dep := range deps {
		err = dep.SaveLicense(dp, msgFunc)
		if err != nil {
			Warnf("Unable to save license for %s; %s", dep.Import, err.Error())
		}
	}
end:
	return err
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

var reportHeaderRow = []string{"Dependency", "Repository", "License", "Added"}
var reportLicenseCol = 2

func (Dependencies) GetReportHeader() []string {
	return reportHeaderRow
}

func (dep *Dependency) GetColorizedReportRow() []string {
	row := dep.GetReportRow()
	row[reportLicenseCol] = dep.GetColorizedLicenseName()
	return row
}

// GetColorizedLicenseName reGetRepoName()turns a colorized name
func (dep *Dependency) GetColorizedLicenseName() (name string) {
	return color.New(GetLicenseColor(dep.LicenseID)).Sprintf(name)
}
