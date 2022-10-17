package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/ribice/glice/v3/pkg"
)

var (
	fileWrite = flag.Bool("f", false, "Write all licenses to files")
	indirect  = flag.Bool("i", false, "Gets indirect modules as well")
	path      = flag.String("p", "", `Path of desired directory to be scanned with Glice (e.g. "github.com/ribice/glice/v3")`)
	thanks    = flag.Bool("t", false, "Stars dependent repos. Needs GITHUB_API_KEY env variable to work")
	verbose   = flag.Bool("v", false, "Adds verbose logging")
	format    = flag.String("fmt", "table", "Output format [table | json | csv]")
	notext    = flag.Bool("notext", false, "Allows ignoring the capture of license text")
	output    = flag.String("o", "stdout", "Output location [stdout | file]")
)

func initOptions() *glice.Options {
	flag.Parse()
	o := glice.GetOptions()
	o.WriteFile = *fileWrite
	o.IncludeIndirect = *indirect
	o.SourceDir = *path
	//o.GiveThanks = *thanks
	o.LogVerbosely = *verbose
	o.OutputFormat = *format
	o.NoCaptureLicenseText = *notext
	o.OutputDestination = *output

	if o.SourceDir == "" {
		cf, err := os.Getwd()
		checkErr(err)
		o.SourceDir = cf
	}

	if !o.LogVerbosely {
		log.SetOutput(io.Discard)
		log.SetFlags(0)
	}

	return o
}

func main() {

	o := initOptions()

	client := glice.NewClient(o)
	cl, err := client, client.Validate()
	checkErr(err)

	//_, err = cl.ParseDependencies(o)
	//checkErr(err)

	checkErr(cl.GenerateOutput())

}

func checkErr(err error) {
	if err != nil {
		log.SetOutput(os.Stderr)
		if el, ok := err.(glice.ErrorList); ok {
			for _, err := range el {
				log.Print(err.Error())
			}
			log.Fatalf("%d errors occured", len(el))
		}
		log.Fatal(err.Error())

	}
}
