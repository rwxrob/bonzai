package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/redu"
)

// Combine is a [bonzai.Completer] that combines completions from
// multiple [bonzai.Completer]s. It does not handle duplicates.
type Combine []bonzai.Completer

func (completers Combine) Complete(args ...string) []string {
	var list []string
	for _, comp := range completers {
		list = append(list, comp.Complete(args...)...)
	}
	return redu.Unique(list)
}

func (completers Combine) SetCmd(a *bonzai.Cmd) { current = a }
func (completers Combine) Cmd() *bonzai.Cmd     { return current }

var (
	CmdsOpts               = Combine{Cmds, Opts}
	CmdsAliases            = Combine{Cmds, Aliases}
	CmdsAliasesOpts        = Combine{Cmds, Aliases, Opts}
	FileDirCmdsOpts        = Combine{FileDir, CmdsOpts}
	FileDirCmdsAliasesOpts = Combine{FileDir, CmdsAliasesOpts}
)
