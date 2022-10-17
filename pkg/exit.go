package glice

import (
	"log"
	"os"
)

const (
	exitCannotGetWorkingDir       = 1
	exitCannotGetCacheDir         = 2
	exitCannotCreateCacheDir      = 3
	exitCannotParseDependencies   = 4
	exitCannotInitializeYAMLFile  = 5
	exitCannotStatFile            = 6
	exitYAMLExistsCannotOverwrite = 7
)

func LogAndExit(status int, msg string, args ...interface{}) {
	log.Printf(msg, args...)
	log.SetOutput(os.Stderr)
	log.Printf(msg, args...)
	os.Exit(status)
}
