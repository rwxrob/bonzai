package env_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/comp/completers/env"
)

func ExampleCompVars_Complete() {
	setupTestEnvironment(map[string]string{
		"ONE":   "1",
		"one":   "1",
		"TWO":   "2",
		"THREE": "3",
	})
	defer cleanupTestEnvironment()

	fmt.Println(env.CompVars.Complete())
	fmt.Println(env.CompVars.Complete(``))
	fmt.Println(env.CompVars.Complete(`o`))
	fmt.Println(env.CompVars.Complete(`O`))
	fmt.Println(env.CompVars.Complete(`T`))
	fmt.Println(env.CompVars.Complete(`TW`))
	// Output:
	// []
	// [$ONE $THREE $TWO $one]
	// [$one]
	// [$ONE]
	// [$THREE $TWO]
	// [$TWO]
}

func ExampleCompVars_CompleteCaseInsensitive() {
	setupTestEnvironment(map[string]string{
		"one": "1",
		"ONE": "1",
		"TWO": "2",
	})
	defer cleanupTestEnvironment()

	c := env.NewCompVars("", true)
	fmt.Println(c.Complete())
	fmt.Println(c.Complete(``))
	fmt.Println(c.Complete(`O`))
	fmt.Println(c.Complete(`t`))
	// Output:
	// []
	// [$ONE $TWO $one]
	// [$ONE $one]
	// [$TWO]
}

func ExampleCompVars_CompletePrefix() {
	setupTestEnvironment(map[string]string{
		"A":     "a",
		"API_A": "a",
		"API_a": "a",
		"API_B": "b",
		"api_c": "c",
		"OTHER": "o",
	})
	defer cleanupTestEnvironment()

	c := env.NewCompVars("API_", true)
	fmt.Println(c.Complete())
	fmt.Println(c.Complete(``))
	fmt.Println(c.Complete(`A`))
	fmt.Println(c.Complete(`API_`))
	fmt.Println(c.Complete(`API_A`))
	fmt.Println(c.Complete(`c`))
	// Output:
	// []
	// [$API_A $API_B $API_a $api_c]
	// [$API_A $API_B $API_a $api_c]
	// [$API_A $API_B $API_a $api_c]
	// [$API_A $API_a]
	// [$api_c]
}
