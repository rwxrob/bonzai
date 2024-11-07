package is_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/is"
)

func ExampleAllLatinASCIILower() {
	tests := []string{
		"hello",
		"world",
		"Hello",
		"123",
		"abc123",
		"abcdEF",
		"",
		"abc-def",
		"latinlower",
		"ábc",
	}

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
	tests := []string{
		"hello",
		"world",
		"Hello",
		"123",
		"abc123",
		"abcdEF",
		"",
		"abc-def",
		"latinlower",
		"ábc",
		"-fail",
		"its-all-fine",
		"even-this-",
	}

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
	tests := []string{
		"HELLO",
		"WORLD",
		"Hello",
		"123",
		"abc123",
		"abcdEF",
		"",
		"ABC-DEF",
		"LATINLOWER",
		"ÁBC",
	}

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

func ExampleTruthy() {
	fmt.Println("true", is.Truthy("true"))
	fmt.Println("t", is.Truthy("t"))
	fmt.Println("on", is.Truthy("on"))
	fmt.Println("1", is.Truthy("1"))
	fmt.Println("5", is.Truthy("5"))
	fmt.Println("500", is.Truthy("500"))

	fmt.Println("false", is.Truthy("false"))
	fmt.Println("f", is.Truthy("f"))
	fmt.Println("off", is.Truthy("off"))
	fmt.Println("0", is.Truthy("0"))
	fmt.Println("-1", is.Truthy("-1"))
	fmt.Println("-5", is.Truthy("-5"))
	fmt.Println("-500", is.Truthy("-500"))

	fmt.Println("", is.Truthy(""))
	fmt.Println(" ", is.Truthy(" "))
	fmt.Println("\t", is.Truthy("\t"))
	fmt.Println("\n", is.Truthy("\n"))
	fmt.Println("foo", is.Truthy("foo"))
	fmt.Println("f4g5g5", is.Truthy("f4g5g5"))
	fmt.Println("~:", is.Truthy("~:"))

	// Output:
	// true true
	// t true
	// on true
	// 1 true
	// 5 true
	// 500 true
	// false false
	// f false
	// off false
	// 0 false
	// -1 false
	// -5 false
	// -500 false
	//  false
	//   false
	// 	 false
	//
	//  false
	// foo false
	// f4g5g5 false
	// ~: false
}
