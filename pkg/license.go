package glice

type LicenseStatus string

const (
	AllowedStatus     LicenseStatus = "allowed"
	DisallowedStatus  LicenseStatus = "disallowed"
	WhiteListedStatus LicenseStatus = "whitelisted"
	UnspecifiedStatus LicenseStatus = "unspecified"
)

type LicenseIDs []string
type LicenseIDMap map[string]exists

func (ids LicenseIDs) ToMap() LicenseIDMap {
	idMap := make(LicenseIDMap, len(ids))
	for _, id := range ids {
		idMap[id] = exists{}
	}
	return idMap
}

type License struct {
	ID              string        `json:"licenseId"`
	Reference       string        `json:"reference"`
	IsDeprecated    bool          `json:"isDeprecatedLicenseId"`
	DetailsURL      string        `json:"detailsUrl"`
	ReferenceNumber int           `json:"referenceNumber"`
	Name            string        `json:"name"`
	SeeAlso         []string      `json:"seeAlso"`
	IsOSIApproved   bool          `json:"isOsiApproved"`
	IsFSFLibre      bool          `json:"isFsfLibre,omitempty"`
	Status          LicenseStatus `json:"status"`
	LastEdited      string        `json:"lastUpdated"`
	UpdatedBy       string        `json:"updatedBy"`
	Text            string        `json:"-"`
}

func NewLicense(args LicenseArgs) *License {
	return &License{}
}

// GetID returns the license's ID value which is one of the identifiers found at spdx.org/licenses
func (l *License) GetID() string {
	return l.ID
}

// GetSPDXID is an alias of GetID
func (l *License) GetSPDXID() string {
	return l.ID
}

// GetReference returns the value of the object's Reference property
func (l *License) GetReference() string {
	return l.Reference
}

// GetIsDeprecated returns the value of the object's IsDeprecated property
func (l *License) GetIsDeprecated() bool {
	return l.IsDeprecated
}

// GetDetailsURL returns the value of the object's DetailsURL property
func (l *License) GetDetailsURL() string {
	return l.DetailsURL
}

// GetReferenceNumber returns the value of the object's ReferenceNumber property
func (l *License) GetReferenceNumber() int {
	return l.ReferenceNumber
}

// GetName returns the value of the object's Name property
func (l *License) GetName() string {
	return l.Name
}

// GetSeeAlso returns the value of the object's SeeAlso property
func (l *License) GetSeeAlso() []string {
	return l.SeeAlso
}

// GetIsOSIApproved returns the value of the object's IsOSIApproved property
func (l *License) GetIsOSIApproved() bool {
	return l.IsOSIApproved
}

// GetIsFSFLibre returns the value of the object's IsFSFLibre property
func (l *License) GetIsFSFLibre() bool {
	return l.IsFSFLibre
}

// GetText returns the value of the object's Text property
func (l *License) GetText() string {
	return l.Text
}
