package is_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/pkg/core/is"
)

func ExampleAllLatinASCIILower() {

	tests := []string{"hello", "world", "Hello", "123", "abc123", "abcdEF", "", "abc-def", "latinlower", "ábc"}

	for _, test := range tests {
		result := is.AllLatinASCIILower(test)
		fmt.Printf("%q %v\n", test, result)
	}

	// Output:
	// "hello" true
	// "world" true
	// "Hello" false
	// "123" false
	// "abc123" false
	// "abcdEF" false
	// "" true
	// "abc-def" false
	// "latinlower" true
	// "ábc" false
}

func ExampleAllLatinASCIIUpper() {

	tests := []string{"HELLO", "WORLD", "Hello", "123", "abc123", "abcdEF", "", "ABC-DEF", "LATINLOWER", "ÁBC"}

	for _, test := range tests {
		result := is.AllLatinASCIIUpper(test)
		fmt.Printf("%q %v\n", test, result)
	}

	// Output:
	// "HELLO" true
	// "WORLD" true
	// "Hello" false
	// "123" false
	// "abc123" false
	// "abcdEF" false
	// "" true
	// "ABC-DEF" false
	// "LATINLOWER" true
	// "ÁBC" false

}
