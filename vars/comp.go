package vars

import (
	"github.com/rwxrob/bonzai"
)

type _Comp struct{}

var Comp = new(_Comp)

// Complete takes a [*bonzai.Cmd] and then calls
func (_Comp) Complete(_ bonzai.Cmd, args ...string) (list []string) {
	if Data == nil || len(args) == 0 {
		return
	}
	list, _ = Data.KeysWithPrefix(args[0])
	return
}
