package glice

import (
	"os"
	"path/filepath"
)

const CacheSubDir = "glice"
const CacheFilename = "cache.json"

var cacheFilepath = filepath.Join(CacheDir(), CacheFilename)

func CacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		LogAndExit(exitCannotGetCacheDir,
			"Unable to get cache dir as %s",
			err.Error())
	}
	dir = filepath.Join(dir, CacheSubDir)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		LogAndExit(exitCannotCreateCacheDir,
			"Unable to get create cache subdir %s: %s",
			err.Error())
	}
	return dir
}

func CacheFilepath() string {
	return cacheFilepath
}

func SetCacheFilepath(fp string) {
	cacheFilepath = fp
}
