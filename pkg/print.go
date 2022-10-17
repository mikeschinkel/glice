package glice

import (
	"context"
	"io"
)

var (
	headerRow = []string{"Dependency", "RepoURL", "License"}
)

func Print(options *Options, writeTo io.Writer) error {
	o := options.Clone()
	o.OutputFormat = "table"
	o.OutputDestination = "stdout"
	return PrintTo(o, writeTo)
}

func PrintTo(options *Options, writeTo io.Writer) error {
	c := NewClient(options)
	err := c.Validate()
	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = ScanRepositories(ctx, options)
	if err != nil {
		return err
	}

	_ = c.Print(writeTo)
	return nil
}
