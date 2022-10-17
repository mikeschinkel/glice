package glice

import "github.com/fatih/color"

type licenseFormat struct {
	name  string
	color color.Attribute
}

var licenseColor = map[string]licenseFormat{
	"other":      {name: "Other", color: color.FgBlue},
	"mit":        {name: "MIT", color: color.FgGreen},
	"lgpl-3.0":   {name: "LGPL-3.0", color: color.FgCyan},
	"mpl-2.0":    {name: "MPL-2.0", color: color.FgHiBlue},
	"agpl-3.0":   {name: "AGPL-3.0", color: color.FgHiCyan},
	"unlicense":  {name: "Unlicense", color: color.FgHiRed},
	"apache-2.0": {name: "Apache-2.0", color: color.FgHiGreen},
	"gpl-3.0":    {name: "GPL-3.0", color: color.FgHiMagenta},
}
