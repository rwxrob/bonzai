package vars_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/pkg/futil"
	"github.com/rwxrob/bonzai/pkg/vars"
)

func ExampleVars() {
	m := vars.New()
	m.Id = `foo`
	m.Dir = `testdata`
	m.File = `vars`
	fmt.Println(m.Path())
	fmt.Println(m.DirPath())
	// Output:
	// testdata/foo/vars
	// testdata/foo
}

func ExampleVars_Init() {

	m := vars.New()
	m.Id = `foo`
	m.Dir = `testdata`
	m.File = `vars`

	defer func() { os.RemoveAll(m.DirPath()) }()

	m.Init()
	fmt.Println(futil.Exists(`testdata/foo/vars`))

	// Output:
	// true
}

func ExampleVars_Set() {

	m := vars.New()
	m.Id = `foo`
	m.Dir = `testdata`
	m.File = `vars`

	defer func() { os.RemoveAll(m.DirPath()) }()

	m.Init()
	if err := m.Set("some", "thing\nhere"); err != nil {
		fmt.Println(err)
	}
	byt, _ := os.ReadFile(`testdata/foo/vars`)
	fmt.Println(string(byt) == `some=thing\nhere`+"\n")

	// Output:
	// true
}

func ExampleVars_Get() {

	m := vars.New()
	m.Id = `foo`
	m.Dir = `testdata`
	m.File = `vars`

	defer func() { os.RemoveAll(m.DirPath()) }()

	m.Init()
	if err := m.Set("some", "thing\nhere"); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%q\n", m.Get(`some`))

	// Output:
	// "thing\nhere"
}

func ExampleVars_UnmarshalText() {

	in := `
some=thing here
another=one over here
`

	m := vars.New()
	m.UnmarshalText([]byte(in))
	fmt.Println(len(m.M))
	fmt.Println(m.M["some"])
	fmt.Println(m.M["another"])

	// Output:
	// 2
	// thing here
	// one over here
}

func ExampleVars_MarshalText() {

	m := vars.New()
	m.M["some"] = "thing here"
	m.M["another"] = "one\rhere\nbut all good"

	byt, err := m.MarshalText()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(byt))

	// Ordered Output:
	// some:thing here
	// another:one\rhere\nbut all good
}
