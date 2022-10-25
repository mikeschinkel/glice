package glice

type RepositoryLicense struct {
	ID   string
	Text string
	URL  string
}

type LicenseArgs struct {
	ID  string
	URL string
}

func NewRepositoryLicense(args LicenseArgs) *RepositoryLicense {
	return &RepositoryLicense{
		ID:  args.ID,
		URL: args.URL,
	}
}

func (l *RepositoryLicense) GetSPDXID() string {
	return l.ID
}
func (l *RepositoryLicense) GetID() string {
	return l.ID
}
func (l *RepositoryLicense) GetText() string {
	return l.Text
}
