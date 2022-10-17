package glice

const LogFilename = "glice.log"

var logFilepath = SourceDir(LogFilename)

func LogFilepath() string {
	return logFilepath
}
func SetLogFilepath(fp string) {
	logFilepath = fp
}
