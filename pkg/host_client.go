package glice

func NewHostClient() *HostClient {
	return &HostClient{}
}

type HostClient struct {
	RepositoryAdapter
	CanLogIn bool
}
