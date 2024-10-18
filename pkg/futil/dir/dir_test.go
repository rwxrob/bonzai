package dir_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/fs/dir"
	"github.com/rwxrob/fs/file"
)

func ExampleEntries() {
	list := dir.Entries("testdata")
	fmt.Println(list)
	// Output:
	// [testdata/adir testdata/emptyfiles testdata/file testdata/hiddenfiles testdata/nofiles testdata/notemptyfiles testdata/other]

}

func ExampleExists() {
	fmt.Println(dir.Exists("testdata/exists"))
	fmt.Println(dir.Exists("testdata"))
	// Output:
	// false
	// true
}

func ExampleName() {
	os.Chdir(`testdata`)
	fmt.Println(dir.Name())
	os.Chdir(`..`)
	// Output:
	// testdata
}

func ExampleHereOrAbove_here() {
	_dir, _ := os.Getwd()
	defer func() { os.Chdir(_dir) }()
	os.Chdir("testdata/adir")

	path, err := dir.HereOrAbove("anotherdir")
	if err != nil {
		fmt.Println(err)
	}
	d := strings.Split(path, string(filepath.Separator))
	fmt.Println(d[len(d)-2:])

	// Output:
	// [adir anotherdir]

}

func ExampleHereOrAbove_above() {
	_dir, _ := os.Getwd()
	defer func() { os.Chdir(_dir) }()
	os.Chdir("testdata/adir/anotherdir")

	path, err := dir.HereOrAbove("adir")
	if err != nil {
		fmt.Println(err)
	}
	d := strings.Split(path, string(filepath.Separator))
	fmt.Println(d[len(d)-2:])

	// Output:
	// [testdata adir]

}

/*
func ExampleAbs() {
	fmt.Println(dir.Abs())
	// Output:
	// ignored
}
*/

func ExampleIsEmpty_notexist() {
	fmt.Println(dir.IsEmpty(`testdata/notexist`))
	// Output:
	// false
}

func ExampleIsEmpty_nofiles() {
	fmt.Println(dir.IsEmpty(`testdata/nofiles`))
	// Output:
	// true
}

func ExampleIsEmpty_emptyfiles() {
	fmt.Println(dir.IsEmpty(`testdata/emptyfiles`))
	// Output:
	// true
}

func ExampleIsEmpty_notemptyfiles() {
	fmt.Println(file.Size(`testdata/notemptyfiles/README.md`))
	fmt.Println(dir.IsEmpty(`testdata/notemptyfiles`))
	// Output:
	// 5
	// false
}

func ExampleIsEmpty_hiddenfiles() {
	fmt.Println(file.Size(`testdata/hiddenfiles/.foo`))
	fmt.Println(dir.IsEmpty(`testdata/hiddenfiles`))
	// Output:
	// 5
	// false
}
