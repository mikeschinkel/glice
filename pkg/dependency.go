package glice

import (
	"context"
	"fmt"
	"github.com/fatih/color"
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
}

func GetDependencyFromRepository(ctx context.Context, r *Repository) *Dependency {
	return &Dependency{
		r:          r,
		Import:     r.Import,
		RepoURL:    r.GetURL(),
		Host:       r.GetHost(),
		Author:     r.GetOrgName(),
		Project:    r.GetRepoName(),
		LicenseID:  r.GetLicenseID(),
		LicenseURL: r.GetLicenseURL(),
	}
}

func (d *Dependency) GetLicenseText() (text string) {
	if d.r == nil {
		goto end
	}
	if d.r.license == nil {
		goto end
	}
	text = d.r.license.GetText()
end:
	return text
}

func (d *Dependency) GetColor() (clr color.Attribute) {
	var lf licenseFormat
	var ok bool
	if lf, ok = licenseColor[d.LicenseID]; !ok {
		clr = color.FgYellow
		goto end
	}
	clr = lf.color
end:
	return clr
}

// GetColorizedLicenseName reGetRepoName()turns a colorized name
func (d *Dependency) GetColorizedLicenseName() (name string) {
	var lf licenseFormat
	var ok bool

	if lf, ok = licenseColor[d.LicenseID]; !ok {
		name = d.LicenseID
	} else {
		name = lf.name
	}
	return color.New(d.GetColor()).Sprintf(name)
}

func ScanDependencies(options *Options) (ds Dependencies, err error) {
	var repos Repositories
	var deps Dependencies

	ctx := context.Background()

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
		deps[i] = GetDependencyFromRepository(ctx, r)
	}
end:
	return deps, err
}
