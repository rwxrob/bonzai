package comp

import (
	"log"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/fn/filt"
)

type _Params struct{}

var Params = new(_Params)

// Returns all  [Cmd.Params] that match [futil.HasPrefix] passed the
// first argument. See [bonzai.Completer].
func (_Params) Complete(an any, args ...string) []string {

	x, is := an.(*bonzai.Cmd)
	if !is {
		log.Printf(`%v is a %T not *bonzai.Cmd`, an, an)
		return []string{}
	}

	list := x.ParamsSlice()

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}
