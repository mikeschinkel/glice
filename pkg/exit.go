package glice

import (
	"log"
	"os"
)

const (
	exitAuditFoundDisallowedLicenses  = 1
	exitCannotGetWorkingDir           = 2
	exitCannotGetCacheDir             = 3
	exitCannotCreateCacheDir          = 4
	exitCannotParseDependencies       = 5
	exitCannotInitializeYAMLFile      = 6
	exitCannotStatFile                = 7
	exitYAMLFileExistsCannotOverwrite = 8
	exitYAMLFileDoesNotExist          = 9
)

func LogAndExit(status int, msg string, args ...interface{}) {
	log.Printf(msg, args...)
	log.SetOutput(os.Stderr)
	log.Printf("\n"+msg, args...)
	os.Exit(status)
}
