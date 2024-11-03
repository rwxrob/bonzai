package comp

import (
	"log"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/maps"
)

type _Vars struct{}

var Vars = new(_Vars)

// Complete takes a [*bonzai.Cmd] and then calls
func (_Vars) Complete(an any, args ...string) (list []string) {
	x, is := an.(*bonzai.Cmd)
	if !is {
		log.Printf(`%v is a %T not *bonzai.Cmd`, an, an)
		return
	}
	if x.Vars == nil {
		return
	}
	list = maps.KeysWithPrefix(x.Vars, args[0])
	// TODO if we have a full match for a key, complete the value as well
	return
}
