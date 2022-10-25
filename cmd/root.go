package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

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

var direct bool
var verbose int
var logOutput bool
var nocache bool
var captureLic bool
var logfile string
var source string
var cachefile string

func init() {
	pf := rootCmd.PersistentFlags()
	pf.BoolVar(&direct, "direct-only", false, "Exclude direct dependencies")
	pf.IntVar(&verbose, "verbose", glice.NoteLevel, fmt.Sprintf("Specify a verbosity level: %s", glice.ValidVerbositiesString))
	pf.Lookup("verbose").NoOptDefVal = strconv.Itoa(glice.InfoLevel)
	pf.BoolVar(&logOutput, "log", false, "Log output to default logging filepath.")
	pf.StringVar(&logfile, "logfile", "", "File to log output to.")
	pf.StringVar(&source, "source", glice.SourceDir(), "Source directory where go.mod.")
	pf.StringVar(&cachefile, "cache-file", glice.CacheFilepath(), "Full filepath to the cachefile to create.")
	pf.BoolVar(&nocache, "nocache", false, "Disable use of caching")
	pf.BoolVar(&captureLic, "capture-license", false, "Download license from host while processing (slower)")

	rootCmd.MarkFlagsMutuallyExclusive("nocache", "cache-file")
}

// initOptions create a glice.Options object and sets its values
// from the command line flags.
func initOptions() {
	glice.SetOptions(&glice.Options{
		VerbosityLevel: verbose,
		DirectOnly:     direct,
		LogOutput:      logOutput,
		NoCache:        nocache,
		LogFilepath:    logfile,
		SourceDir:      source,
		CacheFilepath:  cachefile,
		CaptureLicense: captureLic,
	})
}

// addTTLFlag allows multiple Commands to share the TTL Flag.
func addTTLFlag(cmd *cobra.Command) {
	cmd.Flags().Int("ttl", 24*60*60, "Time-to-Live for data in the cache file allowing recently audited dependencies to be skipped")
}
