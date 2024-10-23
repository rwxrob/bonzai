package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/pkg/core/comp"
)

func ExampleThreeLetterEngWeekday() {
	fmt.Println(comp.ThreeLetterEngWeekday.Complete(nil, `t`))
	// Output:
	// [tue thu]
}
