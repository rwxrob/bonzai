package util_test

import (
	"github.com/rwxrob/bonzai/filter"
	"github.com/rwxrob/bonzai/util"
)

func ExampleFiles() {
	filter.Println(util.Files("testdata/files"))
	// Output:
	// bar
	// blah
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
