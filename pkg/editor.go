package glice

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
)

type Editors []*Editor
type Editor struct {
	Name      string `yaml:"name"`
	Email     string `yaml:"email"`
	Reference string `yaml:"ref"`
}

const numProperties = 3

var _ yaml.Marshaler = (*Editor)(nil)
var _ yaml.Unmarshaler = (*Editor)(nil)

func (e *Editor) MarshalYAML() (interface{}, error) {
	if e.Reference == "" {
		e.Reference = UpToN(e.Email, '@', 1)
	}
	editor := fmt.Sprintf("&%s %s <%s>", e.Reference, e.Name, e.Email)
	return editor, nil
}

var regexParseEditor = regexp.MustCompile(`^\s*&(\S+)\s+(.+)\s+<([^>]+)>\s*$`)

var errMsg = "editor value '%s' is incomplete, or is not formatted correctly"

func (e *Editor) UnmarshalYAML(node *yaml.Node) (err error) {
	var editor []string

	segments := regexParseEditor.FindAllStringSubmatch(node.Value, -1)
	if segments == nil || len(segments) == 0 {
		err = fmt.Errorf(errMsg, node.Value)
		goto end
	}
	editor = segments[0]
	if len(editor) <= numProperties {
		err = errors.New(errMsg)
	}
	*e = Editor{}
	if len(editor) > 1 {
		e.Name = editor[1]
	}
	if len(editor) > 2 {
		e.Email = editor[2]
	}
	if len(editor) > 3 {
		e.Reference = editor[3]
	}
	if e.Reference == "" {
		err = AppendError(err, "editor reference must not be empty")
	}
	if e.Name == "" {
		err = AppendError(err, "editor name must not be empty")
	}
	if e.Email == "" {
		err = AppendError(err, "editor email must not be empty")
	}
end:
	if err != nil {
		err = fmt.Errorf("%s; %s", err.Error(),
			"should be formatted as '&reference FirstName LastName <email@example.com>'")
	}
	return err
}
