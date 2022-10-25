package gllicscan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

func UnionString(s1, s2 []string) (s3 []string) {
	var n, i int

	m := make(map[string]struct{})
	for _, item := range s1 {
		m[item] = struct{}{}
	}
	for _, s := range s2 {
		if _, ok := m[s]; ok {
			continue
		}
		n++
	}
	if n == 0 {
		s3 = s1
		goto end
	}
	s3 = make([]string, n)
	copy(s3, s1)
	i = len(s1)
	for _, s := range s2 {
		if _, ok := m[s]; ok {
			continue
		}
		s3[i] = s
		i++
	}
end:
	return s3
}

func SaveJSONFile(fg FilepathGetter) (err error) {
	var b []byte
	fp := fg.GetFilepath()
	b, err = json.MarshalIndent(fg, "", "\t")
	if err != nil {
		err = fmt.Errorf("unable to marshal JSON to %s; %w", fp, err)
		goto end
	}
	err = ioutil.WriteFile(fp, b, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("unable to write to '%s'; %w", fp, err)
		goto end
	}
end:
	return err
}

var ErrCannotStatFile = errors.New("cannot stat file")

// CheckFileExists checks for file existence, and returns an error if it does not exist.
func CheckFileExists(fp string) (exists bool, err error) {
	_, err = os.Stat(fp)
	if errors.Is(err, fs.ErrNotExist) {
		goto end
	}
	if err != nil {
		err = fmt.Errorf("unable to check existence of %s: %w; %s",
			fp,
			ErrCannotStatFile,
			err)
		goto end
	}
	exists = true
end:
	return exists, err
}
