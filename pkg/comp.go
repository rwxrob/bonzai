package bonzai

import (
	"github.com/rwxrob/bonzai/pkg/core/ds/set"
	"github.com/rwxrob/bonzai/pkg/core/fn/filt"
)

type defcomp struct{}

var DefComp = new(defcomp)

// Complete resolves completion as follows:
//
//  1. If leaf has Comp function, delegate to it
//
//  2. If leaf has no arguments, return all Cmds and Params
//
//  3. If first argument is the name of a Command return it only even
//     if in the Hidden list
//
//  4. Otherwise, return every Command or Param that is not in the
//     Hidden list and HasPrefix matching the first arg
//
// See bonzai.Completer.
func (defcomp) Complete(x *Cmd, args ...string) []string {

	// if has completer, delegate
	if c := x.Comp; c != nil {
		return c.Complete(x, args...)
	}

	// not sure we've completed the command name itself yet
	if len(args) == 0 {
		return []string{x.Name}
	}

	// build list of visible commands and params
	list := []string{}
	list = append(list, x.CommandNames()...)
	list = append(list, x.AliasesSlice()...)
	list = append(list, x.ParamsSlice()...)
	list = set.Minus[string, string](list, x.HiddenSlice())

	//log.Print(list)

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}
