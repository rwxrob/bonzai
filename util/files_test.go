package util_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/loop"
	"github.com/rwxrob/bonzai/util"
)

func ExampleFiles() {
	loop.Println(util.Files("testdata/files"))
	// Output:
	// testdata/files/bar
	// testdata/files/blah
	// testdata/files/dir1/
	// testdata/files/dir2/
	// testdata/files/dir3/
	// testdata/files/foo
	// testdata/files/other
	// testdata/files/some
}

func ExampleFiles_empty() {
	os.Chdir("testdata/files")
	defer os.Chdir("../..")
	loop.Println(util.Files(""))
	// Output:
	// ./bar
	// ./blah
	// ./dir1/
	// ./dir2/
	// ./dir3/
	// ./foo
	// ./other
	// ./some
}

func ExampleFilesWith() {
	loop.Println(util.FilesWith("testdata/files", "b"))
	// Output:
	// testdata/files/bar
	// testdata/files/blah
}

func ExampleFiles_not_Directory() {
	fmt.Println(util.Files("none"))
	// Output:
	// []
}
