package glice

import (
	"fmt"
)

var editorCache = make(map[string]EditorGetter, 0)

// EditorGetter provides an interface for supporting an editor domain,
// e.g. github.org, bitbucket.org, etc.
type EditorGetter interface {
	GetEditor() *Editor
}

// GetEditorGetter returns an object that implements EditorGetter for a given domain
func GetEditorGetter(dep *Dependency) (eg EditorGetter, err error) {
	var ok bool
	if eg, ok = editorCache[dep.Host]; ok {
		goto end
	}
	switch dep.Host {
	case "github.com":
		eg = NewGit()
	default:
		msg :=
			`repository hosts with domain '%s' are not yet supported. ` +
				`Support can be added in ./editor_getter.go"`
		err = fmt.Errorf(msg, dep.Host)
		goto end
	}
	editorCache[dep.Host] = eg
end:
	return eg, err
}
