package glice

import (
	"context"
	"fmt"
)

// RepositoryGetter provides an interface for supporting a repository domain,
// e.g. github.org, bitbucket.org, etc.
type RepositoryGetter interface {
	NameGetter
	RepositoryLicenseGetter
	RepositoryUpVoter
	HostClientSetter
	RepoInfoGetter
	Initializer
}

// RepositoryGetterFunc provides the type to support GetRepositoryGetterFunc's ability
// to return a function value for use by GetRepositoryGetter.
type RepositoryGetterFunc func(context.Context, *Repository) (RepositoryGetter, error)

// GetRepositoryGetterFunc returns a RepositoryGetterFunc given the repository's domain
// TODO Convert this into an opt-in where the domain-specific code opt-in via calling a
//      register function in this file from its an init() function in its own code.
// This function is also used to determine which domains are recognized as providing
// repository services without needed introspection of a ?go-get=1 query on HTTP GET
func GetRepositoryGetterFunc(r *Repository) (f RepositoryGetterFunc) {
	switch r.GetHost() {
	case "github.com":
		f = GetGitHubRepositoryGetter
	}
	return f
}

// GetRepositoryGetter returns an object that implements RepositoryGetter for a given domain
func GetRepositoryGetter(ctx context.Context, r *Repository) (rg RepositoryGetter, err error) {
	rg, err = GetRepositoryGetterFunc(r)(ctx, r)
	if err != nil {
		err = fmt.Errorf("unable to get respository getter for %s; %w", r.GetHost(), err)
		goto end
	}
	if rg == nil {
		Failf(exitHostNotYetSupported,
			"Repositories hosted on %s are not yet supported. Support can be added in ./repository_getter.go.",
			r.GetHost())
	}
end:
	return rg, err
}
