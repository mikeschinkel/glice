package glice

type LicenseList struct {
	LicenseListVersion string    `json:"licenseListVersion"`
	Licenses           []License `json:"licenses"`
	ReleaseDate        string    `json:"releaseDate"`
	ETag               string    `json:"etag"`
	SourceURL          string    `json:"sourceUrl"`
}

func NewLicenseList() *LicenseList {
	return &LicenseList{
		SourceURL: LicenseListSourceURL,
	}
}
