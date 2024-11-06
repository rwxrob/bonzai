package comp

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/vars"
)

type _Vars struct{}

var Vars = new(_Vars)

// Complete takes a [*bonzai.Cmd] and then calls
func (_Vars) Complete(_ bonzai.Cmd, args ...string) (list []string) {
	if vars.Data == nil || len(args) == 0 {
		return
	}
	list, _ = vars.Data.KeysWithPrefix(args[0])
	return
}
