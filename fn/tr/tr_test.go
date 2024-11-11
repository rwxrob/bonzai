package tr_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/fn/tr"
)

func ExamplePrefix() {

	FooPrefixer := tr.Prefix{`foo`}

	fmt.Println(FooPrefixer.Transform(`thing`))

	printType := func(an any) {
		switch v := an.(type) {
		case tr.String:
			fmt.Printf("I'm a string transformer: %v", v)
		default:
			fmt.Printf("unknown type: %T (%v)", v, v)
		}
	}

	printType(FooPrefixer)

	// Output:
	// foothing
	// I'm a string transformer: {foo}

}
