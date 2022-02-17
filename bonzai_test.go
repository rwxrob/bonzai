package bonzai_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
)

func ExampleArgsFrom() {
	fmt.Printf("%q\n", bonzai.ArgsFrom(`greet  hi french`))
	fmt.Printf("%q\n", bonzai.ArgsFrom(`greet hi   french `))
	// Output:
	// ["greet" "hi" "french"]
	// ["greet" "hi" "french" " "]
}
