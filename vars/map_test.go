package vars_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/vars"
)

func ExampleNewMapFrom() {

	m, err := vars.NewMapFrom(`testdata/vars.properties`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m.M[`.pomo.warn`])

	// Output:
	// 1m
}

func ExampleMap_Get() {

	m, err := vars.NewMapFrom(`testdata/vars.properties`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m.Get(`.pomo.warn`))

	// Output:
	// 1m <nil>

}

func ExampleMap_Get_not_found() {

	m, err := vars.NewMapFrom(`testdata/vars.properties`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m.Get(`.pomo.wain`))

	// Output:
	// could not find key: .pomo.wain

}

func ExampleVars_Init() {

	m := vars.NewMap()
	m.File = `testdata/other.properties`

	defer func() {
		if err := os.RemoveAll(m.File); err != nil {
			fmt.Println(err)
		}
	}()

	if err := m.Init(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(futil.Exists(`testdata/other.properties`))

	// Output:
	// true
}

func ExampleVars_Set() {

	m := vars.NewMap()
	m.File = `testdata/settest.properties`
	m.Init()

	// cleanup after the test completes
	defer func() {
		if err := os.RemoveAll(m.File); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println(m.Set("some", "thing\nhere"))
	byt, _ := os.ReadFile(`testdata/settest.properties`)
	fmt.Println(string(byt) == `some=thing\nhere`+"\n")

	// Output:
	// <nil>
	// true
}

func ExampleVars_MarshalText() {

	m := vars.NewMap()
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

func ExampleMap_Clear() {
	m := vars.NewMap()
	m.M = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := m.Clear()
	if err != nil {
		fmt.Println(err)
	}

	// The map should be empty after clearing
	fmt.Println(len(m.M))

	// Output:
	// 0
}

func ExampleMap_GrepK() {
	m := vars.NewMap()
	m.M = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	val, err := m.GrepK(`y2`)
	if err != nil {
		fmt.Println(err)
	}

	// The map should be empty after clearing
	fmt.Println(val)

	// Output:
	// key2=value2
}

func ExampleMap_GrepK_nokey() {
	m := vars.NewMap()
	m.M = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	val, err := m.GrepK(`foo`)
	if err != nil {
		fmt.Println(err)
	}

	// The map should be empty after clearing
	fmt.Printf("%q\n", val)

	// Output:
	// ""

}

func ExampleMap_Data() {

	m, err := vars.NewMapFrom(`testdata/vars.properties`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m.Data())

	// Unordered output:
	// .pomo.warn=1m
	// .pomo.prefix=🍅
	// .pomo.prefixwarn=💢
	// .pomo.duration=52m
	// .pomo.interval=20s
	//  <nil>

}

func ExampleMap_Print() {
	m, err := vars.NewMapFrom(`testdata/vars.properties`)
	if err != nil {
		fmt.Println(err)
	}
	m.Print()

	// Unordered output:
	// .pomo.warn=1m
	// .pomo.prefix=🍅
	// .pomo.prefixwarn=💢
	// .pomo.duration=52m
	// .pomo.interval=20s

}

func ExampleMap_Delete() {
	m, err := vars.NewMapFrom(`testdata/vars.properties`)
	if err != nil {
		fmt.Println(err)
	}
	defer m.Set(`.pomo.prefix`, `🍅`)
	m.Delete(`.pomo.prefix`)
	m.Print()

	// Unordered output:
	// .pomo.warn=1m
	// .pomo.prefixwarn=💢
	// .pomo.duration=52m
	// .pomo.interval=20s
}
