package tr_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/fn/tr"
)

func ExamplePrefix() {

	FooPrefixer := tr.Prefix{`foo`}

	this := []string{`one`, `two`, `three`}

	fmt.Println(FooPrefixer.Transform(this))

	printType := func(an any) {
		switch v := an.(type) {
		case tr.Strings:
			fmt.Printf("I'm a string transformer: %v", v)
		default:
			fmt.Printf("unknown type: %T (%v)", v, v)
		}
	}

	printType(FooPrefixer)

	// Output:
	// [fooone footwo foothree]
	// I'm a string transformer: {foo}

}
