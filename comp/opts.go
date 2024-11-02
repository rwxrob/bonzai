package comp

import (
	"log"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/fn/filt"
)

type _Opts struct{}

var Opts = new(_Opts)

// Complete returns all [Cmd.Opts] that match [futil.HasPrefix] passed the
// first argument. See [bonzai.Completer].
func (_Opts) Complete(an any, args ...string) []string {

	x, is := an.(*bonzai.Cmd)
	if !is {
		log.Printf(`%v is a %T not *bonzai.Cmd`, an, an)
		return []string{}
	}

	list := x.OptsSlice()

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}
