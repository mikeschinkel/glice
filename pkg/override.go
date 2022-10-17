package glice

type Overrides = []*Override

type Override struct {
	DependencyImport string `yaml:"dependency"`
	LicenseID        string `yaml:"license,omitempty"`
	LicenseIDs       string `yaml:"licenses,omitempty"`
	VerifiedBy       string `yaml:"verifier"`
	LastVerified     string `yaml:"verified"`
	Notes            string `yaml:"notes,omitempty"`
}
