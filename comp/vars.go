package comp

import (
	"github.com/rwxrob/bonzai"
)

type _Vars struct{}

var Vars = new(_Vars)

// Complete takes a [*bonzai.Cmd] and then calls
func (_Vars) Complete(_ bonzai.Cmd, args ...string) (list []string) {
	if bonzai.Vars == nil || len(args) == 0 {
		return
	}
	list, _ = bonzai.Vars.KeysWithPrefix(args[0])
	return
}
