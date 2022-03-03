package scan_test

import (
	"fmt"
	"log"
	"os"
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
	//Output:
	// U+0073 's' 1,1-1 (1-1)
	// <EOD>
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
	s.Print()
	s.Scan()
	s.Print()
	// Output:
	// U+006E 'n' 1,9-9 (9-9)
	// U+0067 'g' 1,10-10 (10-10)
	// <EOD>
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

func ExampleExpect_lk() {
	s, _ := scan.New("some thing")
	_, e := s.Expect(is.Lk{"foo"})
	fmt.Println(e)
	c, _ := s.Expect(is.Lk{"foo", 's'})
	c.Print()
	s.ScanN(2)
	s.Print()
	c, _ = s.Expect(is.Lk{is.Rng{'l', 'p'}})
	s.Print() // not advanced
	c, _ = s.Expect(is.In{is.Rng{'l', 'p'}})
	s.Print() // advanced
	// Output:
	// expected ["foo"] at U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
	// U+006D 'm' 1,3-3 (3-3)
	// U+006D 'm' 1,3-3 (3-3)
	// U+0065 'e' 1,4-4 (4-4)
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
	// unexpected "some" at U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_in() {
	s, _ := scan.New("some thing")
	s.Scan()
	c, _ := s.Expect(is.In{'O', 'o', "ome"})
	c.Print()
	s.Print()
	_, err := s.Expect(is.In{'x', 'y', "zee"})
	fmt.Println(err)
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// U+006D 'm' 1,3-3 (3-3)
	// expected one of ['x' 'y' "zee"] at U+006D 'm' 1,3-3 (3-3)
}

func ExampleExpect_seq() {
	s, _ := scan.New("some thing")
	s.Snap()
	s.Expect(is.Seq{"some", ' ', "thing"})
	s.Print()
	s.Back()
	s.Print()
	_, err := s.Expect(is.Seq{"some", "thing"})
	fmt.Println(err)
	// Output:
	// <EOD>
	// U+0073 's' 1,1-1 (1-1)
	// expected rune 't' at U+0020 ' ' 1,5-5 (5-5)
}

func ExampleExpect_opt() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(is.Opt{"thing", "some"})
	c.Print()
	s.Print()
	c, _ = s.Expect(is.Opt{"foo"})
	s.Print() // no change
	// Output:
	// U+0065 'e' 1,4-4 (4-4)
	// U+0020 ' ' 1,5-5 (5-5)
	// U+0020 ' ' 1,5-5 (5-5)
}

func ExampleExpect_mn1() {
	s, _ := scan.New("sommme thing")
	start := s.Mark()
	s.ScanN(2)
	c, _ := s.Expect(is.Mn1{'m'}) // goggles up all three
	c.Print()
	s.Print()
	s.Jump(start)
	c, _ = s.Expect(is.Mn1{'s'}) // yep, just one
	c.Print()
	s.Print()
	// Output:
	// U+006D 'm' 1,5-5 (5-5)
	// U+0065 'e' 1,6-6 (6-6)
	// U+0073 's' 1,1-1 (1-1)
	// U+006F 'o' 1,2-2 (2-2)
}

func ExampleExpect_min() {
	s, _ := scan.New("sommme thing")
	start := s.Mark()
	s.ScanN(2)
	c, _ := s.Expect(is.Min{2, 'm'}) // goggles up all three
	c.Print()
	s.Print()
	s.Jump(start)
	_, e := s.Expect(is.Min{2, 's'}) // nope, only one found
	fmt.Println(e)
	s.Print()
	// Output:
	// U+006D 'm' 1,5-5 (5-5)
	// U+0065 'e' 1,6-6 (6-6)
	// expected min 2 of 's' at U+006F 'o' 1,2-2 (2-2)
	// U+006F 'o' 1,2-2 (2-2)
}

func ExampleExpect_min_Max() {
	s, _ := scan.New("sommme thing")
	s.Snap()
	s.ScanN(2)
	s.Print()
	s.Expect(is.MMx{1, 3, 'm'}) // goggles up all three
	s.Print()
	s.Back()
	s.Expect(is.MMx{1, 3, 's'}) // yep, at least one
	s.Print()
	_, err := s.Expect(is.MMx{1, 3, 'X'}) // nope
	fmt.Println(err)
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// U+0065 'e' 1,6-6 (6-6)
	// U+006F 'o' 1,2-2 (2-2)
	// expected min 1, max 3 of 'X' at U+006F 'o' 1,2-2 (2-2)
}

func ExampleExpect_count() {
	s, _ := scan.New("sommme thing")
	s.Snap()
	s.ScanN(2)
	s.Print()
	s.Expect(is.N{3, 'm'}) // goggles up all three
	s.Print()
	s.Back()
	s.Expect(is.N{1, 's'}) // yes, but silly since 's' is easier
	s.Print()
	_, err := s.Expect(is.N{3, 'X'}) // nope
	fmt.Println(err)
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// U+0065 'e' 1,6-6 (6-6)
	// U+006F 'o' 1,2-2 (2-2)
	// expected exactly 3 of 'X' at U+006F 'o' 1,2-2 (2-2)
}

func ExampleExpect_in_Range() {
	s, _ := scan.New("some thing")
	s.Scan()
	c1, _ := s.Expect(is.Rng{'l', 'p'})
	c1.Print()
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// U+006D 'm' 1,3-3 (3-3)
}

func ExampleExtendExpect() {
	s, _ := scan.New("some thing")
	s.ExtendExpect = func(s *scan.Scanner, a ...any) (*scan.Cur, error) {
		return s.Cur, fmt.Errorf("custom error for type %T handled at %v",
			a[0], s.Cur,
		)
	}
	_, e := s.Expect([]byte{'0'})
	fmt.Println(e)
	// Output:
	// custom error for type []uint8 handled at U+0073 's' 1,1-1 (1-1)
}

func ExampleSnap() {
	s, _ := scan.New("some thing")
	s.ScanN(3)
	s.Snap()
	s.Print()
	s.ScanN(4)
	s.Print()
	s.Back()
	s.Print()
	// Output:
	// U+0065 'e' 1,4-4 (4-4)
	// U+0069 'i' 1,8-8 (8-8)
	// U+0065 'e' 1,4-4 (4-4)
}

func ExampleScan() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	s, _ := scan.New(`s😈me thing`)
	s.Scan()
	s.Print()
	s.Scan()
	s.Print()
	s.Log()
	// Output:
	// U+1F608 '😈' 1,2-2 (2-2)
	// U+006D 'm' 1,3-6 (3-6)
	// U+006D 'm' 1,3-6 (3-6)
}
