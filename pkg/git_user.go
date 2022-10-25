package glice

import (
	"os/exec"
	"strings"
)

var _ UserAdapter = (*GitUser)(nil)

type GitUser struct {
	Name  string
	Email string
}

func (gu *GitUser) GetName() string {
	return gu.Name
}

func (gu *GitUser) GetEmail() string {
	return gu.Email
}

var user *GitUser

func GetGitUser() UserAdapter {
	var cmd *exec.Cmd
	var name, email []byte
	var err error

	if user != nil {
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
	user = &GitUser{
		Name:  strings.TrimSpace(string(name)),
		Email: strings.TrimSpace(string(email)),
	}
end:
	return user
}
