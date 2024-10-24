package comp

import bonzai "github.com/rwxrob/bonzai/pkg"

type Combine []bonzai.Completer

func (completers Combine) Complete(an any, args ...string) []string {
	var list []string
	for _, comp := range completers {
		list = append(list, comp.Complete(an, args...)...)
	}
	return list
}

var FileDirCmdsParams = Combine{FileDir, CmdsParams}
