package mark_test

import (
	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/scan"
)

func ExamplePlain() {
	s := scan.New("some thing here")
	//s.TraceX()
	s.X(mark.Plain)
	s.Print()
	s.Tree.Root.Print()
	// Output:
	// <EOD>
	// {"T":1,"N":[{"T":10,"V":"some thing here"}]}
}

func ExampleBracketed() {
	s := scan.New("<thing>")
	s.TraceX()
	s.X(mark.Bracketed)
	s.Print()
	s.Tree.Root.Print()
	// Output:
	// <EOD>
	// {"T":1,"N":[{"T":10,"V":"some thing here"}]}
}
