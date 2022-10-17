package glice

func NewHostClient() *HostClient {
	return &HostClient{}
}

type HostClient struct {
	RepositoryAccessor
	CanLogIn bool
}
