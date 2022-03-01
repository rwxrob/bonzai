package scanner_test

import "github.com/rwxrob/bonzai/scanner"

func ExampleCur() {
	m := new(scanner.Cur)
	m.Print()
	m.NewLine()
	m.Print()
	//Output:
	// U+0000 '\x00' 0,0-0 (0-1)
	// U+0000 '\x00' 1,1-1 (0-1)
}
