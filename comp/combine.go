package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/redu"
)

type Combine []bonzai.Completer

// Complete calls Complete on all items in its list returning the
// resulting combined list (without removing duplicates).
func (completers Combine) Complete(an any, args ...string) []string {
	var list []string
	for _, comp := range completers {
		list = append(list, comp.Complete(an, args...)...)
	}
	return redu.Unique(list)
}

var (
	CmdsOpts        = Combine{Cmds, Opts}
	FileDirCmdsOpts = Combine{FileDir, CmdsOpts}
)
