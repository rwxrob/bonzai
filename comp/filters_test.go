package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/comp/filt"
)

func ExamplePipe_Complete() {
	cmd := bonzai.Cmd{
		Name: `foo`,
		Opts: `fooc|foobar|foobaz`,
		Comp: comp.Pipe{comp.Opts, filt.Prefix(`foo`)},
	}
	fmt.Println(cmd.Comp.Complete())
	fmt.Println(cmd.Comp.Complete(``))
	fmt.Println(cmd.Comp.Complete(`b`))

	// Output:
	// []
	// [fooc foobar foobaz]
	// [foobar foobaz]
}
