package glice

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// SourceDir returns the current working directory as a string
func SourceDir() string {
	return GetSourceDir("")
}

// GetSourceDir returns the current working directory as a string
// with the path passed appended.
func GetSourceDir(path string) string {
	opt := GetOptions()
	if opt.SourceDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			Failf(ExitCannotGetWorkingDir,
				"Unable to get current working directory: %s",
				err.Error())
		}
		opt.SourceDir = wd
	}
	if path == "" {
		path = opt.SourceDir
	} else {
		path = filepath.Join(opt.SourceDir, path)
	}
	return path
}

// FileExists returns true of the file represented by the passed filepath exists
func FileExists(fp string) (exists bool) {
	_, err := os.Stat(fp)
	if errors.Is(err, fs.ErrNotExist) {
		goto end
	}
	if err != nil {
		Failf(ExitCannotStatFile,
			"Unable to check existence of %s: %s",
			fp,
			err.Error())
	}
	exists = true
end:
	return exists
}

func LoadYAMLFile(fp string, obj interface{}) (fg FilepathGetter, err error) {
	var b []byte

	if !FileExists(fp) {
		err = fmt.Errorf("unable to find %s; %w", fp, ErrLoadableYAMLFile)
		goto end
	}
	b, err = os.ReadFile(fp)
	if err != nil {
		err = fmt.Errorf("unable to read %s; %w", fp, err)
		goto end
	}
	err = yaml.Unmarshal(b, obj)
	if err != nil {
		err = fmt.Errorf("unable to unmashal %s; %w", fp, err)
		goto end
	}
end:
	return obj.(FilepathGetter), err
}

func SaveYAMLFile(fg FilepathGetter) (err error) {
	var f *os.File
	var b []byte

	fp := fg.GetFilepath()
	f, err = os.Create(fp)
	if err != nil {
		err = fmt.Errorf("unable to open file '%s'; %w", fp, err)
		goto end
	}
	defer MustClose(f)

	b, err = yaml.Marshal(fg)
	if err != nil {
		err = fmt.Errorf("unable to encode to %s; %w", fp, err)
		goto end
	}
	_, err = f.Write(b)
	if err != nil {
		err = fmt.Errorf("unable to write to '%s'; %w", fp, err)
		goto end
	}

end:
	return err
}

//goland:noinspection GoUnusedConst
const (
	RetainOriginalFile = true
	DeleteOriginalFile = false
)

// BackupFile creates a backup of the file passed by adding a ".bak"
// or ".<n>.bak" extension while maintaining all prior backups.
func BackupFile(fp string, retainOriginal bool) (bfs []string, err error) {
	var bf string

	if !FileExists(fp) {
		goto end
	}
	bf = fmt.Sprintf("%s.bak", fp)
	bfs = []string{fp, bf}
	for {
		if !FileExists(bf) {
			break
		}
		bf = fmt.Sprintf("%s.%d.bak", fp, len(bfs))
		bfs = append(bfs, bf)
	}
	for fc := len(bfs) - 1; fc > 0; fc-- {
		err = os.Rename(bfs[fc-1], bfs[fc])
		if err != nil {
			err = fmt.Errorf("unable to backup %s; %w", fp, err)
		}
	}
	if !retainOriginal {
		goto end
	}

	err = CopyFile(bf, fp)
	if err != nil {
		err = fmt.Errorf("unable to retain original file '%s' during backup; %w", fp, err)
		goto end
	}

end:
	return bfs, err
}

// CopyFile copies a file from source to destination, returning an error if applicable
func CopyFile(src, dst string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		err = fmt.Errorf("unable to copy FROM file '%s'; %w", src, err)
		goto end
	}

	err = ioutil.WriteFile(dst, content, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("unable to copy TO file '%s'; %w", dst, err)
		goto end
	}

end:
	return err
}

func ReplaceFileExtension(filename string, newExt FileExtension) string {
	ext := path.Ext(filename)
	return filename[:len(filename)-len(ext)+1] + string(newExt)
}
