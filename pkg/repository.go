package glice

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	//TagName    = `meta`
	htmlRegion = `head`
	htmlTag    = `meta`
	sourceName = "go-source"
	importName = "go-import"
)

type Repositories = []*Repository

var _ RepoInfoGetter = (*Repository)(nil)

// Repository holds information about the repository
type Repository struct {
	Import   string
	Imports  []string
	url      *url.URL
	license  *RepositoryLicense
	parsed   bool
	resolved bool
}

// ScanRepositories scans the directory providing in options and returns
// a slice of *Repository objects
func ScanRepositories(ctx context.Context, options *Options) ([]*Repository, error) {
	var repos Repositories
	var el ErrorList

	modules, err := ParseModFile(
		options.SourceDir,
		options.IncludeIndirect,
	)
	if err != nil {
		err = fmt.Errorf("unable to list repositories for '%s'; %w",
			options.SourceDir,
			err)
		goto end
	}

	repos = make(Repositories, len(modules))
	el = make(ErrorList, 0)
	for i, module := range modules {
		r := &Repository{
			Import: module,
		}
		Infof("Resolving repository for '%s'", r.Import)
		err = r.ResolveRepository(ctx, options)
		if err != nil {
			el = append(el, err)
			continue
		}
		Infof("\tas '%s'", r.GetURL())
		repos[i] = r
	}
	if len(el) != 0 {
		err = el
	}
end:
	return repos, err
}

// ResolveRepository accepts a Go import string and 1. parses it as a URL
// then either 2. classifies it at a known repo domain, or 3. uses an
// HTTPS request to resolve it from HTML <meta> tags in the response.
// See https://pkg.go.dev/cmd/go#hdr-Remote_import_paths
func (r *Repository) ResolveRepository(ctx context.Context, options *Options) (err error) {
	// Capture the original import as it will get changed if we need
	// recursion, but we want the original one on final return.
	imp := r.Import
	// Add an HTTPS scheme and use *url.ParseModFile() to validate to be
	// a valid URL, If valid, assign returned *url.URL to r.url.
	err = r.ParseImport(imp)
	if err != nil {
		err = fmt.Errorf("invalid repository import '%s'; %w", r.Import, err)
		goto end
	}

	if r.RecognizeKnownRepoDomain() {
		// It is a known repository domain (github.com, etc.) so no need
		// to resolve an HTTP request to determine the repo URL.
		goto end
	}

	// HTTPS request the import as a URL then parse the HTML <meta> tags to determine
	// its repo URL as per https://pkg.go.dev/cmd/go#hdr-Remote_import_paths
	// See paragraph starting with "For code hosted on other servers"
	err = r.ResolveImport(ctx, options)
	if err != nil {
		err = fmt.Errorf("unable to resolve repository import '%s'; %w", r.Import, err)
	}
end:
	// Restore the value of the original import
	r.Import = imp
	return err
}

// SetImport sets a new Import, but keeps track of prior Import values
func (r *Repository) SetImport(imp string) {
	r.Import = imp
	if r.Imports == nil {
		r.Imports = []string{imp}
	}
	if imp != r.Imports[len(r.Imports)-1] {
		r.Imports = append(r.Imports, imp)
	}
}

func (r *Repository) GetOrgName() string {
	return ExtractFieldWithDelimiter(r.GetPath(), 1, '/')
}

func (r *Repository) GetRepoName() string {
	return ExtractFieldWithDelimiter(r.GetPath(), 2, '/')
}

func (r *Repository) GetRepoURL() string {
	r.checkURL()
	return r.url.String()
}

func (r *Repository) GetHost() string {
	r.checkURL()
	return r.url.Host
}

func (r *Repository) GetPath() string {
	r.checkURL()
	return strings.TrimLeft(r.url.Path, "/")
}

func (r *Repository) GetURL() string {
	r.checkURL()
	return r.url.String()
}

func (r *Repository) checkURL() {
	if r.url == nil {
		panic(fmt.Sprintf("Cannot call Repository.%s() before calling Repository.ResolveRepository()",
			CallerName()))
	}
}

func (r *Repository) GetLicenseID() string {
	if r.license == nil {
		return "License inaccessible"
	}
	return r.license.ID
}

func (r *Repository) GetLicenseURL() string {
	if r.license == nil {
		return "http://inaccessible"
	}
	return r.license.URL
}

// ParseImport parses an import string which is typically a URL w/o the scheme So this func
// adds a scheme and the parses the URL using *url.Parse() from the "net/url" import.
func (r *Repository) ParseImport(imp string) (err error) {
	r.SetImport(imp)
	locator := fmt.Sprintf("https://%s", r.Import)
	err = r.ParseURL(locator)
	if err != nil {
		err = fmt.Errorf("invalid repository import: '%s'; %w", r.Import, err)
	}
	r.parsed = true
	return err
}

// ParseURL parses a URL in string format into a *url.URL and assigns
// to Repository.url if successful, or returns an error if not
func (r *Repository) ParseURL(u string) (err error) {
	r.url, err = url.Parse(u)
	if err != nil {
		err = fmt.Errorf("invalid repository URL: '%s'; %w", u, err)
		r.url = &url.URL{}
	}
	return err
}

var importCache = map[string]*url.URL{}

// ResolveImport indirect repos as described here:
// https://golang.org/cmd/go/#hdr-Remote_import_paths
func (r *Repository) ResolveImport(ctx context.Context, options *Options) (err error) {
	var resp *http.Response
	var meta *Meta
	var repoURL string
	var imp = r.Import

	if repoURL, ok := importCache[imp]; ok {
		r.url = repoURL
		goto end
	}

	resp, err = HTTPGetWithContext(ctx, fmt.Sprintf("%s?go-get=1", r.GetURL()))
	if err != nil {
		err = fmt.Errorf("failed to retrieve %s; %w", r.GetURL(), err)
		goto end
	}

	defer MustClose(resp.Body)

	meta, err = GetMetaFromHTMLReader(resp.Body)
	if err != nil {
		err = fmt.Errorf(`failed to extract <meta> from HTML at %s; %w`,
			r.GetURL(),
			err)
		goto end
	}

	repoURL, err = meta.ResolveGoImport(r.Import)
	if err == nil {
		err = r.ParseURL(repoURL)
		if err != nil {
			err = fmt.Errorf(`failed to parse URL '%s' from <meta> in HTML at %s; %w`,
				repoURL,
				r.GetURL(),
				err)
			goto end
		}
	}
	if err == errRequestPrefixInstead || err == nil {
		r.SetImport(StripURLScheme(repoURL))
		err = r.ResolveRepository(ctx, options)
		if err != nil {
			err = fmt.Errorf("unable to resolve repository '%s'; %w", repoURL, err)
			goto end
		}
	}
	if err != nil {
		err = fmt.Errorf(`failed to parse <meta> from HTML at %s; %w`,
			r.GetURL(),
			err)
		goto end
	}

	importCache[imp] = r.url

end:
	return err
}

// ResolveLicense requests the license for the repository
func (r *Repository) ResolveLicense(ctx context.Context, options *Options) (err error) {
	var ra RepositoryGetter

	if r.license != nil {
		goto end
	}

	ra, err = GetRepositoryGetter(ctx, r)
	if err != nil {
		Failf(exitCannotGetRepositoryGetter,
			"unable to get repository getter for %s",
			r.GetHost())
	}

	r.license, err = ra.GetRepositoryLicense(ctx, options)
	if err != nil {
		err = fmt.Errorf("unable to get license for '%s'; %w",
			ra.GetName(),
			err)
		goto end
	}

	// TODO IMO this should be moved out of getting a license and should
	//      be handled by calling a bespoke method to thank all repos.
	//ra.UpVoteRepository(ctx)

end:
	return err
}

// RecognizeKnownRepoDomain inspects r.url.Host and returns true if the domain is
// recognized explicitly by this program, or false otherwise. If it is recognized
// and there are no <meta> tags for ?go-get=1 then it may also update the Host
// and Path of r.url, as applicable.

func (r *Repository) RecognizeKnownRepoDomain() (recognized bool) {
	return GetRepositoryGetterFunc(r) != nil
}

//func (r *Repository) RecognizeKnownRepoDomain() (recognized bool) {
//	_url := r.url
//	switch _url.Host {
//	case "github.com":
//		recognized = true
//	case "google.golang.com", "google.golang.org":
//		//_url.Host = "github.com"
//		//_url.Path = "googleapis/google-api-go-client"
//		//recognized = true
//		recognized = false
//	case "cloud.google.com", "go.octolab.org", "go.uber.org", "golang.org", "gopkg.in", "gotest.tools", "k8s.io", "sigs.k8s.io":
//		// This case is not actually needed, but it documents explicitly the domains that are have tested with.
//		// These will get resolve by requesting the import as a URL with a `go-get=1` query parameter
//		// and then inspecting <meta name="go-import"> and/or <meta name="go-source">
//		recognized = false
//	default:
//		print("")
//	}
//	return recognized
//}
