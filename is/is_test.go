package is_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/is"
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
	// "" false
	// "abc-def" false
	// "latinlower" true
	// "ábc" false
}

func ExampleAllLatinASCIILowerWithDashes() {

	tests := []string{"hello", "world", "Hello", "123", "abc123", "abcdEF", "", "abc-def", "latinlower", "ábc", "-fail", "its-all-fine", "even-this-"}

	for _, test := range tests {
		result := is.AllLatinASCIILowerWithDashes(test)
		fmt.Printf("%q %v\n", test, result)
	}

	// Output:
	// "hello" true
	// "world" true
	// "Hello" false
	// "123" false
	// "abc123" false
	// "abcdEF" false
	// "" false
	// "abc-def" true
	// "latinlower" true
	// "ábc" false
	// "-fail" false
	// "its-all-fine" true
	// "even-this-" false
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
	// "" false
	// "ABC-DEF" false
	// "LATINLOWER" true
	// "ÁBC" false

}
