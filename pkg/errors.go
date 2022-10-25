package glice

import (
	"errors"
)

var (
	// ErrNoGoModFile is returned when path doesn't contain go.mod file
	ErrNoGoModFile = errors.New("no go.mod file present")

	ErrLoadableYAMLFile = errors.New("loadable YAML file does not exist")

	// ErrNoAPIKey is returned when thanks flag is enabled without providing GITHUB_API_KEY env variable
	ErrNoAPIKey = errors.New("the GITHUB_API_KEY environment variable is empty")

	ErrCannotLogin = errors.New("host cannot login likely because of lacking credentials")

	ErrRequestPrefixInstead = errors.New("prefix is a subset of import; request prefix instead")
)
