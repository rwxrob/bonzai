package inyaml_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/persisters/inyaml"
)

func ExampleSetup() {

	defer func() {
		os.Remove(`testdata/temp.yaml`)
	}()

	p := inyaml.Persister{`testdata/temp.yaml`}

	exists := func(path string) bool {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return false
		}
		return err == nil
	}

	fmt.Printf("before: %v\n", exists(`testdata/temp.yaml`))
	if err := p.Setup(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("after: %v\n", exists(`testdata/temp.yaml`))

	// Output:
	// before: false
	// after: true
}

func ExampleSetup_useExisting() {

	p := inyaml.Persister{`testdata/some.yaml`}

	exists := func(path string) bool {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return false
		}
		return err == nil
	}

	fmt.Printf("before: %v\n", exists(`testdata/some.yaml`))
	beforebuf, _ := os.ReadFile(`testdata/some.yaml`)
	if err := p.Setup(); err != nil {
		fmt.Println(err)
	}
	afterbuf, _ := os.ReadFile(`testdata/some.yaml`)
	fmt.Printf("after: %v\n", exists(`testdata/some.yaml`))

	if string(beforebuf) == string(afterbuf) {
		fmt.Println(`same text before and after`)
	}

	// Output:
	// before: true
	// after: true
	// same text before and after
}

func ExampleGet() {
	p := inyaml.Persister{`testdata/some.yaml`}
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
		os.Remove(`testdata/temp.yaml`)
	}()

	p := inyaml.Persister{`testdata/temp.yaml`}

	p.Set(`some`, `thing`)
	p.Set(`another`, `one`)
	out, _ := os.ReadFile(`testdata/temp.yaml`)
	fmt.Println(string(out))

	// Output:
	// another: one
	// some: thing
}
