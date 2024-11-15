package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
)

type aliases struct{}

// Aliases is a [bonzai.Completer] for all available [bonzai.Cmd.Aliases]
// including [bonzai.Cmd.Name].
var Aliases = new(aliases)

func (aliases) Complete(args ...string) []string {
	list := []string{}
	if len(args) == 0 || current == nil {
		return list
	}
	for _, c := range current.Aliases() {
		if len(c) > 0 {
			list = append(list, c)
		}
	}
	return filt.HasPrefix(list, args[0])
}

func (aliases) SetCmd(a *bonzai.Cmd) { current = a }
func (aliases) Cmd() *bonzai.Cmd     { return current }
