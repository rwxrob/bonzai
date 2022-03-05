package util_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/each"
	"github.com/rwxrob/bonzai/util"
)

func ExampleFiles() {
	each.Println(util.Files("testdata/files"))
	// Output:
	// testdata/files/bar
	// testdata/files/blah
	// testdata/files/dir1/
	// testdata/files/foo
	// testdata/files/other
	// testdata/files/some
}

func ExampleFiles_spaces() {
	each.Println(util.Files("testdata/files/dir1"))
	// Output:
	// testdata/files/dir1/some\ thing
}

func ExampleFiles_empty() {
	os.Chdir("testdata/files")
	defer os.Chdir("../..")
	each.Println(util.Files(""))
	// Output:
	// bar
	// blah
	// dir1/
	// foo
	// other
	// some
}

func ExampleFiles_not_Directory() {
	fmt.Println(util.Files("none"))
	// Output:
	// []
}
