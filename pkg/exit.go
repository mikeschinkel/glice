package glice

const (
	ExitUnexpectedError              = -1
	ExitAuditFoundDisallowedLicenses = 1
	ExitCannotGetWorkingDir          = 3
	ExitCannotGetCacheDir            = 4
	ExitCannotCreateCacheDir         = 5
	ExitCannotScanDependencies       = 6
	ExitCannotSaveFile               = 7
	ExitCannotStatFile               = 8
	ExitFileExistsCannotOverwrite    = 9
	ExitFileDoesNotExist             = 10
	ExitRepoInfoGetterIsNil          = 11
	ExitCannotGetRepositoryAdapter   = 12
	ExitCannotSetOptions             = 13
	ExitHostNotYetSupported          = 14
	ExitCannotWriteReport            = 15
	ExitCannotCreateFile             = 16
	ExitCannotGetReportWriterAdapter = 17
	ExitCannotReadFile               = 18
	ExitCannotUnmarshalJSON          = 19
)
