package glice

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

var _ RepositoryAccessor = (*GitHubRepoClient)(nil)

type GitHubRepoClient struct {
	*github.Client
	hostClient     *HostClient
	repoInfoGetter RepoInfoGetter
}

func GetGitHubRepositoryAccessor(ctx context.Context, r *Repository) (_ RepositoryAccessor, err error) {
	gc := &GitHubRepoClient{}
	hc := NewHostClient()
	gc.hostClient = hc
	gc.repoInfoGetter = r
	hc.RepositoryAccessor = gc
	gc.LogIn(ctx)
	return gc, err
}

func (c *GitHubRepoClient) LogIn(ctx context.Context) {
	var _c *http.Client
	var isCredentialed bool

	apiKey := os.Getenv("GITHUB_API_KEY")
	if apiKey == "" {
		log.Println("The environment variable GITHUB_API_KEY has not been set, or is empty; license lookups may fail.")
	} else {
		ts := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: apiKey,
		})
		_c = oauth2.NewClient(ctx, ts)
		isCredentialed = true
	}
	c.Client = github.NewClient(_c)
	c.hostClient.CanLogIn = isCredentialed
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
		panic(fmt.Sprintf("Cannot call GitHubRepoClient.%s() before setting Repository.repoInfoGetter",
			CallerName()))
	}
}

func (c *GitHubRepoClient) UpVoteRepository(ctx context.Context, options *Options) {
	var r *github.Response
	var err error

	if !c.hostClient.CanLogIn {
		// Have not logged in yet so can't increase star count
		goto end
	}

	// Increment star count for the repository
	r, err = c.Client.Activity.Star(ctx, c.GetOrgName(), c.GetRepoName())
	if err != nil {
		log.Printf("Unable to increment star count for repostory '%s': %s", c.GetName(), r.Status)
	}
end:
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
			if options.NoCaptureLicenseText {
				// CLI switch requested we ignore capturing license content
				break
			}
			lic.Text = rl.GetContent()
			break
		}
		er, ok := err.(*github.ErrorResponse)
		if !ok {
			// Hmm. Some other kind of error was returned
			break
		}
		// Pass it along in case its helpful
		err = er
		if er.Response.StatusCode != http.StatusNotFound {
			// Anything other than a Not Found is an error
			break
		}
	}
	return lic, err
}
