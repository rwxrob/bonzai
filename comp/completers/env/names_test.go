package env_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/comp/completers/env"
)

func ExampleCompNames_Complete() {
	setupTestEnvironment(map[string]string{
		"ONE":   "1",
		"one":   "1",
		"TWO":   "2",
		"THREE": "3",
	})
	defer cleanupTestEnvironment()

	fmt.Println(env.CompNames.Complete())
	fmt.Println(env.CompNames.Complete(``))
	fmt.Println(env.CompNames.Complete(`o`))
	fmt.Println(env.CompNames.Complete(`O`))
	fmt.Println(env.CompNames.Complete(`T`))
	fmt.Println(env.CompNames.Complete(`TW`))
	// Output:
	// []
	// [ONE THREE TWO one]
	// [one]
	// [ONE]
	// [THREE TWO]
	// [TWO]
}

func ExampleCompNames_CompleteCaseInsensitive() {
	setupTestEnvironment(map[string]string{
		"one": "1",
		"ONE": "1",
		"TWO": "2",
	})
	defer cleanupTestEnvironment()

	c := env.NewCompNames("", true)
	fmt.Println(c.Complete())
	fmt.Println(c.Complete(``))
	fmt.Println(c.Complete(`O`))
	fmt.Println(c.Complete(`t`))
	// Output:
	// []
	// [ONE TWO one]
	// [ONE one]
	// [TWO]
}

func ExampleCompNames_CompletePrefix() {
	setupTestEnvironment(map[string]string{
		"A":     "a",
		"API_A": "a",
		"API_a": "a",
		"API_B": "b",
		"api_c": "c",
		"OTHER": "o",
	})
	defer cleanupTestEnvironment()

	c := env.NewCompNames("API_", true)
	fmt.Println(c.Complete())
	fmt.Println(c.Complete(``))
	fmt.Println(c.Complete(`A`))
	fmt.Println(c.Complete(`API_`))
	fmt.Println(c.Complete(`API_A`))
	fmt.Println(c.Complete(`c`))
	// Output:
	// []
	// [API_A API_B API_a api_c]
	// [API_A API_B API_a api_c]
	// [API_A API_B API_a api_c]
	// [API_A API_a]
	// [api_c]
}

func setupTestEnvironment(env map[string]string) {
	os.Clearenv()
	for key, value := range env {
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
}

func cleanupTestEnvironment() {
	os.Clearenv()
}
