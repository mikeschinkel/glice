package glice

import (
	"log"
	"os"
)

const LogFilename = "glice.log"

var logFilepath = SourceDir(LogFilename)

func LogFilepath() string {
	return logFilepath
}
func SetLogFilepath(fp string) {
	logFilepath = fp
}

func LogPrintFunc(show func()) {
	var out = log.Writer()
	opt := GetOptions()
	if !opt.LogVerbosely && opt.LogOuput {
		// If not Logging Verbosely we would be outputting
		// the same information twice
		show()
	}
	log.SetOutput(os.Stderr)
	show()
	log.SetOutput(out)
}
