package trace_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/trace"
)

func ExampleTrace() {
	fmt.Println(trace.Flags)
	fmt.Println(trace.AutoUpdate)
	trace.Flags |= trace.AutoUpdate
	fmt.Println(trace.Flags == trace.Flags|trace.AutoUpdate)
	fmt.Println(trace.Flags == trace.Flags|trace.Emerg)
	// Output:
	// 0
	// 512
	// true
	// false
}
