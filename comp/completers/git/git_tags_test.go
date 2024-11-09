package git_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"

	"github.com/rwxrob/bonzai/comp/completers/git"
)

func ExampleCompTags_Complete() {
	setupTags()
	defer teardownGitRepo()
	cmd := bonzai.Cmd{}
	fmt.Println(git.CompTags.Complete(cmd))
	fmt.Println(git.CompTags.Complete(cmd, ``))
	fmt.Println(git.CompTags.Complete(cmd, `tag-`))
	fmt.Println(git.CompTags.Complete(cmd, `tag-b`))
	fmt.Println(git.CompTags.Complete(cmd, `tag-f`))
	fmt.Println(git.CompTags.Complete(cmd, `tag-ba`))
	// Output:
	// []
	// [tag-bar tag-blah tag-foo]
	// [tag-bar tag-blah tag-foo]
	// [tag-bar tag-blah]
	// [tag-foo]
	// [tag-bar]
}

func setupTags() {
	setupGitRepo()
	safely(futil.Touch(`./foo`))
	run.Out(`git`, `add`, `.`)
	run.Out(`git`, `commit`, `-m`, `commit foo`)
	run.Out(`git`, `tag`, `tag-foo`)
	safely(futil.Touch(`./bar`))
	run.Out(`git`, `add`, `.`)
	run.Out(`git`, `commit`, `-m`, `commit bar`)
	run.Out(`git`, `tag`, `tag-bar`)
	safely(futil.Touch(`./blah`))
	run.Out(`git`, `add`, `.`)
	run.Out(`git`, `commit`, `-m`, `commit blah`)
	run.Out(`git`, `tag`, `tag-blah`)
}

func setupGitRepo() {
	safely(os.Mkdir(`./test_repo`, 0o755))
	safely(os.Chdir(`./test_repo`))
	run.Out(`git`, `init`)
	run.Out(`git`, `config`, `--local`, `user.email`, `test@test.com`)
	run.Out(`git`, `config`, `--local`, `user.name`, `test`)
}

func teardownGitRepo() {
	os.Chdir(`..`)
	safely(os.RemoveAll(`./test_repo`))
}

func safely(err error) {
	if err != nil {
		panic(err)
	}
}
