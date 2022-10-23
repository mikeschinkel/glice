package glice

type Overrides []*Override
type OverrideMap map[string]*Override
type Override struct {
	editor       *Editor
	Import       string   `yaml:"import"`
	LicenseIDs   []string `yaml:"licenses"`
	VerifiedBy   string   `yaml:"verifier"`
	LastVerified string   `yaml:"verified"`
	Notes        string   `yaml:"notes,omitempty"`
}

type MarshallableOverride Override

func NewOverride(dep *Dependency, ed *Editor) *Override {
	return &Override{
		Import:       dep.Import,
		LicenseIDs:   []string{dep.LicenseID},
		VerifiedBy:   ed.ID,
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
