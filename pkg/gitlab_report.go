package glice

import (
	"github.com/ribice/glice/v3/pkg/gllicscan"
)

// The purpose of GitLabReport is to provide a bridge between glice.Dependencies
// and gllicscan.Report; i.e. to allow one to be converted to another but to
// only be decoupled via the gllicscan.GitLabReportAdapter interface.

var _ gllicscan.GitLabDependencyAdapter = (*GitLabDependency)(nil)

var _ gllicscan.GitLabReportAdapter = (*GitLabReport)(nil)

type GitLabReport struct {
	Dependencies Dependencies
	LicenseIDMap OverridesLicenseIDMap
	adapters     []gllicscan.GitLabDependencyAdapter
}

func NewGitLabReport(deps Dependencies, overrides Overrides) *GitLabReport {
	return &GitLabReport{
		Dependencies: deps,
		LicenseIDMap: deps.GetOverridesLicenseIDs(overrides),
	}
}

func (gr *GitLabReport) GetLicensesIDs() []string {
	return gr.Dependencies.GetLicenseIDs()
}

func (gr *GitLabReport) GetDependencyAdapters() []gllicscan.GitLabDependencyAdapter {
	if gr.adapters != nil {
		goto end
	}
	gr.adapters = make([]gllicscan.GitLabDependencyAdapter, len(gr.Dependencies))
	for i, dep := range gr.Dependencies {
		gd := NewGitLabDependency(dep)
		gr.adapters[i] = gd
		licIDs, ok := gr.LicenseIDMap[gd.Import]
		if !ok {
			continue
		}
		gd.SetLicenseIDs(licIDs)
	}
end:
	return gr.adapters
}
