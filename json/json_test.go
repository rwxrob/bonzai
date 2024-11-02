package json_test

import (
	stdjson "encoding/json"
	"fmt"

	"github.com/rwxrob/bonzai/fn"
	json "github.com/rwxrob/bonzai/json"
)

func ExampleEscape() {
	set := fn.A[string]{
		`<`, `>`, `&`, `"`, `'`,
		"\t", "\b", "\f", "\n", "\r",
		"\\", "\"", "💢", "д",
	}
	set.Map(json.Escape).Print()
	// Output:
	// <>&\"'\t\b\f\n\r\\\"💢д
}

func ExampleMarshal() {
	m := map[string]string{"<foo>": "&bar"}

	// the good way
	buf, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))

	// the broken encoding/json way
	buf, err = stdjson.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))

	// Output:
	// {"<foo>":"&bar"}
	// {"\u003cfoo\u003e":"\u0026bar"}

}

func ExampleMarshalIndent() {
	m := map[string]string{"<foo>": "&bar"}

	// the good way
	buf, err := json.MarshalIndent(m, " ", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))

	// the broken encoding/json way
	buf, err = stdjson.MarshalIndent(m, " ", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))

	// Output:
	// {
	//   "<foo>": "&bar"
	//  }
	// {
	//   "\u003cfoo\u003e": "\u0026bar"
	//  }
}

func ExampleUnmarshal() {
	m := new(map[string]string)
	if err := json.Unmarshal([]byte(`{"<foo>":"&bar"}`), m); err != nil {
		fmt.Println(err)
	}
	fmt.Println(m)
	if err := json.Unmarshal([]byte(`{"<foo>":"&bar"}`), m); err != nil {
		fmt.Println(err)
	}
	fmt.Println(m)
	// Output:
	// &map[<foo>:&bar]
	// &map[<foo>:&bar]
}

func ExampleThis_string() {
	this := json.This{"foo"}
	this.Print()
	this.This = "!some"
	this.Print()
	// Output:
	// "foo"
	// "!some"
}

func ExampleThis_numbers() {
	this := json.This{23434}
	this.Print()
	this.This = -24.24234
	this.Print()
	// Output:
	// 23434
	// -24.24234
}

func ExampleThis_bool() {
	this := json.This{true}
	this.Print()
	// Output:
	// true
}

func ExampleThis_nil() {
	this := json.This{nil}
	this.Print()
	// Output:
	// null
}

func ExampleThis_slice() {
	this := json.This{[]string{"foo", "bar"}}
	this.Print()
	// Output:
	// ["foo","bar"]
}

func ExampleThis_map() {
	this := json.This{map[string]string{"foo": "bar"}}
	this.Print()
	// Output:
	// {"foo":"bar"}
}

func ExampleThis_struct() {
	this := json.This{struct {
		Foo   string
		Slice []string
	}{"foo", []string{"one", "two"}}}
	this.Print()
	// Output:
	// {"Foo":"foo","Slice":["one","two"]}
}
