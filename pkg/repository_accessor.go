package glice

import (
	"context"
	"fmt"
)

// RepositoryAccessor provides an interface for supporting a repository domain,
// e.g. github.org, bitbucket.org, etc.
type RepositoryAccessor interface {
	NameGetter
	RepositoryLicenseGetter
	RepositoryUpVoter
	HostClientSetter
	RepoInfoGetter
}

// RepositoryAccessorFunc provides the type to support GetRepositoryAccessorFunc's ability
// to return a function value for use by GetRepositoryAccessor.
type RepositoryAccessorFunc func(context.Context, *Repository) (RepositoryAccessor, error)

// GetRepositoryAccessorFunc returns a RepositoryAccessorFunc given the repository's domain
// TODO Convert this into an opt-in where the domain-specific code opt-in via calling a
//      register function in this file from its an init() function in its own code.
// This function is also used to determine which domains are recognized as providing
// repository services without needed introspection of a ?go-get=1 query on HTTP GET
func GetRepositoryAccessorFunc(r *Repository) (f RepositoryAccessorFunc) {
	switch r.GetHost() {
	case "github.com":
		f = GetGitHubRepositoryAccessor
	}
	return f
}

// GetRepositoryAccessor returns an object that implements RepositoryAccessor for a given domain
func GetRepositoryAccessor(ctx context.Context, r *Repository) (ra RepositoryAccessor, err error) {
	ra, err = GetRepositoryAccessorFunc(r)(ctx, r)
	if err != nil {
		err = fmt.Errorf("unable to get respository accessor for %s; %w", r.GetHost(), err)
		goto end
	}
	if ra == nil {
		panic(fmt.Sprintf("Repositories hosted on %s are not yet supported. They can be added in ./repository_accessor.go.",
			r.GetHost()))
	}
end:
	return ra, err
}
