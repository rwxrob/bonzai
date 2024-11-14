package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
	"github.com/rwxrob/bonzai/fn/redu"
	"github.com/rwxrob/bonzai/fn/tr"
)

type Combine []any

func (ops Combine) Complete(args ...string) []string {
	var list []string
	for _, this := range ops {
		switch v := this.(type) {
		case bonzai.Completer:
			list = append(list, v.Complete(args...)...)
		case filt.Strings:
			list = v.Filter(list)
		case tr.Strings:
			list = v.Transform(list)
		case redu.Strings:
			list = v.Reduce(list)
		}
	}
	return redu.Unique(list)
}

func (Combine) SetCmd(a *bonzai.Cmd) { current = a }
func (Combine) Cmd() *bonzai.Cmd     { return current }
