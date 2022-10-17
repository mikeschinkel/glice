package glice

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"os"
	"path/filepath"
)

type Client struct {
	dependencies Dependencies
	path         string
	format       string
	output       string
	NoText       bool
	validated    bool
}

var outputExtension = map[string]string{
	"table": "txt",
	"json":  "json",
	"csv":   "csv",
}

func NewClient(options *Options) *Client {
	SetOptions(options)
	return &Client{
		path:   options.SourceDir,
		format: options.OutputFormat,
		output: options.OutputDestination}
}

func (c *Client) GenerateOutput() (err error) {
	options := GetOptions()
	switch options.OutputDestination {
	case "stdout":
		err = c.Print(os.Stdout)
	case "file":
		fp := fmt.Sprintf("dependencies.%s", outputExtension[options.OutputFormat])
		f, err := os.Create(fp)
		if err != nil {
			err = fmt.Errorf("unable to create dependencies file %s; %w", fp, err)
			goto end
		}
		err = c.Print(f)
		if err != nil {
			err = fmt.Errorf("error occurred when attempting to print dependency report; %w", err)
			goto end
		}
		err = f.Close()
		if err != nil {
			err = fmt.Errorf("unable to close dependencies file %s; %w", fp, err)
			goto end
		}
	}
	if !options.WriteFile {
		goto end
	}
	err = c.WriteLicensesToFile()
	if err != nil {
		err = fmt.Errorf("unable to write licenses to individual files; %w", err)
		goto end
	}
end:
	return err
}

func (c *Client) Validate() (err error) {
	if _, ok := validFormats[c.format]; !ok {
		err = fmt.Errorf("invalid format provided (%s) - allowed ones are [table, json, csv]",
			c.format)
		goto end
	}

	if _, ok := validOutputs[c.output]; !ok {
		err = fmt.Errorf("invalid output provided (%s) - allowed ones are [stdout, file]",
			c.output)
		goto end
	}

	if !ModFileExists(c.path) {
		err = ErrNoGoMod
		goto end
	}

	c.validated = true
end:
	return err
}

func (c *Client) Print(writeTo io.Writer) error {
	if len(c.dependencies) < 1 {
		return nil
	}

	switch c.format {
	case "table":
		tw := tablewriter.NewWriter(writeTo)
		tw.SetHeader(headerRow)
		//for _, d := range c.dependencies {
		//	tw.Append([]string{d.Import, color.BlueString(d.url), d.Shortname})
		//}
		tw.Render()
	case "json":
		return json.NewEncoder(writeTo).Encode(c.dependencies)
	case "csv":
		csvW := csv.NewWriter(writeTo)
		defer csvW.Flush()
		err := csvW.Write(headerRow)
		if err != nil {
			return err
		}
		//for _, d := range c.dependencies {
		//	//err = csvW.Write([]string{d.Project, d.url, d.License})
		//	if err != nil {
		//		return err
		//	}
		//}
		return csvW.Error()
	}

	// shouldn't be possible to get this error
	return fmt.Errorf("invalid output provided (%s) - allowed ones are [stdout, json, csv]", c.output)
}

func (c *Client) WriteLicensesToFile() (err error) {
	if len(c.dependencies) < 1 {
		goto end
	}

	_ = os.Mkdir("licenses", os.ModePerm)

	for _, d := range c.dependencies {
		text := d.GetLicenseText()
		if text == "" {
			continue
		}

		var dec []byte
		dec, err = base64.StdEncoding.DecodeString(text)
		if err != nil {
			err = fmt.Errorf("unable to decode license text for '%s'; %w", d.Import, err)
			goto end
		}

		var f *os.File
		file := filepath.Join(c.path, "licenses", fmt.Sprintf("%s-%s-license.MD", d.Author, d.Project))
		f, err = os.Create(file)
		if err != nil {
			err = fmt.Errorf("unable to create file %s; %w", file, err)
			goto end
		}

		if _, err = f.Write(dec); err != nil {
			err = fmt.Errorf("unable to write to file %s; %w", file, err)
			goto end
		}

		if err = f.Sync(); err != nil {
			err = fmt.Errorf("unable to synchronize file %s; %w", file, err)
			goto end
		}

		if err := f.Close(); err != nil {
			err = fmt.Errorf("unable to close file %s; %w", file, err)
			goto end
		}
	}
end:
	return err
}
