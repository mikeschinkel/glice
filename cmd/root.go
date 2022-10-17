package cmd

import (
	"os"

	"github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

/*
Commands & Switches:
	--verbose
	--indirect
	--cache-file={cache_file}
	--path={repo_dir}
	--log
	--logfile
	init - Initialize glice.yaml for a directory
		--cache-file={cache_file}
		--path={repo_dir}
	scan - CI check
		--cache-file={cache_file}
		--path={repo_dir}
		--ttl={cache_ttl}
	report - Generate a license report
			- print - Print license report to stdout
			- write - Write license report to file
				--file={report_file}
  text - Write licenses to text files
		--output={output_dir}
	thank - Give thanks by starring repositories
*/

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "help",
	Short: "Glice inspects your Go-specific source dependencies for disallowed open-source licenses",
	Long: `Glice is a tool for inspecting your Go-specific source code dependencies to check for ` +
		`disallowed open-source licenses with functionality specifically defined to address ` +
		`Continuous Integration (CI) needs.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(initOptions)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var indirect bool
var verbose bool
var logOutput bool
var nocache bool
var logfile string
var source string
var cachefile string

func init() {
	pf := rootCmd.PersistentFlags()
	pf.BoolVar(&indirect, "indirect", false, "Include indirect dependencies")
	pf.BoolVar(&verbose, "verbose", false, "Generate verbose output")
	pf.BoolVar(&logOutput, "log", false, "Log output to default logging filepath.")
	pf.StringVar(&logfile, "logfile", "", "File to log output to.")
	pf.StringVar(&source, "source", glice.SourceDir(""), "Directory where go.mod for the repo to scan is located.")
	pf.StringVar(&cachefile, "cache-file", glice.CacheFilepath(), "Full filepath to the cachefile to create.")
	pf.BoolVar(&nocache, "nocache", false, "Disable use of caching")
}

func initOptions() {
	glice.SetOptions(&glice.Options{
		LogVerbosely:    verbose,
		IncludeIndirect: indirect,
		LogOuput:        logOutput,
		NoCache:         nocache,
		LogFilepath:     logfile,
		SourceDir:       source,
		CacheFilepath:   cachefile,
	})
}
