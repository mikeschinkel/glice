package glice

func NewHostClient() *HostClient {
	return &HostClient{}
}

type HostClient struct {
	RepositoryGetter
	CanLogIn bool
}
