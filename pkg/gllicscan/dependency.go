package gllicscan

import "sort"

type DependencyMap map[string]*Dependency
type Dependencies []*Dependency
type Dependency struct {
	Name           string   `json:"name"`
	Version        string   `json:"version"`
	PackageManager string   `json:"package_manager"`
	Path           string   `json:"path"`
	LicenseIDs     []string `json:"licenses"`
}

func NewDependency(dep *Dependency) *Dependency {
	dep.PackageManager = "go"
	if dep.LicenseIDs == nil {
		dep.LicenseIDs = make([]string, 0)
	}
	return dep
}

func (dep Dependency) Amend(adapter GitLabDependencyAdapter) {
	licIDs := UnionString(dep.LicenseIDs, adapter.GetLicenseIDs())
	sort.Strings(licIDs)
	dep.LicenseIDs = licIDs
}

func GetGitLabDependencyMap(adapters []GitLabDependencyAdapter) GitLabDependencyAdapterMap {
	newMap := make(GitLabDependencyAdapterMap, len(adapters))
	for _, adapter := range adapters {
		newMap[adapter.GetName()] = adapter
	}
	return newMap
}

// Amend reconciles the GitLab Report Dependencies with
// passed-in Dependencies from another source, i.e. go.mod scanning
func (deps Dependencies) Amend(adapters []GitLabDependencyAdapter) Dependencies {
	adaptersMap := GetGitLabDependencyMap(adapters)
	depsMap := deps.ToMap()
	for name, adapter := range adaptersMap {
		dep, ok := depsMap[name]
		if ok {
			dep.Amend(adapter)
			continue
		}
		// GitLab scanning missed a dependency Glice found
		// oo add it into the GitLab Dependency Report.
		//goland:noinspection GoAssignmentToReceiver
		deps = append(deps, NewDependency(&Dependency{
			Name:       name,
			Version:    adapter.GetVersion(),
			Path:       adapter.GetPath(),
			LicenseIDs: adapter.GetLicenseIDs(),
		}))
	}
	return deps
}

// ToMap creates a map indexed by Dependency Import of Dependencies
func (deps Dependencies) ToMap() DependencyMap {
	depMap := make(DependencyMap, len(deps))
	for _, dep := range deps {
		depMap[dep.Name] = dep
	}
	return depMap
}
