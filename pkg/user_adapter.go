package glice

import (
	"fmt"
)

var userCache = make(map[string]UserAdapter, 0)

// UserAdapter provides an interface for supporting an editor domain,
// e.g. github.org, bitbucket.org, etc.
type UserAdapter interface {
	GetName() string
	GetEmail() string
}

// GetUserAdapter returns an object that implements UserAdapter for a given domain
func GetUserAdapter(dep *Dependency) (ua UserAdapter, err error) {
	var ok bool
	if ua, ok = userCache[dep.Host]; ok {
		goto end
	}
	switch dep.Host {
	case "github.com":
		ua = GetGitUser()
	default:
		msg :=
			`repository hosts with domain '%s' are not yet supported. ` +
				`Support can be added in ./user_adapter.go"`
		err = fmt.Errorf(msg, dep.Host)
		goto end
	}
	userCache[dep.Host] = ua
end:
	return ua, err
}
