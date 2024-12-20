package inprops_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/persisters/inprops"
)

func ExampleSetup() {

	defer func() {
		os.Remove(`testdata/temp.props`)
	}()

	p := inprops.Persister{`testdata/temp.props`}

	exists := func(path string) bool {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return false
		}
		return err == nil
	}

	fmt.Printf("before: %v\n", exists(`testdata/temp.props`))
	if err := p.Setup(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("after: %v\n", exists(`testdata/temp.props`))

	// Output:
	// before: false
	// after: true
}

func ExampleSetup_useExisting() {

	p := inprops.Persister{`testdata/some.props`}

	exists := func(path string) bool {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return false
		}
		return err == nil
	}

	fmt.Printf("before: %v\n", exists(`testdata/some.props`))
	beforebuf, _ := os.ReadFile(`testdata/some.props`)
	if err := p.Setup(); err != nil {
		fmt.Println(err)
	}
	afterbuf, _ := os.ReadFile(`testdata/some.props`)
	fmt.Printf("after: %v\n", exists(`testdata/some.props`))

	if string(beforebuf) == string(afterbuf) {
		fmt.Println(`same text before and after`)
	}

	// Output:
	// before: true
	// after: true
	// same text before and after
}

func ExampleGet() {
	p := inprops.Persister{`testdata/some.props`}
	fmt.Println(p.Get(`some`))
	fmt.Println(p.Get(`another`))
	fmt.Printf("%q\n", p.Get(`not-exist`))
	// Output:
	// thing
	// one
	// ""
}

func ExampleSet() {

	defer func() {
		os.Remove(`testdata/temp.props`)
	}()

	p := inprops.Persister{`testdata/temp.props`}

	p.Set(`some`, `thing`)
	p.Set(`another`, `one`)
	out, _ := os.ReadFile(`testdata/temp.props`)
	fmt.Println(string(out))

	// Unordered Output:
	// some=thing
	// another=one
}
