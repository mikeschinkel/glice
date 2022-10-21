package glice

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
)

type Editors []*Editor
type EditorMap map[string]*Editor
type Editor struct {
	Name   string
	Email  string
	Anchor string
}

var (
	defaultName  = "Username Goes Here"
	defaultEmail = "email-alias@example.com"
)

const numProperties = 2

var _ yaml.Marshaler = (*Editor)(nil)
var _ yaml.Unmarshaler = (*Editor)(nil)

func (em EditorMap) ToEditors() Editors {
	editors := make(Editors, len(em))
	index := 0
	for _, ed := range em {
		editors[index] = ed
		index++
	}
	return editors
}

func (e *Editor) MarshalYAML() (interface{}, error) {
	return e.String(), nil
}

var regexParseEditor = regexp.MustCompile(`^\s*(.+)\s+<([^>]+)>\s*$`)

var errMsg = "editor value '%s' is incomplete, or is not formatted correctly"

func (e *Editor) UnmarshalYAML(node *yaml.Node) (err error) {
	var editor []string

	segments := regexParseEditor.FindAllStringSubmatch(node.Value, 1)
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
	if e.Name == "" {
		err = AppendError(err, "editor name must not be empty")
	}
	if e.Email == "" {
		err = AppendError(err, "editor email must not be empty")
	}
	e.Anchor = node.Anchor
end:
	if err != nil {
		err = fmt.Errorf("%s; %s", err.Error(),
			"should be formatted as '&anchor FirstName LastName <email@example.com>'")
	}
	return err
}

func (e *Editor) Alias() string {
	e.ensureAnchor()
	return fmt.Sprintf("*%s", e.Anchor)
}

func (e *Editor) String() string {
	e.ensureAnchor()
	return fmt.Sprintf("&%s %s <%s>", e.Anchor, e.Name, e.Email)
}

func (e *Editor) ensureAnchor() {
	if e.Anchor == "" {
		e.Anchor = UpToN(e.Email, '@', 1)
	}
}
