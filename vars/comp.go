package vars

type cmp struct{}

var Comp = new(cmp)

// Complete takes a [*bonzai.Cmd] and then calls
func (cmp) Complete(args ...string) (list []string) {
	if Data == nil || len(args) == 0 {
		return
	}
	list, _ = Data.KeysWithPrefix(args[0])
	return
}
