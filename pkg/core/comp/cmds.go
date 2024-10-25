package comp

import (
	"log"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/fn/filt"
)

type _Cmds struct{}

var Cmds = new(_Cmds)

// Complete returns all visible [Cmd.Cmds] that match [futil.HasPrefix]
// for arg[0] . See [bonzai.Completer].
func (_Cmds) Complete(an any, args ...string) []string {

	x, is := an.(*bonzai.Cmd)
	if !is {
		log.Printf(`%v is a %T not *bonzai.Cmd`, an, an)
		return []string{}
	}

	list := []string{}
	for _, c := range x.Cmds {
		if c.IsHide() {
			continue
		}
		list = append(list, c.Name)
	}

	return filt.HasPrefix(list, args[0])
}
