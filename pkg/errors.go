package glice

import "errors"

var (
	// ErrNoGoMod is returned when path doesn't contain go.mod file
	ErrNoGoMod = errors.New("no go.mod file present")

	// ErrNoAPIKey is returned when thanks flag is enabled without providing GITHUB_API_KEY env variable
	ErrNoAPIKey = errors.New("cannot use thanks feature without github api key")
)
