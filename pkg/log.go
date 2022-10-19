package glice

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	AllLevel   = 0
	InfoLevel  = 1
	NoteLevel  = 2
	WarnLevel  = 3
	ErrorLevel = 4
	FailLevel  = 5
)

var levels = []string{
	"ALL",
	"INFO",
	"NOTE",
	"WARN",
	"ERROR",
	"FAIL",
}

const LogFilename = "glice.log"

var logFilepath = SourceDir(LogFilename)
var logFile *os.File

func LevelLabel(level int) (ll string) {
	if ll = levels[level]; ll == "" {
		ll = fmt.Sprintf("INVALID_LEVEL[%d]", level)
	}
	return ll
}

func InitializeLogging(fp string) (err error) {
	logFile, err = os.Create(fp)
	if err != nil {
		Warnf("Unable initialize logging; %s", err.Error())
		goto end
	}
	log.SetOutput(logFile)
end:
	return
}

func LogFilepath() string {
	return logFilepath
}

func SetLogFilepath(fp string) {
	logFilepath = fp
}

func Infof(format string, args ...interface{}) {
	LogWithLabelPrintf(InfoLevel, format, args...)
}

func Notef(format string, args ...interface{}) {
	LogWithLabelPrintf(NoteLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	LogWithLabelPrintf(WarnLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	LogWithLabelPrintf(ErrorLevel, format, args...)
}

func Failf(status int, format string, args ...interface{}) {
	LogWithLabelPrintf(FailLevel, format, args...)
	os.Exit(status)
}

func LogWithLabelPrintf(level int, format string, args ...interface{}) {
	if len(format) > 0 {
		if format[0] == '\n' {
			format = fmt.Sprintf("\n%s: %s", LevelLabel(level), format[1:])
		} else {
			format = fmt.Sprintf("%s: %s", LevelLabel(level), format)
		}
	}
	LogPrintf(level, format, args...)
}

func LogPrintf(level int, format string, args ...interface{}) {
	LogPrintFunc(level, func() {
		fmt.Printf(format, args...)
	})
}

func LoggingToFile() bool {
	return logFile != nil
}

func LogPrintFunc(level int, showFunc func()) {
	var out = log.Writer()
	var writer io.Writer

	opt := GetOptions()

	logFunc := func() {
		log.SetOutput(writer)
		showFunc()
	}

	if level >= opt.VerbosityLevel {
		switch {
		case level > NoteLevel:
			writer = os.Stderr
		default:
			writer = os.Stdout
		}
		logFunc()
	}

	if LoggingToFile() {
		log.SetOutput(logFile)
		logFunc()
	}

	log.SetOutput(out)
}
