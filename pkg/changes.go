package glice

import (
	"log"
)

type Changes struct {
	Old []string
	New []string
}

func NewChanges() *Changes {
	return &Changes{
		Old: make([]string, 0),
		New: make([]string, 0),
	}
}

// HasChanges returns true if there are either old or new changes
func (c *Changes) HasChanges() bool {
	return len(c.Old) > 0 || len(c.New) > 0
}

// Print outputs all changes, old and new
func (c *Changes) Print() {
	LogPrintFunc(func() {
		showChanges(c.Old, "Old", "These imports were not found in glice.yaml but were found when scanning.")
		showChanges(c.New, "New", "These imports were not found when scanning but were found in glice.yaml.")
	})
}

func showChanges(list []string, _type, descr string) {
	log.Printf("\nChanges: %s", _type)
	log.Println("------------")
	log.Println(descr)
	for _, imp := range list {
		log.Printf("  - %s\n", imp)
	}
	log.Println("")
}
