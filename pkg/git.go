package glice

import (
	"os/exec"
	"strings"
)

type Git struct{}

func NewGit() *Git {
	return &Git{}
}

var editor *Editor

func (g *Git) GetEditor() *Editor {
	var cmd *exec.Cmd
	var name, email []byte
	var err error

	if editor != nil {
		goto end
	}

	cmd = exec.Command("git", "config", "--default", defaultName, "user.name")
	name, err = cmd.CombinedOutput()
	if err != nil {
		Warnf("Failed when calling `git config`; %s", err.Error())
		name = []byte(defaultName)
	}
	cmd = exec.Command("git", "config", "--default", defaultEmail, "user.email")
	email, err = cmd.CombinedOutput()
	if err != nil {
		Warnf("Failed when calling `git config`; %s", err.Error())
		email = []byte(defaultName)
	}
	editor = &Editor{
		Name:  strings.TrimSpace(string(name)),
		Email: strings.TrimSpace(string(email)),
	}
end:
	return editor
}
