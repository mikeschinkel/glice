package glice_test

import (
	"context"
	"fmt"
	"github.com/ribice/glice/v3/pkg"
	"log"
	"os"
	"testing"
)

const (
	SourceDirectory = "/home/mschinkel-ctr/Projects/helios/go/src/freya"
	TestDataDir     = "test/data"
	TestLogDir      = "test/log"
)

var (
	GitHubAPIKey = os.Getenv("GITHUB_API_KEY")

	LogFilepath = glice.SourceDir(fmt.Sprintf("%s/test.log", TestLogDir))

	yamlFileInitTests = []*glice.Options{
		{
			DirectOnly:           false,
			SourceDir:            SourceDirectory,
			VerbosityLevel:       glice.WarnLevel,
			OutputFormat:         "json",
			NoCaptureLicenseText: true,
			OutputDestination:    "file",
		},
	}
)

func TestProjectFileInit(t *testing.T) {
	var err error
	var file *os.File

	ctx := context.Background()

	if GitHubAPIKey == "" {
		t.Fatal("The GITHUB_API_KEY envionment variable has not been set")
	}
	file, err = os.Create(LogFilepath)
	if err != nil {
		t.Errorf("failed to create logfile %s: %s", LogFilepath, err.Error())
	}
	log.SetOutput(file)
	for _, options := range yamlFileInitTests {
		t.Run(options.SourceDir, func(t *testing.T) {
			pf := glice.NewProjectFile(options.SourceDir)
			pf.Generated = glice.Timestamp()[:10]
			pf.AllowedLicenses = glice.DefaultAllowedLicenses
			pf.Editors = glice.Editors{
				{Name: "Mike Schinkel", Email: "mschinkel-ctr@singlestore.com"},
			}
			pf.Overrides = glice.Overrides{
				{
					Import:       "github.com/Masterminds/squirrel",
					LicenseID:    "MIT",
					VerifiedBy:   "*mschinkel-ctr",
					LastVerified: glice.Timestamp()[:10],
					Notes:        "Verification not real, done by unit test",
				},
				{
					Import:       "github.com/miekg/dns",
					LicenseID:    "BSD-3-Clause",
					VerifiedBy:   "*mschinkel-ctr",
					LastVerified: glice.Timestamp()[:10],
					Notes:        "Verification not real, done by unit test",
				},
			}
			pf.Dependencies, err = glice.ScanDependencies(ctx, options)
			if err != nil {
				t.Errorf("failed to parse dependencies: %s", err.Error())
			}
			pf.Filepath = glice.GetProjectFilepath(TestDataDir)
			err = pf.Save()
			if err != nil {
				t.Errorf("failed to create `glice.yaml` file %s: %s",
					options.SourceDir,
					err.Error())
			}
		})
	}
}

var yamlFileLoadTests = []struct {
	Directory string
}{
	{Directory: glice.SourceDir(TestDataDir)},
}

func TestProjectFileLoad(t *testing.T) {
	for _, test := range yamlFileLoadTests {
		t.Run(test.Directory, func(t *testing.T) {
			pf, err := glice.LoadProjectFile(test.Directory)
			if err != nil {
				t.Errorf("failed to load `glice.yaml` file %s; %s",
					pf.Filepath,
					err.Error())
			}
		})
	}
}
