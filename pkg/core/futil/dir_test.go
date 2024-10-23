package futil_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/pkg/core/futil"
)

func ExampleEntries() {
	list := futil.DirEntries("testdata")
	fmt.Println(list)
	// Output:
	// [testdata/adir testdata/anotherfile testdata/emptyfile testdata/emptyfiles testdata/exists testdata/fieldtest testdata/headtail testdata/hiddenfiles testdata/ints testdata/nofiles testdata/notemptyfiles testdata/other testdata/preserve testdata/testfs testdata/tmpfile]

}

func ExampleDirName() {
	os.Chdir(`testdata`)
	fmt.Println(futil.DirName())
	os.Chdir(`..`)
	// Output:
	// testdata
}

func ExampleDirIsEmpty_notexist() {
	fmt.Println(futil.DirIsEmpty(`testdata/notexist`))
	// Output:
	// false
}

/*
// impossible to test while under git management
func ExampleDirIsEmpty_nofiles() {
	fmt.Println(futil.DirIsEmpty(`testdata/nofiles`))
	// Output:
	// true
}
*/

func ExampleDirIsEmpty_emptyfiles() {
	fmt.Println(futil.DirIsEmpty(`testdata/emptyfiles`))
	// Output:
	// true
}

func ExampleDirIsEmpty_notemptyfiles() {
	fmt.Println(futil.FileSize(`testdata/notemptyfiles/README.md`))
	fmt.Println(futil.DirIsEmpty(`testdata/notemptyfiles`))
	// Output:
	// 5
	// false
}

func ExampleDirIsEmpty_hiddenfiles() {
	fmt.Println(futil.FileSize(`testdata/hiddenfiles/.foo`))
	fmt.Println(futil.DirIsEmpty(`testdata/hiddenfiles`))
	// Output:
	// 5
	// false
}
