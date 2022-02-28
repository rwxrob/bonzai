package scanner_test

import "github.com/rwxrob/bonzai/scanner"

func ExampleCur() {
	m := new(scanner.Cur)
	m.Print()
	//Output:
	// U+0000 '\x00' 0,0-0 (0-1)
}
