package glice_test

import (
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
			IncludeIndirect:      true,
			SourceDir:            SourceDirectory,
			VerbosityLevel:       glice.WarnLevel,
			OutputFormat:         "json",
			NoCaptureLicenseText: true,
			OutputDestination:    "file",
		},
	}
)

func TestYAMLFileInit(t *testing.T) {
	var err error
	var file *os.File

	err = os.Setenv("GITHUB_API_KEY", GitHubAPIKey)
	if err != nil {
		t.Errorf("failed to set GitHub API key in environment: %s", err.Error())
	}
	file, err = os.Create(LogFilepath)
	if err != nil {
		t.Errorf("failed to create logfile %s: %s", LogFilepath, err.Error())
	}
	log.SetOutput(file)
	for _, options := range yamlFileInitTests {
		t.Run(options.SourceDir, func(t *testing.T) {
			yf := glice.NewYAMLFile(options.SourceDir)
			yf.Generated = glice.Timestamp()[:10]
			yf.AllowedLicenses = glice.DefaultAllowedLicenses
			yf.Editors = glice.Editors{
				{Name: "Mike Schinkel", Email: "mschinkel-ctr@singlestore.com"},
			}
			yf.Overrides = glice.Overrides{
				{
					DependencyImport: "github.com/Masterminds/squirrel",
					LicenseID:        "MIT",
					VerifiedBy:       "*mschinkel-ctr",
					LastVerified:     glice.Timestamp()[:10],
					Notes:            "Verification not real, done by unit test",
				},
				{
					DependencyImport: "github.com/miekg/dns",
					LicenseID:        "BSD-3-Clause",
					VerifiedBy:       "*mschinkel-ctr",
					LastVerified:     glice.Timestamp()[:10],
					Notes:            "Verification not real, done by unit test",
				},
			}
			yf.Dependencies, err = glice.ScanDependencies(options)
			if err != nil {
				t.Errorf("failed to parse dependencies: %s", err.Error())
			}
			yf.Filepath = glice.YAMLFilepath(TestDataDir)
			err = yf.Init()
			if err != nil {
				t.Errorf("failed to create YAML file %s: %s",
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

func TestYAMLFileLoad(t *testing.T) {
	for _, test := range yamlFileLoadTests {
		t.Run(test.Directory, func(t *testing.T) {
			yf, err := glice.LoadYAMLFile(test.Directory)
			if err != nil {
				t.Errorf("failed to load YAML file %s; %s",
					yf.Filepath,
					err.Error())
			}
		})
	}
}
