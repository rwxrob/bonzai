package scanner_test

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/rwxrob/bonzai/pkg/core/scanner"
)

func ExampleS_init() {

	// * extremely minimal initialization
	// * order guaranteed never to change

	s := scanner.New(`some thing`)
	fmt.Println(s)

}

func ExampleS_package_Trace() {

	// take over stderr just for this test
	defer log.SetFlags(log.Flags())
	defer log.SetOutput(os.Stderr)
	defer func() { scanner.Trace = 0 }()
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	s := scanner.New(`foo`)
	scanner.Trace++
	s.Scan()
	s.Scan()
	s.Scan()

	// Output:
	// 'f' 0-1 "oo"
	// 'o' 1-2 "o"
	// 'o' 2-3 ""

}

func ExampleS_Scan() {

	s := scanner.New(`foo`)
	s.Print() // equivalent of a "zero value"

	fmt.Println(s.Scan())
	s.Print()
	fmt.Println(s.Scan())
	s.Print()
	fmt.Println(s.Scan())
	s.Print()
	fmt.Println(s.Scan()) // does not advance
	s.Print()             // same as before

	// Output:
	// '\x00' 0-0 "foo"
	// true
	// 'f' 0-1 "oo"
	// true
	// 'o' 1-2 "o"
	// true
	// 'o' 2-3 ""
	// false
	// 'o' 2-3 ""

}

func ExampleS_Scan_loop() {
	s := scanner.New(`abcdefgh`)
	for s.Scan() {
		fmt.Print(string(s.Rune()))
		if !s.Finished() {
			fmt.Print("-")
		}
	}
	// Output:
	// a-b-c-d-e-f-g-h
}

func ExampleS_Scan_risky_Jump() {

	s := scanner.New(`foo1234`)
	fmt.Println(s.Scan())
	s.Print()
	s.E += 2              //😟 WARNING: s.R and s.B, not yet updated!
	fmt.Println(s.Scan()) //😊 s.R and s.B now updated
	s.Print()

	// Output:
	// true
	// 'f' 0-1 "oo1234"
	// true
	// '1' 3-4 "234"
}

func ExampleS_Is() {

	s := scanner.New(`foo`)

	s.Scan() // never forget to scan with Is (use Peek otherwise)

	fmt.Println(s.Is("fo"))
	fmt.Println(s.Is("bar"))

	// Output:
	// true
	// false
}

func ExampleS_Peek() {

	s := scanner.New(`foo`)

	fmt.Println(s.Peek("fo"))
	s.Scan()
	fmt.Println(s.Peek("fo"))
	fmt.Println(s.Peek("oo"))

	// Output:
	// true
	// false
	// true
}

func ExampleS_Is_not() {

	s := scanner.New("\r\n")

	s.Scan() // never forget to scan with Is (use Peek otherwise)

	fmt.Println(s.Is("\r"))
	fmt.Println(s.Is("\r\n"))
	fmt.Println(s.Is("\n"))

	// Output:
	// true
	// true
	// false

}

func ExampleS_Match() {

	s := scanner.New(`foo`)

	s.Scan() // never forget to scan (use PeekMatch otherwise)

	f := regexp.MustCompile(`f`)
	F := regexp.MustCompile(`F`)
	o := regexp.MustCompile(`o`)

	fmt.Println(s.Match(f))
	fmt.Println(s.Match(F))
	fmt.Println(s.Match(o))

	// Output:
	// 1
	// -1
	// -1
}

func ExampleS_Pos() {

	//😟 WARNING: uses risky jumps (assigning s.E)

	s := scanner.New("one line\nand another\r\nand yet another")

	s.E = 2
	s.Pos().Print()

	s.E = 0
	s.Scan()
	s.Scan()
	s.Pos().Print()

	s.E = 12
	s.Pos().Print()

	s.E = 27
	s.Pos().Print()

	// Output:
	// U+006E 'n' 1,2-2 (2-2)
	// U+006E 'n' 1,2-2 (2-2)
	// U+0064 'd' 2,3-3 (12-12)
	// U+0079 'y' 3,5-5 (27-27)

}

func ExampleS_Positions() {

	s := scanner.New("one line\nand another\r\nand yet another")

	for _, p := range s.Positions(2, 12, 27) {
		p.Print()
	}

	// Output:
	// U+006E 'n' 1,2-2 (2-2)
	// U+0064 'd' 2,3-3 (12-12)
	// U+0079 'y' 3,5-5 (27-27)

}

/*
func ExampleS_Report() {

	//😟 WARNING: uses risky jumps (assigning s.E)

	defer log.SetFlags(log.Flags())
	defer log.SetOutput(os.Stderr)
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	s := scanner.New("one line\nand another\r\nand yet another")

	s.Scan()
	s.Report()

	s.E = 14
	s.Report()

	s.Error("sample error")
	s.Report()

	// Output:
	// U+006F 'o' 1,1-1 (1-1)
	// U+0061 'a' 2,5-5 (14-14)
	// error: sample error at U+0061 'a' 2,5-5 (14-14)

}
*/

func ExampleS_Finished() {

	s := scanner.New(`foo`)

	s.Print()

	s.Scan()
	s.Print()
	fmt.Println(s.Finished())

	s.Scan()
	s.Print()
	fmt.Println(s.Finished())

	s.Scan()
	s.Print()
	fmt.Println(s.Finished())

	// Output:
	// '\x00' 0-0 "foo"
	// 'f' 0-1 "oo"
	// false
	// 'o' 1-2 "o"
	// false
	// 'o' 2-3 ""
	// true

}

func ExampleMark() {

	s := scanner.New(`foo`)

	m := s.Mark()
	fmt.Println(m)
	s.Print()
	s.Scan()
	s.Print()
	mm := s.Mark()
	fmt.Println(mm)
	s.Goto(m)
	s.Print()

	// Output:
	// '\x00' 0-0
	// '\x00' 0-0 "foo"
	// 'f' 0-1 "oo"
	// 'f' 0-1
	// '\x00' 0-0 "foo"

}
