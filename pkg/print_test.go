package glice_test

//import (
//	"bytes"
//	"github.com/ribice/glice/v3"
//	"os"
//	"path/filepath"
//	"reflect"
//	"testing"
//)
//
//func wd() string {
//	d, _ := os.Getwd()
//	return d
//}
//
//var gliceDeps = []string{"github.com/fatih/color", "github.com/google/go-github",
//	"github.com/keighl/metabolize", "github.com/olekukonko/tablewriter",
//	"golang.org/x/mod", "golang.org/x/oauth2"}
//
//func TestGetOtherRepo(t *testing.T) {
//	if glice.ResolveVanityURL("golang.org/x/net/context/ctxhttp").url != "https://go.googlesource.com/net" {
//		t.Error("Wrong URL")
//	}
//}
//
//func TestClient_ParseDependencies(t *testing.T) {
//	tests := map[string]struct {
//		path            string
//		includeIndirect bool
//		thanks          bool
//		wantRepos       []string
//		wantErr         bool
//	}{
//		"thanks without api key": {
//			thanks:  true,
//			wantErr: true,
//		},
//		"Invalid path": {
//			path:    "invalid",
//			wantErr: true,
//		},
//		"Valid path": {
//			path:      wd(),
//			wantRepos: gliceDeps,
//		},
//	}
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			c := &glice.Client{path: tt.path, format: "table", output: "stdout"}
//			if err := c.ParseDependencies(tt.includeIndirect, tt.thanks); (err != nil) != tt.wantErr {
//				t.Errorf("ParseDependencies() error = %v, wantErr %v", err, tt.wantErr)
//			}
//			if len(c.dependencies) != len(tt.wantRepos) {
//				t.Error("expected number of repos and urlParsed not the same")
//			}
//			if len(tt.wantRepos) > 1 {
//				var gotRepos []string
//				for _, v := range c.dependencies {
//					gotRepos = append(gotRepos, v.Import)
//				}
//				if !reflect.DeepEqual(tt.wantRepos, gotRepos) {
//					t.Error("got and want repos do not match")
//				}
//			}
//		})
//	}
//}
//
//func TestPrint(t *testing.T) {
//	tests := map[string]struct {
//		path            string
//		wantWriteOutput bool
//		wantErr         bool
//	}{
//		"invalid path": {
//			path:    "invalid",
//			wantErr: true,
//		},
//		"valid path": {
//			path:            wd(),
//			wantWriteOutput: true,
//		},
//	}
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			writeTo := &bytes.Buffer{}
//			err := Print(tt.path, false, writeTo)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if (writeTo.String() != "") != tt.wantWriteOutput {
//				t.Error("wantWriteOutput and gotOutput do not match")
//			}
//		})
//	}
//}
//
//func TestPrintTo(t *testing.T) {
//	tests := map[string]struct {
//		path            string
//		format          string
//		wantWriteOutput bool
//		wantErr         bool
//	}{
//		"invalid path": {
//			path:    "invalid",
//			wantErr: true,
//		},
//		"json format": {
//			path:            wd(),
//			wantWriteOutput: true,
//			format:          "json",
//		},
//		"csv format": {
//			path:            wd(),
//			wantWriteOutput: true,
//			format:          "csv",
//		},
//		"valid path": {
//			path:            wd(),
//			wantWriteOutput: true,
//			format:          "table",
//		},
//	}
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			writeTo := &bytes.Buffer{}
//			err := PrintTo(tt.path, tt.format, "stdout", false, writeTo)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if (writeTo.String() != "") != tt.wantWriteOutput {
//				t.Error("wantWriteOutput and gotOutput do not match")
//			}
//		})
//	}
//}
//
//func TestClient_Print(t *testing.T) {
//	tests := map[string]struct {
//		dependencies []*glice.Dependency
//		wantOutput   bool
//	}{
//		"without dependencies": {},
//		"with dependencies": {
//			dependencies: []*glice.Dependency{{Import: "Glice", url: "github.com/ribice/glice", Shortname: "MIT"}},
//			wantOutput:   true,
//		},
//	}
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			c := &glice.Client{dependencies: tt.dependencies, format: "table", output: "stdout"}
//			output := &bytes.Buffer{}
//			c.Print(output)
//			if (output.String() != "") != tt.wantOutput {
//				t.Error("wantOutput and gotOutput do not match")
//			}
//		})
//	}
//}
//
//func TestClient_WriteLicensesToFile(t *testing.T) {
//	tests := map[string]struct {
//		dependencies   []*glice.Dependency
//		wantErr        bool
//		wantOutputFile bool
//	}{
//		"no dependencies": {},
//		"a dependency with invalid license text": {
//			dependencies: []*glice.Dependency{{
//				Author:  "ribice",
//				Project: "glice",
//				Text:    "license-text",
//			}},
//			wantErr: true},
//		"a dependency without license text": {
//			dependencies: []*glice.Dependency{{
//				Author:  "ribice",
//				Project: "glice",
//			}},
//		},
//		"valid dependency": {
//			dependencies: []*glice.Dependency{{
//				Author:  "ribice",
//				Project: "glice",
//				Text:    "bGljZW5zZS10ZXh0",
//			}},
//			wantOutputFile: true},
//	}
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			c := &glice.Client{dependencies: tt.dependencies, format: "table", output: "stdout"}
//			err := c.WriteLicensesToFile()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if tt.wantOutputFile {
//				licensePath := filepath.Join(wd(), "licenses", "ribice-glice-license.MD")
//
//				if _, err := os.Stat(licensePath); err != nil {
//					if os.IsNotExist(err) {
//						t.Errorf("License file is missing, but should be there")
//					}
//				}
//
//				mustRemove(t, licensePath)
//				mustRemove(t, filepath.Join(wd(), "licenses"))
//			}
//		})
//	}
//
//}
//
//func TestListRepositories(t *testing.T) {
//	_, err := ScanRepositories("path", false)
//	if err == nil {
//		t.Errorf("expected err, got: %v", err)
//	}
//
//	wd, err := os.Getwd()
//	if err != nil {
//		t.Error(err)
//	}
//
//	repos, err := ScanRepositories(wd, false)
//	if err != nil {
//		t.Error(err)
//	}
//
//	var gotNames []string
//	for _, r := range repos {
//		gotNames = append(gotNames, r.Import)
//	}
//
//	if !reflect.DeepEqual(gliceDeps, gotNames) {
//		t.Errorf("listRepositories() = %v, want %v", gotNames, gliceDeps)
//	}
//}
//
//func TestNewClient(t *testing.T) {
//	tests := map[string]struct {
//		path    string
//		output  string
//		format  string
//		wantErr bool
//	}{
//		"invalid format": {
//			format:  "invalid",
//			wantErr: true,
//		},
//		"invalid output": {
//			format:  "csv",
//			output:  "invalid",
//			wantErr: true,
//		},
//		"invalid path": {
//			format:  "csv",
//			output:  "stdout",
//			path:    "invalid",
//			wantErr: true,
//		},
//		"invalid  path": {
//			format: "csv",
//			output: "stdout",
//			path:   wd(),
//		},
//	}
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			_, err := NewClient(tt.path, tt.format, tt.output)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//		})
//	}
//}
//
//func TestGetRepository(t *testing.T) {
//	tests := map[string]struct {
//		module string
//		want   *glice.Dependency
//	}{
//		"github.com/ribice": {
//			module: "github.com/ribice",
//			want:   &glice.Dependency{Import: "github.com/ribice"},
//		},
//		"github.com/ribice/glice": {
//			module: "github.com/ribice/glice",
//			want:   &glice.Dependency{Import: "github.com/ribice/glice", url: "https://github.com/ribice/glice", Host: "github.com", Author: "ribice", Project: "glice"},
//		},
//		"gopkg.in/ribice": {
//			module: "gopkg.in/ribice",
//			want:   &glice.Dependency{Import: "gopkg.in/ribice"},
//		},
//		"gopkg.in/ribice/glice": {
//			module: "gopkg.in/ribice/glice",
//			want:   &glice.Dependency{Import: "gopkg.in/ribice/glice", url: "https://github.com/ribice/glice", Host: "github.com", Author: "ribice", Project: "glice"},
//		},
//		"fmt": {
//			module: "fmt",
//			want:   &glice.Dependency{Import: "fmt"},
//		},
//	}
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			if got := glice.GetRepository(tt.module); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetRepository() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func mustRemove(t *testing.T, path string) {
//	err := os.Remove(path)
//	if err != nil {
//		t.Errorf("Failed to remove '%s'", path)
//	}
//}
