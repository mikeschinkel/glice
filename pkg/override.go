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
		VerifiedBy:       ed.Alias(),
		LastVerified:     Timestamp()[:10],
		Notes: "Your dependency-specific notes go here,\n" +
			"e.g. links where you read the license you verified, etc.",
	}
}
