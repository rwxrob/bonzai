package futil_test

import (
	"fmt"
	"log"
	"net/http"
	ht "net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai/pkg/futil"
)

func ExampleTouch_create() {
	fmt.Println(NotExists("testdata/foo"))
	futil.Touch("testdata/foo")
	fmt.Println(Exists("testdata/foo"))
	os.Remove("testdata/foo")
	// Output:
	// true
	// true
}

func ExampleTouch_update() {

	// first create it and capture the time as a string
	futil.Touch("testdata/tmpfile")
	u1 := ModTime("testdata/tmpfile")
	log.Print(u1)

	// touch it and capture the new time
	futil.Touch("testdata/tmpfile")
	u2 := ModTime("testdata/tmpfile")
	log.Print(u2)

	// check that they are not equiv
	fmt.Println(u1 == u2)

	// Output:
	// false
}

func ExampleFetch() {

	// serve get
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `random file content`)
		})
	svr := ht.NewServer(handler)
	defer svr.Close()
	defer os.Remove(`testdata/file`)

	// not found
	handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		})
	notfound := ht.NewServer(handler)
	defer notfound.Close()

	if err := futil.Fetch(svr.URL, `testdata/file`); err != nil {
		fmt.Println(err)
	}

	it, _ := os.ReadFile(`testdata/file`)
	fmt.Println(string(it))

	if err := futil.Fetch(notfound.URL, `testdata/file`); err != nil {
		fmt.Println(err)
	}

	// Output:
	// random file content
	// 404 Not Found
}

func ExampleReplace() {

	// serve get
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `something random`)
	})
	svr := ht.NewServer(handler)
	defer svr.Close()

	// create a file to replace
	os.Create(`testdata/replaceme`)
	defer os.Remove(`testdata/replaceme`)
	os.Chmod(`testdata/replaceme`, 0400)

	// show info about control file
	info, err := os.Stat(`testdata/replaceme`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info.Mode())
	fmt.Println(info.Size())

	// replace it with local url
	futil.Replace(`testdata/replaceme`, svr.URL)

	// check that it is new
	info, err = os.Stat(`testdata/replaceme`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info.Mode())
	fmt.Println(info.Size())

	// Output:
	// -r--------
	// 0
	// -r--------
	// 16

}

func ExampleExists() {
	fmt.Println(futil.Exists("testdata/exists"))
	fmt.Println(futil.Exists("testdata"))
	// Output:
	// true
	// false
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

func ExampleHead() {

	lines, err := futil.Head(`testdata/headtail`, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [one two]
}

func ExampleHead_over() {

	lines, err := futil.Head(`testdata/headtail`, 20)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [one two three four five]
}

func ExampleTail() {

	lines, err := futil.Tail(`testdata/headtail`, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [four five]
}

func ExampleTail_over() {

	lines, err := futil.Tail(`testdata/headtail`, 20)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [one two three four five]
}

func ExampleTail_negative() {

	lines, err := futil.Tail(`testdata/headtail`, -2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [three four five]
}

/*
func ExampleRepaceAllString() {
	err := futil.ReplaceAllString(`testdata/headtail`, `three`, `THREE`)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// ignored
}
*/

func ExampleFindString() {
	str, err := futil.FindString(`testdata/headtail`, `thre+`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
	// Output:
	// three
}

func ExampleOverwrite() {
	err := futil.Overwrite(`testdata/overwritten`, `hello`)
	defer os.Remove(`testdata/overwritten`)
	if err != nil {
		fmt.Println(err)
	}
	futil.Cat(`testdata/overwritten`)
	// Output:
	// hello
}

func ExampleIsEmpty() {
	fmt.Println(futil.IsEmpty(`testdata/overwritten`))
	fmt.Println(futil.IsEmpty(`testdata/ovewritten`))
	futil.Touch(`testdata/emptyfile`)
	fmt.Println(futil.IsEmpty(`testdata/emptyfile`))
	// Output:
	// false
	// false
	// true
}

func ExampleSize() {
	fmt.Println(futil.Size(`testdata/headtail`))
	// Output:
	// 24
}

func ExampleField() {
	fmt.Println(futil.Field(`testdata/fieldtest`, 2))
	// Output:
	// [foo bar baz]
}
