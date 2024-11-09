package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
)

type _Opts struct{}

var Opts = new(_Opts)

// Complete returns all [Cmd.Opts] that match [filt.HasPrefix] passed
// the first argument. See [bonzai.Completer].
func (_Opts) Complete(x bonzai.Cmd, args ...string) []string {
	list := x.OptsSlice()

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}
