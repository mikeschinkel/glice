package gllicscan

type LicenseMap map[string]*License
type Licenses []*License
type License struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewLicense(id string) *License {
	return &License{
		ID: id,
	}
}

func (ls Licenses) Amend(licenses []string) Licenses {
	licMap := ls.ToMap()
	for _, lic := range licenses {
		_, ok := licMap[lic]
		if ok {
			continue
		}
		//goland:noinspection GoAssignmentToReceiver
		ls = append(ls, NewLicense(lic))
	}
	return ls
}

// ToMap creates a map indexed by License ID of Licenses
func (ls Licenses) ToMap() LicenseMap {
	licMap := make(LicenseMap, len(ls))
	for _, lic := range ls {
		licMap[lic.ID] = lic
	}
	return licMap
}
