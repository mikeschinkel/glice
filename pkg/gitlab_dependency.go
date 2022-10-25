package glice

import "github.com/ribice/glice/v3/pkg/gllicscan"

type GitLabDependencyMap map[string]*GitLabDependency
type GitLabDependencies []*GitLabDependency

// ToMap creates a map indexed by GitLabDependency.Name of GitLabDependencies
func (gds *GitLabDependencies) ToMap() GitLabDependencyMap {
	gdMap := make(GitLabDependencyMap, len(*gds))
	for _, gd := range *gds {
		gdMap[gd.GetName()] = gd
	}
	return gdMap
}

type GitLabDependency struct {
	*Dependency
	Version    string
	Path       string
	LicenseIDs []string
}

var _ gllicscan.GitLabDependencyAdapter = (*GitLabDependency)(nil)

// NewGitLabDependency instantiates a new instance of GitLabDependency
// TODO Need to handle Version and Path
func NewGitLabDependency(dep *Dependency) *GitLabDependency {
	return &GitLabDependency{
		Dependency: dep,
		Version:    "",
		Path:       "",
		LicenseIDs: []string{dep.LicenseID},
	}
}

func (gd *GitLabDependency) GetName() string {
	return gd.Import
}

func (gd *GitLabDependency) GetVersion() string {
	return ""
}

func (gd *GitLabDependency) GetPath() string {
	return ""
}

func (gd *GitLabDependency) GetLicenseIDs() []string {
	return []string{gd.LicenseID}
}

func (gd *GitLabDependency) SetLicenseIDs(licIDs []string) {
	gd.LicenseIDs = licIDs
}
