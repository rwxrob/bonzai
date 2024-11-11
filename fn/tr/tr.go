package tr

type Interface[I any, O any] interface {
	Transform(in I) O
}

type String Interface[string, string]
type Strings Interface[[]string, []string]
type Int Interface[int, int]
type Ints Interface[[]int, []int]
type Any Interface[any, any]
type Anys Interface[[]any, []any]

// ------------------------------ Prefix ------------------------------

type Prefix struct {
	V string
}

func (p Prefix) Transform(a string) string { return p.V + a }

// ----------------------------- PrefixAll ----------------------------

type PrefixAll struct {
	V string
}

func (p PrefixAll) Transform(a []string) []string {
	out := make([]string, len(a))
	for n, it := range a {
		out[n] = p.V + it
	}
	return out
}
