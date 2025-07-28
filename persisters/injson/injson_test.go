package injson_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/persisters/injson"
)

func ExampleSetup() {

	defer func() {
		os.Remove(`testdata/temp.json`)
	}()

	p := injson.Persister{`testdata/temp.json`}

	exists := func(path string) bool {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return false
		}
		return err == nil
	}

	fmt.Printf("before: %v\n", exists(`testdata/temp.json`))
	if err := p.Setup(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("after: %v\n", exists(`testdata/temp.json`))

	// Output:
	// before: false
	// after: true
}

func ExampleSetup_useExisting() {

	p := injson.Persister{`testdata/some.json`}

	exists := func(path string) bool {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return false
		}
		return err == nil
	}

	fmt.Printf("before: %v\n", exists(`testdata/some.json`))
	beforebuf, _ := os.ReadFile(`testdata/some.json`)
	if err := p.Setup(); err != nil {
		fmt.Println(err)
	}
	afterbuf, _ := os.ReadFile(`testdata/some.json`)
	fmt.Printf("after: %v\n", exists(`testdata/some.json`))

	if string(beforebuf) == string(afterbuf) {
		fmt.Println(`same text before and after`)
	}

	// Output:
	// before: true
	// after: true
	// same text before and after
}

func ExampleGet() {
	p := injson.Persister{`testdata/some.json`}
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
		os.Remove(`testdata/temp.json`)
	}()

	p := injson.Persister{`testdata/temp.json`}

	p.Set(`some`, `thing`)
	p.Set(`another`, `one`)
	out, _ := os.ReadFile(`testdata/temp.json`)
	fmt.Println(string(out))

	// Output:
	// {"another":"one","some":"thing"}
}
