/*
Package compcmd is a completion driver for Bonzai command trees and
fulfills the bonzai.Completer package interface. See Complete method for
details.
*/
package compcmd

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/pkg/fn/filt"
	"github.com/rwxrob/bonzai/pkg/set/text/set"
)

// New returns a private struct that fulfills the bonzai.Completer
// interface. See Complete method for details.
func New() *comp { return new(comp) }

type comp struct{}

// Complete resolves completion as follows:
//
//  1. If leaf has Comp function, delegate to it
//
//  2. If leaf has no arguments, return all Commands and Params
//
//  3. If first argument is the name of a Command return it only even
//     if in the Hidden list
//
//  4. Otherwise, return every Command or Param that is not in the
//     Hidden list and HasPrefix matching the first arg
//
// See bonzai.Completer.
func (comp) Complete(x bonzai.Command, args ...string) []string {

	// if has completer, delegate
	if c := x.GetComp(); c != nil {
		return c.Complete(x, args...)
	}

	// not sure we've completed the command name itself yet
	if len(args) == 0 {
		return []string{x.GetName()}
	}

	// build list of visible commands and params
	list := []string{}
	list = append(list, x.GetCommandNames()...)
	list = append(list, x.GetParams()...)
	list = append(list, x.GetShortcuts()...)
	list = set.Minus[string, string](list, x.GetHidden())

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}
