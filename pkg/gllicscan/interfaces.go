package gllicscan

type GitLabDependencyAdapterMap map[string]GitLabDependencyAdapter
type GitLabDependencyAdapter interface {
	GetName() string
	GetVersion() string
	GetPath() string
	GetLicenseIDs() []string
	SetLicenseIDs([]string)
}

type GitLabReportAdapter interface {
	GetLicensesIDs() []string
	GetDependencyAdapters() []GitLabDependencyAdapter
}

type FilepathGetter interface {
	GetFilepath() string
}
