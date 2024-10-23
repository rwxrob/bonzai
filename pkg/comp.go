package bonzai

import (
	"github.com/rwxrob/bonzai/pkg/core/fn/filt"
)

type defcomp struct{}

var DefComp = new(defcomp)

// Complete resolves completion as follows:
//
//  1. If leaf has Comp function, delegate to it
//
//  2. If leaf has no arguments, return all visible Commands and Params
//
//  4. Otherwise, return all visible Commands or Params that
//     [Cmd.HasPrefix] matching the first arg
//
// See [bonzai.Completer].
func (defcomp) Complete(x *Cmd, args ...string) []string {

	// if has completer, delegate
	if c := x.Comp; c != nil {
		return c.Complete(x, args...)
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
