package mapf_test

import (
	"github.com/rwxrob/bonzai/fn"
	"github.com/rwxrob/bonzai/mapf"
)

func ExampleHashComment() {
	fn.A[string]{"foo", "bar"}.M(mapf.HashComment).P()
	// Output:
	// # foo
	// # bar
}
