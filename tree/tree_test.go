package tree_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/rwxrob/bonzai/tree"
)

func ExampleTree() {
	t := tree.New([]string{"foo", "bar"})
	t.Print()
	// Output:
	// {"Root":[1],"Types":["UNKNOWN","foo","bar"]}
}

func ExampleTree_Log() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	t := tree.New([]string{"foo", "bar"})
	t.Log()
	// Output:
	// {"Root":[1],"Types":["UNKNOWN","foo","bar"]}
}

func ExampleTree_Seed() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	t := tree.New([]string{"foo", "bar"})
	s := t.Seed(2, "")
	fmt.Println(s.IsRoot(), s.IsDetached(), s.IsNull())
	// Output:
	// true true true
}

func TestParse(t *testing.T) {
	types := []string{"root", "foo", "bar"}
	tests := [][]string{
		{`  []`, `[]`},
		/*
			{`[]   `, `[]`},
			{"\n[]", `[]`},
			{"[]\n", `[]`},
			{"\t[]", `[]`},
			{`[]`},
			{`["root"]`, `[1]`},
			{`[2]`},
			{`[2,"some val"]`},
		*/
	}
	for _, tt := range tests {
		tr, err := tree.Parse(tt[0], types)
		if err != nil {
			t.Error(err)
		}
		want := tt[0]
		got := tr.Root.String()
		if len(tt) > 1 {
			want = tt[1]
		}
		if got != want {
			t.Errorf("\nwant: %v\ngot:  %v\n", want, got)
		}
	}
}
