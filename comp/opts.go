package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
)

type _Opts struct{}

// Opts is a [bonzai.Completer] for all available [bonzai.Cmd.Opts].
var Opts = new(_Opts)

func (_Opts) Complete(x bonzai.Cmd, args ...string) []string {
	list := x.OptsSlice()

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}
