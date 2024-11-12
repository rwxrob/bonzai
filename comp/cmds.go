package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
)

type cmds struct{}

// Cmds is a [bonzai.Completer] for all available [bonzai.Cmd.Cmds]. It
// excludes hidden commands.
var Cmds = new(cmds)

func (cmds) Complete(args ...string) []string {
	list := []string{}
	if len(args) == 0 || current == nil {
		return list
	}
	for _, c := range current.Cmds {
		if c.IsHidden() {
			continue
		}
		list = append(list, c.Name)
	}
	return filt.HasPrefix(list, args[0])
}

func (cmds) SetCmd(a *bonzai.Cmd) { current = a }
func (cmds) Cmd() *bonzai.Cmd     { return current }
