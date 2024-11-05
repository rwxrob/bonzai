package vars_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/vars"
)

func ExampleSet() {

	file := `testdata/settest.properties`

	defer func() {
		err := os.Remove(file)
		fmt.Println(err)
	}()

	if err := vars.Set(`somekey`, `someval`, file); err != nil {
		fmt.Println(err)
	}

	if err := vars.Set(`otherkey`, ``, file); err != nil {
		fmt.Println(err)
	}

	futil.Cat(file)

	// Unordered output:
	// somekey=someval
	// otherkey=
	// <nil>
}

func ExampleGet() {

	file := `testdata/vars.properties`

	value, err := vars.Get(`.pomo.warn`, file)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(value)

	// Output:
	// 1m
}
