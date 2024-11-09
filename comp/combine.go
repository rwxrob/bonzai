package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/redu"
)

type Combine []bonzai.Completer

// Complete calls Complete on all items in its list returning the
// resulting combined list (without removing duplicates).
func (completers Combine) Complete(
	x bonzai.Cmd,
	args ...string,
) []string {
	var list []string
	for _, comp := range completers {
		list = append(list, comp.Complete(x, args...)...)
	}
	return redu.Unique(list)
}

var (
	CmdsOpts               = Combine{Cmds, Opts}
	CmdsAliases            = Combine{Cmds, Aliases}
	CmdsAliasesOpts        = Combine{Cmds, Aliases, Opts}
	FileDirCmdsOpts        = Combine{FileDir, CmdsOpts}
	FileDirCmdsAliasesOpts = Combine{FileDir, CmdsAliasesOpts}
)
