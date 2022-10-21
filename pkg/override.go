package glice

type Overrides []*Override
type OverrideMap map[string]*Override
type Override struct {
	Import       string   `yaml:"dependency"`
	LicenseID    string   `yaml:"license,omitempty"`
	LicenseIDs   []string `yaml:"licenses,omitempty"`
	VerifiedBy   string   `yaml:"verifier"`
	LastVerified string   `yaml:"verified"`
	Notes        string   `yaml:"notes,omitempty"`
}

func NewOverride(dep *Dependency, ed *Editor) *Override {
	return &Override{
		Import:       dep.Import,
		LicenseID:    dep.LicenseID,
		LicenseIDs:   []string{},
		VerifiedBy:   ed.Alias(),
		LastVerified: Timestamp()[:10],
		Notes: "Your dependency-specific notes go here,\n" +
			"e.g. links where you read the license you verified, etc.",
	}
}

// ToMap creates a map of Overrides indexed by Dependency Import
func (ors Overrides) ToMap() OverrideMap {
	newOrs := make(OverrideMap, len(ors))
	for _, or := range ors {
		newOrs[or.Import] = or
	}
	return newOrs
}
