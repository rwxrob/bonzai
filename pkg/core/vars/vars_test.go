package vars_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/pkg/core/futil"
	"github.com/rwxrob/bonzai/pkg/core/vars"
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
	fmt.Println(m.Set("some", "thing\nhere"))
	byt, _ := os.ReadFile(`testdata/foo/vars`)
	fmt.Println(string(byt) == `some=thing\nhere`+"\n")

	// Output:
	// 1
	// true
}

func ExampleVars_Get() {

	m := vars.New()
	m.Id = `foo`
	m.Dir = `testdata`
	m.File = `vars`

	defer func() { os.RemoveAll(m.DirPath()) }()

	m.Init()
	fmt.Println(m.Set("some", "thing\nhere"))
	val, code := m.Get(`some`)
	fmt.Printf("%q %v\n", val, code)

	// Output:
	// 1
	// "thing\nhere" 1
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
