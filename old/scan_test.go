package scan_test

import (
	"fmt"
	"log"
	"os"

	"github.com/rwxrob/bonzai/scan"
	z "github.com/rwxrob/bonzai/scan/is"
	"github.com/rwxrob/bonzai/scan/tk"
)

/*

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

func ExampleNewLine() {
	s, _ := scan.New("some thing")
	s.Print()
	s.NewLine()
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 2,1-1 (1-1)
}

/*
func ExampleR_Expect_parse_Single_Success() {
	const FOO = `FOO`
	s, _ := scan.New("some thing")
	//c, _ := s.Expect(z.P{"some", ' ', z.I{'t', 'T'}})//  FIXME
	c, _ := s.Expect(z.P{FOO, "some", ' ', 't'})
	c.Print() // same as "some t", points at 'h'
	s.Print()
	// Output:
	// U+0074 't' 1,6-6 (6-6)
	// U+0068 'h' 1,7-7 (7-7)

}
*/

/*

func ExampleR_Expect_parse_Success_One_Deep() {
	const P = `PHRASE`
	const S = `STRING`
	const R = `RUNE`
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.P{P, z.P{S, "some"}, z.P{R, ' '}, z.P{R, 't'}})
	c.Print() // same as "some t", points at 't'
	//s.Print() // advances to 'h
	for _, v := range s.Nodes {
		fmt.Println(v)
	}
	// Output:
	// U+0074 't' 1,6-6 (6-6)
	// U+0068 'h' 1,7-7 (7-7)
	// not null
}

func ExampleR_Expect_parse_Succes_Mixed() {
	const P = `PHRASE`
	const R = `RUNE`
	s, _ := scan.New("some thing")
	c, e := s.Expect(z.P{P, "some", z.P{R, ' '}, "thing"})
	fmt.Println(e)
	c.Print() // same as "some t", points at 'h'
	s.Print()
	//s.CurNode.Print()
	//fmt.Println(s.CurNode == s.Nodes[0])
	// Output:
	// U+0068 'h' 1,7-7 (7-7)
	// U+0068 'h' 1,7-7 (7-7)
	// {"T":"FOO","V":"some t"}
	// true
}
*/

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

func ExampleExpect_string() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect("som")
	c.Print() // same as s.ScanN(2), last is 'm'
	s.Print() // point to next scan 'e'
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// U+0065 'e' 1,4-4 (4-4)
}

func ExampleExpect_compound_Expr_String() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.X{"some ", "thin"})
	c.Print()
	s.Print()
	// Output:
	// U+006E 'n' 1,9-9 (9-9)
	// U+0067 'g' 1,10-10 (10-10)
}

func ExampleExpect_compound_Expr_Rune() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.X{"some", ' ', "thin"})
	c.Print()
	s.Print()
	// Output:
	// U+006E 'n' 1,9-9 (9-9)
	// U+0067 'g' 1,10-10 (10-10)
}

func ExampleExpect_it_Success() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.Y{"some"})
	c.Print() // even though true, not moved
	s.Print() // scanner also not moved
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_it_Success_Middle() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.X{"some", z.Y{' '}})
	c.Print() // advanced up to (but not including) ' '
	s.Print() // scanner also not moved
	// Output:
	// U+0020 ' ' 1,5-5 (5-5)
	// U+0020 ' ' 1,5-5 (5-5)
}

func ExampleExpect_it_Fail() {
	s, _ := scan.New("some thing")
	_, err := s.Expect(z.X{"some", z.Y{"thing"}})
	fmt.Println(err)
	s.Print() // but scanner did get "some" so advanced
	// Output:
	// expected "thing" at U+0020 ' ' 1,5-5 (5-5)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_it_Fail_X() {
	s, _ := scan.New("some thing")
	_, err := s.Expect(z.X{"some", z.Y{"thing"}})
	fmt.Println(err)
	s.Print() // but scanner did get "some" so advanced
	// Output:
	// expected "thing" at U+0020 ' ' 1,5-5 (5-5)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_not_Success() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.N{"foo"})
	c.Print() // not advanced, but also not <nil>
	s.Print() // not advanced at all
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_not_Fail() {
	s, _ := scan.New("some thing")
	_, err := s.Expect(z.N{"some"})
	fmt.Println(err)
	s.Print() // not advanced at all
	// Output:
	// unexpected "some" at U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_not_X_Fail() {
	s, _ := scan.New("some thing wonderful")
	_, err := s.Expect(z.X{z.N{'s'}, 'o'})
	fmt.Println(err) // sees the s so fails
	s.Print()        // not advanced
	// Output:
	// unexpected 's' at U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_not_X_Success() {
	s, _ := scan.New("some thing wonderful")
	c, _ := s.Expect(z.X{z.N{`n`}, z.I{`som`}})
	c.Print() // pointing to last in match 'm'
	s.Print() // advanced to next after match 'e'
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// U+0065 'e' 1,4-4 (4-4)
}

func ExampleExpect_to_Success_Mid() {
	s, _ := scan.New("some thing wonderful")
	c, _ := s.Expect(z.T{"wo"})
	c.Print() // "wo" not inc, same as "some thing ", so ' '
	s.Print() // advances to 'w'
	// Output:
	// U+0020 ' ' 1,11-11 (11-11)
	// U+0077 'w' 1,12-12 (12-12)
}

func ExampleExpect_avoid_Not_with_In() {
	s, _ := scan.New("some thing")
	s.Snap()
	c, _ := s.Expect(z.I{z.N{'s'}, z.R{'a', 'z'}})
	c.Print() // unexpected success
	s.Print() // advanced to 'o'
	s.Back()
	// use z.X instead
	_, err := s.Expect(z.X{z.N{'s'}, z.R{'a', 'z'}})
	fmt.Println(err)
	s.Print() // not advanced
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+006F 'o' 1,2-2 (2-2)
	// unexpected 's' at U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
}

/*
func ExampleExpect_seq_Success() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.P{"some", ' ', "thin"})
	c.Print() // same as "some thin", points at 'n'
	s.Print() // advanced to 'g'
	// Output:
	// U+006E 'n' 1,9-9 (9-9)
	// U+0067 'g' 1,10-10 (10-10)
}
*/

func ExampleExpect_seq_Fail() {
	s, _ := scan.New("some thing")
	_, err := s.Expect(z.X{"some", "thin"})
	fmt.Println(err)
	s.Print() // not advanced at all
	// Output:
	// expected rune 't' at U+0020 ' ' 1,5-5 (5-5)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_seq_Not_Success() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.X{"some", ` `, z.N{`T`}, "thin"})
	c.Print() // same as "some thin"
	s.Print() // advanced to next after ('g')
	// Output:
	// U+006E 'n' 1,9-9 (9-9)
	// U+0067 'g' 1,10-10 (10-10)
}

func ExampleExpect_seq_Not_Fail() {
	s, _ := scan.New("some Thing")
	_, err := s.Expect(z.X{"some", ' ', z.N{`T`}, "ignored"})
	fmt.Println(err)
	s.Print() // not advanced at all
	// Output:
	// unexpected "T" at U+0054 'T' 1,6-6 (6-6)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleExpect_token_ANY() {
	s, _ := scan.New("some thing wonderful")
	c, _ := s.Expect(tk.ANY)
	c.Print() // same as `s` or s.Scan()
	s.Print() // advances
	c, _ = s.Expect(tk.A)
	c.Print() // same as `o` or s.Scan()
	s.Print() // advances
	c, _ = s.Expect(tk.A1)
	// we'll skip tk.A2 - tk.A8
	s.Print()
	c, _ = s.Expect(tk.A9)
	s.Print() // should advance 9 to pos 13
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+006F 'o' 1,2-2 (2-2)
	// U+006F 'o' 1,2-2 (2-2)
	// U+006D 'm' 1,3-3 (3-3)
	// U+0065 'e' 1,4-4 (4-4)
	// U+006F 'o' 1,13-13 (13-13)
}

func ExampleExpect_any_Success() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.A{5})
	c.Print() // same as "some "
	s.Print() // advanced to next after ('t')
	// Output:
	// U+0074 't' 1,6-6 (6-6)
	// U+0068 'h' 1,7-7 (7-7)
}

func ExampleExpect_o_Optional_Success() {
	s, _ := scan.New("some thing")
	//c, _ := s.Expect(z.O{"thing", "some"})
	c, _ := s.Expect("some")
	c.Print() // same as "some", points to 'e'
	s.Print() // advances to space ' '
	// Output:
	// U+0065 'e' 1,4-4 (4-4)
	// U+0020 ' ' 1,5-5 (5-5)
}

func ExampleExpect_minimum_One() {
	s, _ := scan.New("sommme thing")
	start := s.Mark()
	s.ScanN(2)
	c, _ := s.Expect(z.M1{'m'}) // goggles up all three
	c.Print()
	s.Print()
	s.Jump(start)
	c, _ = s.Expect(z.M1{'s'}) // yep, just one
	c.Print()
	s.Print()
	// Output:
	// U+006D 'm' 1,5-5 (5-5)
	// U+0065 'e' 1,6-6 (6-6)
	// U+0073 's' 1,1-1 (1-1)
	// U+006F 'o' 1,2-2 (2-2)
}

func ExampleExpect_minimum() {
	s, _ := scan.New("sssoommme thing")
	c, _ := s.Expect(z.M{2, 's'})
	c.Print() // needs 2, but will consume all three to last 's'
	s.Print() // advances to next after ('o')
	// Output:
	// U+0073 's' 1,3-3 (3-3)
	// U+006F 'o' 1,4-4 (4-4)
}

func ExampleExpect_mMx() {
	s, _ := scan.New("sommme thing")
	s.Snap()
	s.ScanN(2)
	s.Print()
	s.Expect(z.MM{1, 3, 'm'}) // goggles up all three
	s.Print()
	s.Back()
	s.Expect(z.MM{1, 3, 's'}) // yep, at least one
	s.Print()
	_, err := s.Expect(z.MM{1, 3, 'X'}) // nope
	fmt.Println(err)
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// U+0065 'e' 1,6-6 (6-6)
	// U+006F 'o' 1,2-2 (2-2)
	// expected min 1, max 3 of 'X' at U+006F 'o' 1,2-2 (2-2)
}

func ExampleExpect_c() {
	s, _ := scan.New("sommme thing")
	s.Snap()
	s.ScanN(2)
	s.Print()
	s.Expect(z.C{3, 'm'}) // goggles up all three
	s.Print()
	s.Back()
	s.Expect(z.C{1, 's'}) // yes, but silly since 's' is easier
	s.Print()
	_, err := s.Expect(z.C{3, 'X'}) // nope
	fmt.Println(err)
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// U+0065 'e' 1,6-6 (6-6)
	// U+006F 'o' 1,2-2 (2-2)
	// expected rune 'X' at U+006F 'o' 1,2-2 (2-2)
}

func ExampleExpect_rng() {
	s, _ := scan.New("some thing")
	s.Scan()
	c1, _ := s.Expect(z.R{'l', 'p'})
	c1.Print()
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// U+006D 'm' 1,3-3 (3-3)
}

func FailHook(s *scan.R) error { return fmt.Errorf("imma fail") }

func ExampleExpect_hook() {

	// plain function signature
	WouldSave := scan.Hook(func(s *scan.R) error {
		fmt.Println("would save")
		return nil
	})

	// as scan.Hook
	WouldScan := scan.Hook(func(s *scan.R) error {
		s.Scan()
		return nil
	})

	// FailHook defined outside of Example function (see source)

	s, _ := scan.New("some thing")
	s.Scan()
	s.Expect(WouldSave)
	s.Print() // hook didn't advance
	s.Expect(WouldScan)
	s.Print() // hook advanced scan by one
	_, e := s.Expect(FailHook)
	fmt.Println(e)

	// Output:
	// would save
	// U+006F 'o' 1,2-2 (2-2)
	// U+006D 'm' 1,3-3 (3-3)
	// failhook: imma fail at U+006D 'm' 1,3-3 (3-3)

}

func ExampleExpect_to_Success() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.T{'e'})
	c.Print() // same as "som", points to 'm'
	s.Print() // scanned next after ('e')
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// U+0065 'e' 1,4-4 (4-4)
}

func ExampleExpect_to_Inclusive() {
	s, _ := scan.New("some thing")
	c, _ := s.Expect(z.Ti{'e'})
	c.Print() // same as "some", points to 'e'
	s.Print() // scanned next after (' ')
	// Output:
	// U+0065 'e' 1,4-4 (4-4)
	// U+0020 ' ' 1,5-5 (5-5)
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

func ExampleStr() {
	s, _ := scan.New("some thing")
	s.Str("some")
	s.Print()
	s.Str(" ", "th")
	s.Print()
	// Output:
	// U+0020 ' ' 1,5-5 (5-5)
	// U+0069 'i' 1,8-8 (8-8)
}

func ExampleAny() {
	s, _ := scan.New("some thing")
	s.Any(4)
	s.Print()
	// Output:
	// U+0020 ' ' 1,5-5 (5-5)
}

func ExampleOpt() {
	s, _ := scan.New("some thing")
	defer s.PrintPanic()
	s.Opt("S", "s")
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
}

func Example_all() {
	s, _ := scan.New("some thing")
	defer s.PrintPanic()
	s.Opt("S", "s")
	s.Str("ome", " ", "thi")
	s.Print()
	// Output:
	// U+006E 'n' 1,9-9 (9-9)
}

func Example() {
	s, _ := scan.New("some thing")
	defer s.PrintPanic()
	s.Opt("S", "s")
	s.Str("ome", " ", "thi")
	s.Print()
}
