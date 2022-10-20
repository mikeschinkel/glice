package glice

type Overrides []*Override
type Override struct {
	DependencyImport string   `yaml:"dependency"`
	LicenseID        string   `yaml:"license,omitempty"`
	LicenseIDs       []string `yaml:"licenses,omitempty"`
	VerifiedBy       string   `yaml:"verifier"`
	LastVerified     string   `yaml:"verified"`
	Notes            string   `yaml:"notes,omitempty"`
}

func NewOverride(dep *Dependency, ed *Editor) *Override {
	return &Override{
		DependencyImport: dep.Import,
		LicenseID:        dep.LicenseID,
		LicenseIDs:       []string{},
		VerifiedBy:       ed.String(),
		LastVerified:     Timestamp()[:10],
		Notes:            "Your specific notes go here\ne.g. links where you verified a license, etc.",
	}
}
