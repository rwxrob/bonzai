package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
)

type _Aliases struct{}

// Aliases is a [bonzai.Completer] that returns all [bonzai.Cmd.Alias]
// including the [bonzai.Cmd.Name].
var Aliases = new(_Aliases)

// Complete returns all [bonzai.Cmd.Alias] and [bonzai.Cmd.Name] that
// match [filt.HasPrefix] for first argument. See [bonzai.Completer].
func (_Aliases) Complete(x bonzai.Cmd, args ...string) []string {
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
