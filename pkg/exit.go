package glice

// Not using iota here because these values should never change once set
// and iota makes it too easy to accidentally change them.
const (
	ExitUnexpectedError              = -1
	ExitAuditFoundDisallowedLicenses = 1
	ExitCannotGetWorkingDir          = 2
	ExitCannotGetCacheDir            = 3
	ExitCannotCreateCacheDir         = 4
	ExitCannotScanDependencies       = 5
	ExitCannotSaveFile               = 6
	ExitCannotStatFile               = 7
	ExitFileExistsCannotOverwrite    = 8
	ExitFileDoesNotExist             = 9
	ExitRepoInfoGetterIsNil          = 10
	ExitCannotGetRepositoryAdapter   = 11
	ExitCannotSetOptions             = 12
	ExitHostNotYetSupported          = 13
	ExitCannotWriteReport            = 14
	ExitCannotCreateFile             = 15
	ExitCannotGetReportWriterAdapter = 16
	ExitCannotReadFile               = 17
	ExitCannotUnmarshalJSON          = 18
)
