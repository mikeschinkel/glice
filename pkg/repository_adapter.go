package glice

import (
	"context"
	"fmt"
)

// RepositoryAdapter provides an interface for supporting a repository domain,
// e.g. github.org, bitbucket.org, etc.
type RepositoryAdapter interface {
	NameGetter
	RepositoryLicenseGetter
	RepositoryUpVoter
	HostClientSetter
	RepoInfoGetter
	Initializer
}

// RepositoryAdapterFunc provides the type to support GetRepositoryAdapterFunc
// to return a function value for use by GetRepositoryAdapter.
type RepositoryAdapterFunc func(context.Context, *Repository) (RepositoryAdapter, error)

// GetRepositoryAdapterFunc returns a RepositoryAdapterFunc given the repository's domain
// TODO Convert this into an opt-in where the domain-specific code opt-in via calling a
//      register function in this file from its an init() function in its own code.
// This function is also used to determine which domains are recognized as providing
// repository services without needed introspection of a ?go-get=1 query on HTTP GET
func GetRepositoryAdapterFunc(r *Repository) (f RepositoryAdapterFunc) {
	switch r.GetHost() {
	case "github.com":
		f = GetGitHubRepositoryAdapter
	}
	return f
}

// GetRepositoryAdapter returns an object that implements RepositoryAdapter for a given domain
func GetRepositoryAdapter(ctx context.Context, r *Repository) (ra RepositoryAdapter, err error) {
	ra, err = GetRepositoryAdapterFunc(r)(ctx, r)
	if err != nil {
		err = fmt.Errorf("unable to get respository getter for %s; %w", r.GetHost(), err)
		goto end
	}
	if ra == nil {
		Failf(ExitHostNotYetSupported,
			"Repositories hosted on %s are not yet supported. Support can be added in ./repository_adapter.go.",
			r.GetHost())
	}
end:
	return ra, err
}
