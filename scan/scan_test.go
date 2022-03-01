package scan_test

import (
	"fmt"
	"strings"

	"github.com/rwxrob/bonzai/scan"
	"github.com/rwxrob/bonzai/scan/is"
)

func ExampleNew_string() {
	s, err := scan.New("some thing")
	if err != nil {
		fmt.Println(err)
	}
	s.Print()
	s.ScanN(3)
	s.Print()
	//Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+0065 'e' 1,4-4 (4-4)
}

func ExampleNew_bytes() {
	s, err := scan.New([]byte{'s', 'o', 'm'})
	if err != nil {
		fmt.Println(err)
	}
	s.Print()
	s.ScanN(3)
	s.Print()
	fmt.Println(s.Done())
	//Output:
	// U+0073 's' 1,1-1 (1-1)
	// <EOD>
	// true
}

func ExampleNew_reader() {
	r := strings.NewReader("some thing")
	s, err := scan.New(r)
	if err != nil {
		fmt.Println(err)
	}
	s.Print()
	s.ScanN(3)
	s.Print()
	//Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+0065 'e' 1,4-4 (4-4)
}

func ExampleInit() {
	s, err := scan.New("some thing")
	if err != nil {
		fmt.Println(err)
	}
	s.Init("other stuff entirely")
	s.Print()
	s.ScanN(3)
	s.Print()
	s.Scan()
	s.Print()
	//Output:
	// U+006F 'o' 1,1-1 (1-1)
	// U+0065 'e' 1,4-4 (4-4)
	// U+0072 'r' 1,5-5 (5-5)
}

func ExampleMark() {
	s, err := scan.New("some thing")
	if err != nil {
		fmt.Println(err)
	}
	m := s.Mark()
	//log.Printf("%p", s.Cur)
	//log.Printf("%p", m)
	fmt.Println(s.Cur != m)
	// Output:
	// true
}

func ExampleJump() {
	s1, _ := scan.New("some thing")
	s1.ScanN(5)
	s1.Print() // t

	s2, _ := scan.New("other thing")
	s2.ScanN(6)
	s2.Print()      // t
	s1.Jump(s2.Cur) // WRONG, must be same source buffer
	s1.Print()

	s3, _ := scan.New("some thing") // identical
	s3.ScanN(6)
	s3.Print() // h
	s1.Jump(s3.Cur)
	s1.Print()
	s3.ScanN(2)
	s1.Jump(s3.Cur)
	s1.Print()
	s3.Print()

	// Output:
	// U+0074 't' 1,6-6 (6-6)
	// U+0074 't' 1,7-7 (7-7)
	// U+0074 't' 1,7-7 (7-7)
	// U+0068 'h' 1,7-7 (7-7)
	// U+0068 'h' 1,7-7 (7-7)
	// U+006E 'n' 1,9-9 (9-9)
	// U+006E 'n' 1,9-9 (9-9)

}

func ExamplePeek() {
	s, _ := scan.New("some thing")
	s.ScanN(6)
	fmt.Println(s.Peek(3))
	// Output:
	// hin
}

func ExampleLook() {
	s, _ := scan.New("some thing")
	s.Scan()
	m1 := s.Mark()
	m1.Print()
	s.ScanN(3)
	fmt.Printf("%q\n", s.Look(m1)) // look behind
	s.ScanN(4)
	m2 := s.Mark()
	m2.Print()
	s.Jump(m1)
	fmt.Printf("%q\n", s.Look(m2)) // look ahead
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// "ome "
	// U+006E 'n' 1,9-9 (9-9)
	// "ome thin"
}

func ExampleLookSlice() {
	s, _ := scan.New("some thing")
	s.Scan()
	m1 := s.Mark()
	m1.Print()
	s.ScanN(7)
	m2 := s.Mark()
	m2.Print()
	fmt.Println(s.LookSlice(m1, m2))
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// U+006E 'n' 1,9-9 (9-9)
	// ome thin
}

func ExampleNewLine() {
	s, _ := scan.New("some thing")
	s.Print()
	s.NewLine()
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 2,1-1 (1-1)
}

func ExampleErrorExpected() {
	s, _ := scan.New("some thing")
	fmt.Println(s.ErrorExpected("foo"))
	fmt.Println(s.ErrorExpected('f'))
	fmt.Println(s.ErrorExpected([]byte{'f', 'o', 'o'}))
	// Output:
	// expected string "foo" at U+0073 's' 1,1-1 (1-1)
	// expected rune 'f' at U+0073 's' 1,1-1 (1-1)
	// expected []uint8 "foo" at U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_basic() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect("some", ' ', "thin")
	c.Print()
	fmt.Println(s.Done())
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(s.Done())
	// Output:
	// U+006E 'n' 1,9-9 (9-9)
	// false
	// U+0067 'g' 1,10-10 (10-10)
	// <EOD>
	// true
}

func ExampleCheck() {
	s, _ := scan.New("some thing")
	c, _ := s.Check("some", ' ', "thin") // same as Expect ...
	c.Print()                            // ... with cur return ...
	s.Print()                            // ... just doesn't advance
	// Output:
	// U+006E 'n' 1,9-9 (9-9)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_not() {
	s, _ := scan.New("some thing")
	c1, e1 := s.Expect(is.Not{"foo"})
	c1.Print()
	fmt.Println(e1)
	c2, e2 := s.Expect(is.Not{"some"})
	c2.Print()
	fmt.Println(e2)
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// <nil>
	// <nil>
	// not expecting "some" at U+0073 's' 1,1-1 (1-1)
}
