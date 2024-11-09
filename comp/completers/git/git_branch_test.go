package git_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"

	"github.com/rwxrob/bonzai/run"

	"github.com/rwxrob/bonzai/comp/completers/git"
)

func ExampleCompBranches_Complete() {
	setupBranches()
	defer teardownGitRepo()
	cmd := bonzai.Cmd{}
	fmt.Println(git.CompBranches.Complete(cmd))
	fmt.Println(git.CompBranches.Complete(cmd, ``))
	fmt.Println(git.CompBranches.Complete(cmd, ``))
	fmt.Println(git.CompBranches.Complete(cmd, `b`))
	fmt.Println(git.CompBranches.Complete(cmd, `f`))
	fmt.Println(git.CompBranches.Complete(cmd, `ba`))
	// Output:
	// []
	// [bar blah foo]
	// [bar blah foo]
	// [bar blah]
	// [foo]
	// [bar]
}

func setupBranches() {
	setupGitRepo()
	run.Out(`git`, `switch`, `-c`, `foo`)
	run.Out(`git`, `switch`, `-c`, `bar`)
	run.Out(`git`, `switch`, `-c`, `blah`)
}
