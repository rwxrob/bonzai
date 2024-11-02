package yaml_test

import "github.com/rwxrob/bonzai/yaml"

func ExampleThis_string() {
	this := yaml.This{"foo"}
	this.Print()
	this.This = "!some"
	this.Print()
	// Output:
	// foo
	// '!some'
}

func ExampleThis_numbers() {
	this := yaml.This{23434}
	this.Print()
	this.This = -24.24234
	this.Print()
	// Output:
	// 23434
	// -24.24234
}

func ExampleThis_bool() {
	this := yaml.This{true}
	this.Print()
	// Output:
	// true
}

func ExampleThis_nil() {
	this := yaml.This{nil}
	this.Print()
	// Output:
	// null
}

func ExampleThis_slice() {
	this := yaml.This{[]string{"foo", "bar"}}
	this.Print()
	// Output:
	// - foo
	// - bar
}

func ExampleThis_map() {
	this := yaml.This{map[string]string{"foo": "bar"}}
	this.Print()
	// Output:
	// foo: bar
}

func ExampleThis_struct() {
	this := yaml.This{struct {
		Foo   string
		Slice []string
	}{"foo", []string{"one", "two"}}}
	this.Print()
	// Output:
	// foo: foo
	// slice:
	//     - one
	//     - two
}
