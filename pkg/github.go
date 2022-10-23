package glice

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strings"
)

var _ RepositoryAdapter = (*GitHubRepoClient)(nil)

type GitHubRepoClient struct {
	*github.Client
	hostClient     *HostClient
	repoInfoGetter RepoInfoGetter
}

var getter RepositoryAdapter

func GetGitHubRepositoryAdapter(ctx context.Context, r *Repository) (_ RepositoryAdapter, err error) {
	var gc *GitHubRepoClient
	var hc *HostClient
	if getter != nil {
		getter.(*GitHubRepoClient).repoInfoGetter = r
		goto end
	}
	gc = &GitHubRepoClient{}
	hc = NewHostClient()
	gc.hostClient = hc
	gc.repoInfoGetter = r
	hc.RepositoryAdapter = gc
	err = gc.Initialize(ctx)
	getter = gc
end:
	return getter, err
}

func HasGitHubAPIKey() bool {
	return GitHubAPIKey() != ""
}

func GitHubAPIKey() string {
	return strings.TrimSpace(os.Getenv("GITHUB_API_KEY"))
}

func (c *GitHubRepoClient) Initialize(ctx context.Context) (err error) {
	var _c *http.Client
	var isCredentialed bool

	if !HasGitHubAPIKey() {
		Warnf("\nThe environment variable GITHUB_API_KEY has not been set, or is empty; license lookups may fail.\n")
	} else {
		apiKey := GitHubAPIKey()
		ts := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: apiKey,
		})
		_c = oauth2.NewClient(ctx, ts)
		isCredentialed = true
	}
	c.Client = github.NewClient(_c)
	c.hostClient.CanLogIn = isCredentialed
	return err
}

func (c *GitHubRepoClient) SetHostClient(hc *HostClient) {
	c.hostClient = hc
}

func (c *GitHubRepoClient) GetName() string {
	return StripURLScheme(c.repoInfoGetter.GetRepoURL())
}

func (c *GitHubRepoClient) GetRepoURL() string {
	c.checkRepoInfoGetter()
	return c.repoInfoGetter.GetRepoURL()
}

func (c *GitHubRepoClient) GetOrgName() string {
	c.checkRepoInfoGetter()
	return c.repoInfoGetter.GetOrgName()
}

func (c *GitHubRepoClient) GetRepoName() string {
	c.checkRepoInfoGetter()
	return c.repoInfoGetter.GetRepoName()
}

func (c *GitHubRepoClient) checkRepoInfoGetter() {
	if c.repoInfoGetter == nil {
		Failf(ExitRepoInfoGetterIsNil,
			"Must set Repository.repoInfoGetter before calling GitHubRepoClient.%s()",
			CallerName())
	}
}

func (c *GitHubRepoClient) UpVoteRepository(ctx context.Context) (err error) {
	var r *github.Response

	if !c.hostClient.CanLogIn {
		// Have not logged in yet so can't increase star count
		err = fmt.Errorf("unable to star repository '%s'; %w",
			c.GetRepoName(),
			ErrCannotLogin)
		goto end
	}

	// Increment star count for the repository
	r, err = c.Client.Activity.Star(ctx, c.GetOrgName(), c.GetRepoName())
	if err != nil {
		err = fmt.Errorf("unable to increment star count for repository '%s': %s",
			c.GetRepoName(),
			r.Status)
	}
end:
	return err
}

func (c *GitHubRepoClient) GetRepositoryLicense(ctx context.Context, options *Options) (lic *RepositoryLicense, err error) {
	var rl *github.RepositoryLicense

	for {
		rl, _, err = c.Repositories.License(ctx, c.GetOrgName(), c.GetRepoName())
		if err == nil {
			// Request succeeded
			lic = NewRepositoryLicense(LicenseArgs{
				ID:  rl.License.GetSPDXID(),
				URL: rl.GetDownloadURL(),
			})
			if !options.CaptureLicense {
				// CLI switch requested we ignore capturing license content
				break
			}
			lic.Text = rl.GetContent()
			break
		}
		_err, ok := err.(*github.ErrorResponse)
		if !ok {
			// Hmm. Some other kind of error was returned
			// Pass it along in case its helpful
			break
		}
		if _err.Response == nil {
			err = fmt.Errorf("response missing unexpectedly; %w", _err)
			break
		}
		switch _err.Response.StatusCode {
		case http.StatusUnauthorized:
			// Bad credentials?
			err = fmt.Errorf("unauthorized (is your GITHUB_API_KEY correct?); %w", _err)

		case http.StatusNotFound:
			// Anything other than a Not Found or Unauthorized
			err = fmt.Errorf("unexpected error; %w", _err)

		default:
			err = nil
		}
		goto end
	}
end:
	return lic, err
}
