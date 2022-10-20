package glice

import "sort"

type Changes struct {
	Additions Dependencies
	Deletions Dependencies
}

func NewChanges() *Changes {
	return &Changes{
		Additions: make(Dependencies, 0),
		Deletions: make(Dependencies, 0),
	}
}

// HasChanges returns true if there are either old or new changes
func (c *Changes) HasChanges() bool {
	return len(c.Additions) > 0 || len(c.Deletions) > 0
}

// Print outputs all changes, old and new
func (c *Changes) Print() {
	LogPrintFunc(WarnLevel, func() {
		showChanges(c.Additions, "Additions", "These imports were not found in glice.yaml but were found when scanning:")
		showChanges(c.Deletions, "Deletions", "These imports were not found when scanning but were found in glice.yaml:")
	})
}

func showChanges(list Dependencies, _type, descr string) {
	if len(list) == 0 {
		goto end
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Import < list[j].Import
	})
	Notef("\n%s", _type)
	Notef("\n---------")
	Notef("\n%s", descr)
	Notef("\n")
	for _, dep := range list {
		Notef("\n  - %s", dep.Import)
	}
	Notef("\n\n")
end:
}
