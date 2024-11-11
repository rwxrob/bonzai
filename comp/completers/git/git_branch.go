package git

import (
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/run"
)

type _Branches struct{}

// CompBranches is a [bonzai.Completer] that completes for git branches.
var CompBranches = _Branches{}

func (t _Branches) Complete(_ bonzai.Cmd, args ...string) []string {
	list := make([]string, 0)
	if len(args) == 0 {
		return list
	}
	branches := run.Out(
		"git",
		"branch",
		"-l",
		"--no-column",
		"--format",
		"%(refname:short)",
	)
	if len(branches) == 0 {
		return list
	}
	for _, branch := range strings.Split(branches, "\n") {
		tag := strings.TrimSpace(branch)
		if len(tag) == 0 {
			continue
		}
		if strings.HasPrefix(tag, args[0]) {
			list = append(list, tag)
		}
	}
	return list
}
