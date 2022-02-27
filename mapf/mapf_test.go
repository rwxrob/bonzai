package mapf_test

import (
	"os"

	"github.com/rwxrob/bonzai/each"
	"github.com/rwxrob/bonzai/fn"
	"github.com/rwxrob/bonzai/mapf"
)

func ExampleMarkDirs() {
	entries, _ := os.ReadDir("testdata/markdirs")
	each.Println(fn.Map(entries, mapf.MarkDirs))
	//Output:
	// dir1/
	// file1
}

func ExampleHashComment() {
	each.Println(fn.Map([]string{"foo", "bar"}, mapf.HashComment))
	// Output:
	// # foo
	// # bar
}

func ExampleEscSpace() {
	s := []string{"one here", "and another    one"}
	each.Println(fn.Map(s, mapf.EscSpace))
	// Output:
	// one\ here
	// and\ another\ \ \ \ one
}
