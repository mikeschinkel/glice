package glice

import "errors"

var (
	// ErrNoGoModFile is returned when path doesn't contain go.mod file
	ErrNoGoModFile = errors.New("no go.mod file present")

	ErrLoadableYAMLFile = errors.New("loadable YAML file does not exist")

	// ErrNoAPIKey is returned when thanks flag is enabled without providing GITHUB_API_KEY env variable
	ErrNoAPIKey = errors.New("cannot use thanks feature without github api key")
)
