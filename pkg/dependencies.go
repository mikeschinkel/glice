package glice

import (
	"context"
	"fmt"
	"github.com/ribice/glice/v3/pkg/gllicscan"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Dependencies []*Dependency
type DependencyMap map[string]*Dependency

func ScanDependencies(ctx context.Context, options *Options) (ds Dependencies, err error) {
	var repos Repositories
	var deps Dependencies

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

// ToEditorsAndOverrides returns a slice of *Dependency and a slice of unique *Editor
//goland:noinspection GoUnusedParameter
func (deps Dependencies) ToEditorsAndOverrides(ctx context.Context) (editors Editors, overrides Overrides) {
	overrides = make(Overrides, len(deps))
	edMap := make(EditorMap, 0)
	for index, dep := range deps {
		ua, err := GetUserAdapter(dep)
		if err != nil {
			Warnf("Unable to add dependency '%s'; %w",
				dep.Import,
				err)
			continue
		}
		ed := NewEditor(ua)
		overrides[index] = NewOverride(dep, ed)

		id := ed.GetID()
		if _, ok := edMap[id]; !ok {
			edMap[id] = ed
		}
	}
	editors = edMap.ToEditors()
	return editors, overrides
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

// GetLicenseIDs returns a sorted string slice containing the
// list of unique License IDs found in the Dependencies.
// TODO This should be after applying overrides
func (deps Dependencies) GetLicenseIDs() []string {
	licMap := make(map[string]struct{}, 0)
	for _, dep := range deps {
		if _, ok := licMap[dep.LicenseID]; ok {
			continue
		}
		licMap[dep.LicenseID] = struct{}{}
	}
	licIDs := make([]string, 0, len(licMap))
	for licID := range licMap {
		licIDs = append(licIDs, licID)
	}
	sort.Strings(licIDs)
	return licIDs
}

// ToMap creates a map indexed by Dependency Import of Dependencies
func (deps Dependencies) ToMap() DependencyMap {
	depsMap := make(DependencyMap, len(deps))
	for _, dep := range deps {
		depsMap[dep.Import] = dep
	}
	return depsMap
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
		dp = GetSourceDir(dir)
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

var reportHeaderRow = []string{"Dependency", "Repository", "License", "Added"}

func (Dependencies) GetReportHeader() []string {
	return reportHeaderRow
}

// GetOverridesLicenseIDs returns a map of slices of LicenseIDs with map key
// being the Dependency Import and the LicenseIDs being the merger of the
// Dependencies License as reported by the VCS hot (e.g GitHub) and the
// manually entered Licence IDs for the override.
func (deps Dependencies) GetOverridesLicenseIDs(overrides Overrides) (licIDs OverridesLicenseIDMap) {
	orMap := overrides.ToMap()
	depMap := deps.ToMap()
	licIDs = make(OverridesLicenseIDMap, len(overrides))
	for imp, or := range orMap {
		dep, ok := depMap[imp]
		if !ok {
			Warnf("Overridden dependency '%s' not found during scan", imp)
			continue
		}
		var depLicIDs []string
		if dep.LicenseID == NoAssertion {
			depLicIDs = or.LicenseIDs
		} else {
			depLicIDs = gllicscan.UnionStringSlices([]string{dep.LicenseID}, or.LicenseIDs)
		}
		if len(depLicIDs) == 0 {
			depLicIDs = []string{NoAssertion}
			Warnf("Dependency '%s' has '%s' for a license", imp, NoAssertion)
		}
		licIDs[imp] = depLicIDs
	}
	return licIDs
}
