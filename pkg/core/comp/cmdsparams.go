package comp

import (
	"log"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/fn/filt"
)

type _CmdsParams struct{}

var CmdsParams = new(_CmdsParams)

// Returns all visible [Cmd.Commands] and [Cmd.Params] that match
// [Cmd.HasPrefix] passed the first argument. See [bonzai.Completer].
func (_CmdsParams) Complete(an any, args ...string) []string {

	x, is := an.(*bonzai.Cmd)
	if !is {
		log.Printf(`%v is a %T not *bonzai.Cmd`, an, an)
		return []string{}
	}

	// not sure we've completed the command name itself yet
	if len(args) == 0 {
		return []string{x.Name}
	}

	list := []string{}

	// commands
	for _, c := range x.Commands {
		if c.IsHidden() {
			continue
		}
		list = append(list, c.Name)
	}

	// params
	list = append(list, x.ParamsSlice()...)

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}
