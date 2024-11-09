package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
)

type _Cmds struct{}

// Cmds is a [bonzai.Completer] that returns all non-hidden
// [bonzai.Cmd.Cmds]
var Cmds = new(_Cmds)

// Complete returns all visible [Cmd.Cmds] that match [futil.HasPrefix]
// for arg[0] . See [bonzai.Completer].
func (_Cmds) Complete(x bonzai.Cmd, args ...string) []string {
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
