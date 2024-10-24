package comp

import bonzai "github.com/rwxrob/bonzai/pkg"

type Combine []bonzai.Completer

// Complete calls Complete on all items in its list returning the
// resulting combined list (without removing duplicates).
func (completers Combine) Complete(an any, args ...string) []string {
	var list []string
	for _, comp := range completers {
		list = append(list, comp.Complete(an, args...)...)
	}
	return list
}

var (
	CmdsParams        = Combine{Cmds, Params}
	FileDirCmdsParams = Combine{FileDir, CmdsParams}
)
