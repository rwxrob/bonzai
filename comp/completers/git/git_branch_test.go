package git_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"

	"github.com/rwxrob/bonzai/comp/completers/git"
)

func ExampleCompBranches_Complete() {
	setupBranches()
	defer teardownGitRepo()
	fmt.Println(git.CompBranches.Complete())
	fmt.Println(git.CompBranches.Complete(``))
	fmt.Println(git.CompBranches.Complete(``))
	fmt.Println(git.CompBranches.Complete(`b`))
	fmt.Println(git.CompBranches.Complete(`f`))
	fmt.Println(git.CompBranches.Complete(`ba`))
	// Output:
	// []
	// [bar blah foo main]
	// [bar blah foo main]
	// [bar blah]
	// [foo]
	// [bar]
}

func setupBranches() {
	setupGitRepo()
	safely(futil.Touch(`./foo`))
	run.Out(`git`, `add`, `.`)
	run.Out(`git`, `commit`, `-m`, `commit foo`)
	run.Out(`git`, `switch`, `-c`, `foo`)
	run.Out(`git`, `switch`, `-c`, `bar`)
	run.Out(`git`, `switch`, `-c`, `blah`)
}
