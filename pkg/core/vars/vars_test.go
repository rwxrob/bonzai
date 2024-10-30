package vars_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/pkg/core/futil"
	"github.com/rwxrob/bonzai/pkg/core/vars"
)

func ExampleAdd() {

	file := `testdata/addtest.properties`

	defer func() {
		err := os.Remove(file)
		fmt.Println(err)
	}()

	if err := vars.Add(`somekey`, `someval`, file); err != nil {
		fmt.Println(err)
	}

	if err := vars.Add(`otherkey`, ``, file); err != nil {
		fmt.Println(err)
	}

	futil.Cat(file)

	// Output:
	// somekey=someval
	// otherkey=
	// <nil>
}
