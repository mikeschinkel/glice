package gllicscan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

//goland:noinspection GoUnusedGlobalVariable
var GitLabLicenseScanningReportJsonExample = `
{
  "version": "2.1",
  "licenses": [
    {
      "id": "MPL-2.0",
      "name": "Mozilla Public License 2.0",
      "url": "https://opensource.org/licenses/MPL-2.0"
    }
  ],
  "dependencies": [
    {
      "name": "rhino",
      "version": "1.7.10",
      "package_manager": "maven",
      "path": "pom.xml",
      "licenses": [
        "MPL-2.0"
      ]
    }
  ]
}
`

const (
	SchemaVersion  = "2.1"
	ReportFilename = "gl-license-scanning-report.json"
)

var _ FilepathGetter = (*Report)(nil)

func NewReport(fp string) (rpt *Report) {
	return &Report{
		Filepath:     fp,
		Version:      SchemaVersion,
		Licenses:     make(Licenses, 0),
		Dependencies: make(Dependencies, 0),
	}
}

type Report struct {
	Filepath     string       `json:"-"`
	Version      string       `json:"version"`
	Licenses     Licenses     `json:"licenses"`
	Dependencies Dependencies `json:"dependencies"`
}

func (rpt *Report) Amend(adapter GitLabReportAdapter) {
	rpt.Licenses = rpt.Licenses.Amend(adapter.GetLicensesIDs())
	rpt.Dependencies = rpt.Dependencies.Amend(adapter.GetDependencyAdapters())
}

var ErrFileDoesNotExist = errors.New("file does not exist")
var ErrCannotReadFile = errors.New("cannot read file")
var ErrCannotUnmarshalJSON = errors.New("unmarshal JSON")

func LoadReport(fp string) (rpt *Report, err error) {
	var content []byte
	var exists bool

	rpt = NewReport(fp)
	exists, err = CheckFileExists(fp)
	if !exists {
		err = fmt.Errorf("%w,%s", ErrFileDoesNotExist, err.Error())
		goto end
	}
	content, err = ioutil.ReadFile(fp)
	if err != nil {
		err = fmt.Errorf("%w; %s", ErrCannotReadFile, err.Error())
		goto end
	}
	err = json.Unmarshal(content, rpt)
	if err != nil {
		err = fmt.Errorf("%w; %s", ErrCannotUnmarshalJSON, err.Error())
		goto end
	}
end:
	if err != nil {
		err = fmt.Errorf("cannot load report %s; %w", fp, err)
	}
	return rpt, err
}
func (rpt *Report) GetFilepath() string {
	return rpt.Filepath
}
func (rpt *Report) Save() (err error) {
	err = SaveJSONFile(rpt)
	if err != nil {
		err = fmt.Errorf("unable to save %s file; %w", ReportFilename, err)
		goto end
	}
end:
	return err
}
