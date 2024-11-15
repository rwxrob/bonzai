package git

import (
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/run"
)

type tags struct{}

// CompTags is a [bonzai.Completer] that completes for git tags.
var CompTags bonzai.Completer = tags{}

func (t tags) Complete(args ...string) []string {
	list := make([]string, 0)
	if len(args) == 0 {
		return list
	}
	tags := run.Out("git", "tag", "-l", "--no-column")
	if len(tags) == 0 {
		return list
	}
	for _, tag := range strings.Split(tags, "\n") {
		tag := strings.TrimSpace(tag)
		if len(tag) == 0 {
			continue
		}
		if strings.HasPrefix(tag, args[0]) {
			list = append(list, tag)
		}
	}
	return list
}
