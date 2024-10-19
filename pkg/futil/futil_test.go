package futil_test

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai/pkg/futil"
)

//go:embed all:testdata/testfs
var testfs embed.FS

func ExampleExtractEmbed_confirm_Default_Read() {

	// go:embed all:testdata/testfs
	// var testfs embed.FS

	foo, err := os.Open("testdata/testfs/foo")
	if err != nil {
		fmt.Println(err)
	}
	info, err := foo.Stat()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info.Mode())

	// Output:
	// -rwxr-xr-x
}

func ExampleExtractEmbed() {

	// go:embed all:testdata/testfs
	// var testfs embed.FS
	defer os.RemoveAll("testdata/testfsout")

	if err := futil.ExtractEmbed(testfs,
		"testdata/testfs", "testdata/testfsout"); err != nil {
		fmt.Println(err)
	}

	stuff := []string{
		`foo`,
		`_notignored`,
		`.secret`,
		`dir`,
		`dir/README.md`,
	}

	for _, i := range stuff {
		f, err := os.Stat(filepath.Join("testdata/testfsout", i))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%v %v\n", i, f.Mode())
	}

	// Output:
	// foo -rw-------
	// _notignored -rw-------
	// .secret -rw-------
	// dir drwx------
	// dir/README.md -rw-------

}

/*
func ExampleIsDir() {
	fmt.Println(futil.IsDir("testdata"))
	fmt.Println(futil.IsDir("testdata/fail"))
	// Output:
	// true
	// false
}

func ExampleDupPerms() {
	os.Mkdir(`testdata/some`, 0000)
	defer os.Remove(`testdata/some`)
	stats, _ := os.Stat(`testdata/orig`)
	fmt.Println(stats.Mode())
	stats, _ = os.Stat(`testdata/some`)
	fmt.Println(stats.Mode())
	err := futil.DupPerms(`testdata/orig`, `testdata/some`)
	if err != nil {
		fmt.Println(err)
	}
	stats, _ = os.Stat(`testdata/some`)
	fmt.Println(stats.Mode())
	// Output:
	// drw-------
	// d---------
	// drw-------
}
*/

func ExampleExists() {
	fmt.Println(futil.Exists("testdata/exists"))
	fmt.Println(futil.Exists("testdata"))
	// Output:
	// true
	// true
}

func ExampleModTime() {
	fmt.Println(futil.ModTime("testdata/file").IsZero())
	fmt.Println(futil.ModTime("testdata/none"))
	// Output:
	// true
	// 0001-01-01 00:00:00 +0000 UTC
}

func ExampleHereOrAbove_here() {
	dir, _ := os.Getwd()
	defer func() { os.Chdir(dir) }()
	os.Chdir("testdata/adir")

	path, err := futil.HereOrAbove("afile")
	if err != nil {
		fmt.Println(err)
	}
	d := strings.Split(path, string(filepath.Separator))
	fmt.Println(d[len(d)-2:])

	// Output:
	// [adir afile]

}

func ExampleHereOrAbove_above() {
	dir, _ := os.Getwd()
	defer func() { os.Chdir(dir) }()
	os.Chdir("testdata/adir")

	path, err := futil.HereOrAbove("anotherfile")
	if err != nil {
		fmt.Println(err)
	}
	d := strings.Split(path, string(filepath.Separator))
	fmt.Println(d[len(d)-2:])

	// Output:
	// [testdata anotherfile]

}

func ExampleIntDirs() {

	dirs, low, high := futil.IntDirs("testdata/ints")

	fmt.Println(low)
	fmt.Println(high)

	for _, d := range dirs {
		fmt.Printf("%v ", filepath.Base(d.Path))
	}

	// Output:
	// 2
	// 10
	// 10 2 3 4 5 6 7 8 9

}

/*
func ExampleIsosecModTime() {
	dirs, _, _ := futil.IntDirs("testdata/ints")
	fmt.Println(futil.IsosecModTime(dirs[1].Info))
	// Output:
	// 20221117200334
}
*/

/*
func ExampleLatestChange() {
	path, _ := dir.LatestChange("testdata")
	fmt.Println(path)
	// Output:
	// whichever
}
*/

/*

// remaining only work for specific users and systems

func ExampleTilde2Home_good() {
	path := `~/some/path`
	fmt.Println(futil.Tilde2Home(path))
	// Output:
	// /home/rwxrob/some/path
}

func ExampleTilde2Home_no_Tilde() {
	path := `/some/path`
	fmt.Println(futil.Tilde2Home(path))
	// Output:
	// /some/path
}

*/

/*
func ExamplePreserve() {
	file.Touch(`testdata/preserve`)
	name, err := futil.Preserve(`testdata/preserve`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(name)
	// Output:
	// ignored
}
*/

/*
func ExampleRevertIfMissing() {
	file.Overwrite(`testdata/preserve`, `original`)
	name, err := futil.Preserve(`testdata/preserve`)
	os.Remove(`testdata/preserve`)
	defer futil.RevertIfMissing(`testdata/preserve`, name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(name)
	// Output:
	// ignored
}
*/

func ExampleNameIsInt() {
	fmt.Println(futil.NameIsInt(`/some/1`))
	fmt.Println(futil.NameIsInt(`100`))
	fmt.Println(futil.NameIsInt(`-100`))
	fmt.Println(futil.NameIsInt(`/some/-1`))
	fmt.Println(futil.NameIsInt(`/some/`))
	// Output:
	// true
	// true
	// false
	// false
	// false
}
