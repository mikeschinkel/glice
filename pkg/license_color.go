package glice

import (
	"github.com/fatih/color"
	"strings"
)

func GetLicenseColor(licenseID string) (fg color.Attribute) {
	switch {
	case strings.HasPrefix(licenseID, "MIT"):
		fg = color.FgHiCyan
	case strings.HasPrefix(licenseID, "MPL-"):
		fg = color.FgHiYellow
	case strings.HasPrefix(licenseID, "BSD-"):
		fg = color.FgHiGreen
	case strings.HasPrefix(licenseID, "Apache-"):
		fg = color.FgHiBlue
	case strings.HasPrefix(licenseID, "CC-"):
		fg = color.FgHiMagenta
	case strings.Contains(licenseID, "GPL"):
		fg = color.FgHiRed
	case licenseID == "ISC":
		fg = color.FgGreen
	case licenseID == "Unlicense":
		fg = color.FgCyan
	default:
		fg = color.FgWhite
	}
	return fg
}
