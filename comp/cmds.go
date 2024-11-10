package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
)

type cmds struct{}

// Cmds is a [bonzai.Completer] for all available [bonzai.Cmd.Cmds]. It
// excludes hidden commands.
var Cmds = new(cmds)

func (cmds) Complete(x bonzai.Cmd, args ...string) []string {
	list := []string{}
	if len(args) == 0 {
		return list
	}
	for _, c := range x.Cmds {
		if c.IsHidden() {
			continue
		}
		list = append(list, c.Name)
	}
	return filt.HasPrefix(list, args[0])
}
