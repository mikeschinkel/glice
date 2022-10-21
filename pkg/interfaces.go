package glice

import (
	"context"
)

type Chars interface {
	string | []byte | byte | int32
}

type NameGetter interface {
	GetName() string
}
type RepositoryLicenseGetter interface {
	GetRepositoryLicense(context.Context, *Options) (*RepositoryLicense, error)
}
type RepositoryUpVoter interface {
	UpVoteRepository(context.Context)
}
type HostClientSetter interface {
	SetHostClient(client *HostClient)
}
type RepoInfoGetter interface {
	GetOrgName() string
	GetRepoName() string
	GetRepoURL() string
}
type LicenseInfoGetter interface {
	GetSPDXID() string
	GetText() string
	GetURL() string
}
type Initializer interface {
	Initialize(ctx context.Context) error
}

type FilepathGetter interface {
	GetFilepath() string
}
