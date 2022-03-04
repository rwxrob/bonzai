package tk_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/scan/tk"
)

func Example() {
	fmt.Printf("%U\n", tk.EOD)
	fmt.Printf("%U\n", tk.ANY)
	// Output:
	// U+7FFFFFFF
	// U+7FFFFFFE
}
