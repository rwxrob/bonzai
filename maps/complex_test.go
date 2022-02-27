package maps_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/maps"
)

func ExamplePrefix() {
	fmt.Println(maps.Prefix([]string{"foo", "bar"}, "my"))
	// Output:
	// [myfoo mybar]
}

func ExampleKeys() {
	m1 := map[string]int{"two": 2, "three": 3, "one": 1}
	m2 := map[string]string{"two": "two", "three": "three", "one": "one"}
	fmt.Println(maps.Keys(m1))
	fmt.Println(maps.Keys(m2))
	// Output:
	// [one three two]
	// [one three two]
}
