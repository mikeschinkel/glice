package glice

import (
	"context"
	"io"
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
	UpVoteRepository(context.Context) error
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

type WriterSetter interface {
	SetWriter(io.Writer)
}
type ReportWriter interface {
	WriteReport() error
}
type FileExtensionGetter interface {
	FileExtension() FileExtension
}
type FilepathSetter interface {
	SetFilepath(string)
}
type FormatGetter interface {
	GetFormat() OutputFormat
}
type DependenciesSetter interface {
	SetDependencies(Dependencies)
}
