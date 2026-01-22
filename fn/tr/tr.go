package tr

import "github.com/BuddhiLW/bonzai/fn"

type Strings fn.Transformer[string, string]
type Ints fn.Transformer[int, int]
type Anys fn.Transformer[any, any]

// ------------------------------ Prefix ------------------------------

type Prefix struct {
	With string
}

func (p Prefix) Transform(in []string) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = p.With + v
	}
	return out
}
