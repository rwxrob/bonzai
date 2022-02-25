package util_test

import (
	"os"

	"github.com/rwxrob/bonzai/filter"
	"github.com/rwxrob/bonzai/util"
)

func ExampleFiles() {
	filter.Println(util.Files("testdata/files"))
	// Output:
	// bar
	// blah
	// dir1/
	// dir2/
	// dir3/
	// foo
	// other
	// some
}

func ExampleFiles_empty() {
	os.Chdir("testdata/files")
	defer os.Chdir("../..")
	filter.Println(util.Files(""))
	// Output:
	// bar
	// blah
	// dir1/
	// dir2/
	// dir3/
	// foo
	// other
	// some
}

func ExampleFilesWith() {
	filter.Println(util.FilesWith("testdata/files", "b"))
	// Output:
	// bar
	// blah
}
