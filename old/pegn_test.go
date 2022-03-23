package scan_test

import "github.com/rwxrob/scan"

func ExampleR_Any() {
	s, _ := scan.New("so")
	s.Print()
	s.Any(2)
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// <EOD>
}

func ExampleR_Rune() {
	s, _ := scan.New("some thing")
	//s.Any(3)
	s.Rune('s')
	s.Rune('o')
	s.Rune('m')
	s.Print() // same as Any(3)
	// Output:
	// U+0065 'e' 1,4-4 (4-4)
}

func ExampleR_Rune_fail() {
	s, _ := scan.New("some thing")
	defer s.PrintPanic()
	s.Rune('s')
	s.Rune('o')
	s.Rune('m')
	s.Rune('e')
	s.Rune('\t')
	// Output:
	// expected '\t' at U+0020 ' ' 1,5-5 (5-5)
}

func ExampleR_Str() {
	s, _ := scan.New("some thing")
	s.Str("som")
	s.Print()
	// Output:
	// U+0065 'e' 1,4-4 (4-4)
}

func ExampleR_Str_fail() {
	s, _ := scan.New("some thing")
	defer s.PrintPanic()
	s.Str("some\t")
	// Output:
	// expected '\t' at U+0020 ' ' 1,5-5 (5-5)
}
