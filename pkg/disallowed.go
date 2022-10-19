package glice

import (
	"sort"
	"strconv"
)

type Disalloweds []*Disallowed
type Disallowed struct {
	Import    string
	LicenseID string
}

func NewDisallowed(dep *Dependency) *Disallowed {
	return &Disallowed{
		Import:    dep.Import,
		LicenseID: dep.LicenseID,
	}
}

// HasDisalloweds returns true if Disalloweds has one or more rejections
func (ds Disalloweds) HasDisalloweds() bool {
	return len(ds) > 0
}

// ImportWidth returns the length of the longest Import
func (ds Disalloweds) ImportWidth() (width int) {
	for _, d := range ds {
		n := len(d.Import)
		if n <= width {
			continue
		}
		width = n
	}
	return width
}

// LogPrint outputs all rejections in list individually
func (ds Disalloweds) LogPrint() {
	level := ErrorLevel
	LogPrintFunc(level, func() {
		width := strconv.Itoa(ds.ImportWidth() + 2)
		format := "\n%s: - %-" + width + "s %s"
		sort.Slice(ds, func(i, j int) bool {
			return ds[i].Import < ds[j].Import
		})
		for _, d := range ds {
			LogPrintf(level, format, levels[level], d.Import+":", d.LicenseID)
		}
	})
}
