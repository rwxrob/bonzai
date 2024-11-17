// Copyright 2022 Robert S. Muhlestein
// SPDX-License-Identifier: Apache-2.0

package to_test

import (
	"fmt"
	"strings"

	"github.com/rwxrob/bonzai/fn"
	"github.com/rwxrob/bonzai/fn/each"

	"github.com/rwxrob/bonzai/to"
)

type stringer struct{}

func (s stringer) String() string { return "stringer" }

func ExampleString() {
	r := strings.NewReader(`reader`)
	stuff := []any{
		"some",
		[]byte{'t', 'h', 'i', 'n', 'g'},
		1, 2.234,
		stringer{},
		r,
	}
	for _, s := range stuff {
		fmt.Printf("%q ", to.String(s))
	}
	// Output:
	// "some" "thing" "1" "2.234" "stringer" "reader"
}

func Foo() {}

func ExampleFuncName() {
	f1 := func() {}
	f2 := func() {}

	// Foo is defined outside of the ExampleFuncName

	each.Println(fn.Map([]any{f1, f2, Foo, to.Lines}, to.FuncName))

	// Output:
	// func1
	// func2
	// Foo
	// Lines
}

func ExampleLines() {
	buf := `
some

thing 
here

mkay
`
	each.Print(to.Lines(buf))
	// Output:
	// something heremkay
}

type FooStruct struct{}

func (f FooStruct) String() string { return "FOO" }

type HumanFoo struct{}

func (f HumanFoo) Human() string { return "a friendly foo" }

func FooFunc(a any) {}

func ExampleHuman() {
	fmt.Println(to.Human('r'))                    // not number
	fmt.Println(to.Human("string💢good"))          // unescaped
	fmt.Println(to.Human(new(FooStruct)))         // has String()
	fmt.Println(to.Human(new(HumanFoo)))          // has Human()
	fmt.Println(to.Human([]rune{'r', 's', 'm'}))  // commas
	fmt.Println(to.Human([]string{"foo", "bar"})) // also commas
	fmt.Println(to.Human(func() {}))              // func1
	fmt.Println(to.Human(FooFunc))                // FooFunc
	// Output:
	// 'r'
	// "string💢good"
	// FOO
	// a friendly foo
	// ['r','s','m']
	// ["foo","bar"]
	// func1
	// FooFunc
}

func ExampleDedent_simple() {
	fmt.Printf("%q\n", to.Dedented("\n    foo\n    bar"))
	// Output:
	// "foo\nbar"
}

func ExampleDedent_tabs_or_Spaces() {
	fmt.Printf("%q\n", to.Dedented("\n\t\tfoo\n\t\tbar"))
	// Output:
	// "foo\nbar"
}

func ExampleDedent_multiple_Blank_Lines() {
	fmt.Printf("%q\n", to.Dedented("\n\n   \n\n    foo\n    bar"))
	fmt.Printf("%q\n", to.Dedented("\n   \n\n  \n   some"))
	// Output:
	// "foo\nbar"
	// "some"
}

func ExampleDedent_accidental_Chop() {
	fmt.Printf("%q\n", to.Dedented("\n\n   \n\n    foo\n   bar"))
	// Output:
	// "foo\nar"
}

func ExampleDedent_single_blank_line() {
	fmt.Printf("%q\n", to.Dedented("    \n"))
	// Output:
	// ""
}

func ExampleIndentation() {
	fmt.Println(to.Indentation("    some"))
	fmt.Println(to.Indentation("  some"))
	fmt.Println(to.Indentation("some"))
	fmt.Println(to.Indentation(" 💚ome"))
	// Output:
	// 4
	// 2
	// 0
	// 1
}

// wrapped, count = to.Wrapped("There I was not knowing what to do about this exceedingly long line and knowing that certain people would shun me for injecting\nreturns wherever I wanted.", 40)
func ExampleWrapped() {
	wrapped, count := to.Wrapped("some thing", 3)
	fmt.Printf("%q %v\n", wrapped, count)

	// Output:
	// "some\nthing" 2
}

func ExampleWrapped_long() {
	in := `
		The config command allows configuration of the current command in
		YAML and JSON (since all JSON is valid YAML). All changes to
		configuration values are done via the <edit> command since
		configurations may be complex and deeply nested in some cases.
		Querying configuration data, however, can be easily accomplished
		with the <query> command that uses jq-like selection syntax.`

	out, count := to.Wrapped(in, 80)
	fmt.Println(count)
	fmt.Printf("%q", out)

	// Output:
	// 58
	// "The config command allows configuration of the current command in YAML and JSON\n(since all JSON is valid YAML). All changes to configuration values are done via\nthe <edit> command since configurations may be complex and deeply nested in some\ncases. Querying configuration data, however, can be easily accomplished with the\n<query> command that uses jq-like selection syntax."
}

func ExampleWrapped_escapes() {
	in := `Here is a ` + "\033[34m" + `blue` + "\033[0m" + ` thing `
	out, count := to.Wrapped(in, 16)
	fmt.Println(count)
	fmt.Printf("%q", out)
	// Output:
	// 5
	// "Here is a \x1b[34mblue\x1b[0m\nthing"
}

func ExampleIndented() {
	fmt.Println("Indented:\n" + to.Indented("some\nthing", 4))
	// Output:
	// Indented:
	//     some
	//     thing
}

func ExamplePrefixed() {
	fmt.Println(to.Prefixed("some\nthing", "P  "))
	// Output:
	// P  some
	// P  thing
}

func ExampleIndentWrapped() {
	description := `
		The y2j command converts YAML (including references and
		anchors) to compressed JSON (with a single training newline) using
		the popular Go yaml.v3 package and its special <yaml:",inline"> tag.
		Because of this only YAML maps are supported as a base type (no
		arrays). An array can easily be done as the value of a map key.

		`

	fmt.Println("Indented:\n" + to.IndentWrapped(description, 7, 60))

	// Output:
	// Indented:
	//        The y2j command converts YAML (including references
	//        and anchors) to compressed JSON (with a single
	//        training newline) using the popular Go yaml.v3
	//        package and its special <yaml:",inline"> tag. Because
	//        of this only YAML maps are supported as a base type
	//        (no arrays). An array can easily be done as the value
	//        of a map key.
}

func ExampleMergedMaps() {
	m1 := map[string]any{"foo": 1, "bar": 2}
	m2 := map[string]any{"FOO": 1, "BAR": 2}
	m3 := map[string]any{"foo": 10}
	fmt.Println(to.MergedMaps(m1, m2, m3))
	// Output:
	// map[BAR:2 FOO:1 bar:2 foo:10]
}

func ExampleEscReturns() {
	fmt.Println(to.EscReturns("some\rthing\n"))
	// Output:
	// some\rthing\n
}

func ExampleUnEscReturns() {
	fmt.Printf("%q\n", to.UnEscReturns(`some\rthing\n`))
	// Output:
	// "some\rthing\n"
}

func ExampleHTTPS() {
	fmt.Println(to.HTTPS(`https://rwx.gg`))
	fmt.Println(to.HTTPS(`rwx.gg`))
	// Output:
	// https://rwx.gg
	// https://rwx.gg
}

func ExampleRuneCount() {
	// scan.Trace++
	fmt.Println(to.RuneCount("\033some\033"))
	fmt.Println(to.RuneCount("some"))
	fmt.Println(to.RuneCount("\033[32msome\033[0m"))
	// Output:
	// 4
	// 4
	// 4
}

func ExampleBytes_bytes() {
	fmt.Printf("%T", to.Bytes([]byte(`foo`)))
	// Output:
	// []uint8
}

func ExampleBytes_string() {
	fmt.Printf("%T", to.Bytes(`foo`))
	// Output:
	// []uint8
}

func ExampleBytes_runes() {
	fmt.Printf("%T", to.Bytes([]rune(`foo`)))
	// Output:
	// []uint8
}

func ExampleBytes_reader() {
	fmt.Printf("%T", to.Bytes(strings.NewReader(`foo`)))
	// Output:
	// []uint8
}

func ExampleBytes_bork() {
	it := to.Bytes(42)
	if it == nil {
		fmt.Println("yes, it is nil")
	}
	fmt.Printf("%v %T", it, it)
	// Output:
	// yes, it is nil
	// [] []uint8
}

func ExampleCrunchSpace() {
	fmt.Printf("%q\n", to.CrunchSpace(`here    is some`))
	fmt.Printf("%q\n", to.CrunchSpace(`   here is some`))
	fmt.Printf("%q\n", to.CrunchSpace(`here is some   `))
	fmt.Printf("%q\n", to.CrunchSpace("here is\nsome"))
	fmt.Printf("%q\n", to.CrunchSpace("here is\rsome"))
	fmt.Printf("%q\n", to.CrunchSpace("here is\r\n some"))

	// Output:
	// "here is some"
	// " here is some"
	// "here is some "
	// "here is some"
	// "here is some"
	// "here is some"
}

func ExampleTrimCrunchSpace() {
	fmt.Printf("%q\n", to.TrimCrunchSpace(`here    is some`))
	fmt.Printf("%q\n", to.TrimCrunchSpace(`   here is some`))
	fmt.Printf("%q\n", to.TrimCrunchSpace(`here is some   `))
	fmt.Printf("%q\n", to.TrimCrunchSpace("here is\nsome"))
	fmt.Printf("%q\n", to.TrimCrunchSpace("here is\rsome"))
	fmt.Printf("%q\n", to.TrimCrunchSpace("here is\r\n some"))

	// Output:
	// "here is some"
	// "here is some"
	// "here is some"
	// "here is some"
	// "here is some"
	// "here is some"
}

func ExampleTrimCrunchSpaceVisible() {
	fmt.Printf("%q\n", to.TrimCrunchSpaceVisible(" here  \033  is\nsome "))
	fmt.Printf("%q\n", to.TrimCrunchSpaceVisible(" here  \033  is \nsome "))

	// Output:
	// "here is some"
	// "here is some"
}

func ExampleType() {
	fmt.Println("true, false", to.Type("true", false))
	fmt.Println("false, false", to.Type("false", false))
	fmt.Println("false, true", to.Type("false", true))
	fmt.Println(`"", true`, to.Type("", true))
	fmt.Println(`"", false`, to.Type("", false))

	fmt.Println("1, 0", to.Type("1", 0))
	fmt.Println("14, 0", to.Type("14", 0))
	fmt.Println("-14, 0", to.Type("-14", 0))
	fmt.Println(`"", 0`, to.Type("", 0))

	fmt.Println("1, 0.0", to.Type("1", 0.0))
	fmt.Println("14, 0.0", to.Type("14", 0.0))
	fmt.Println("-14, 0.0", to.Type("-14", 0.0))
	fmt.Println(`"", 0.0`, to.Type("", 0.0))

	fmt.Println("hello, hello", to.Type("hello", "default"))

	// Output:
	// true, false true
	// false, false false
	// false, true false
	// "", true false
	// "", false false
	// 1, 0 1
	// 14, 0 14
	// -14, 0 -14
	// "", 0 0
	// 1, 0.0 1
	// 14, 0.0 14
	// -14, 0.0 -14
	// "", 0.0 0
	// hello, hello hello
}
