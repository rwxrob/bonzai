package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
)

type aliases struct{}

// Aliases is a [bonzai.Completer] for all available [bonzai.Cmd.Alias]
// including [bonzai.Cmd.Name].
var Aliases = new(aliases)

func (aliases) Complete(x bonzai.Cmd, args ...string) []string {
	list := []string{}
	if len(args) == 0 {
		return list
	}
	for _, c := range x.Names() {
		if len(c) > 0 {
			list = append(list, c)
		}
	}
	return filt.HasPrefix(list, args[0])
}
