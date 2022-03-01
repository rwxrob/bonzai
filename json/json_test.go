package json_test

import (
	"github.com/rwxrob/bonzai/fn"
	"github.com/rwxrob/bonzai/json"
)

func ExampleEscape() {
	set := fn.A[string]{
		`<`, `>`, `&`, `"`, `'`,
		"\t", "\b", "\f", "\n", "\r",
		"\\", "\"", "💢", "д",
	}
	set.Map(json.Escape).Print()

	/*
		// uncomment to see how bad the default is
		set.Map(func(i string) string {
			byt, _ := bork.Marshal(i)
			return string(byt)
		}).Print()
	*/

	// Output:
	// <>&\"'\t\b\f\n\r\\\"💢д
}

/*
func ExampleScanner_Errorf() {
	s, _ := json.NewScanner(`{"some": "thing"}`)
	fmt.Println(s.Errorf("oh no"))
	// Output:
	// scanner: oh no (line:0 col:0)
}

func ExampleJsonParser_EmptyString() {
	n := new(pegn.Node)
	s := ""
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// []
}

func ExampleJsonParser_EmptyArray() {
	n := new(pegn.Node)
	s := `[]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// []
}

func ExampleJsonParser_NoTypeArray() {
	n := new(pegn.Node)
	s := `[0,""]     `
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// []
}

func ExampleJsonParser_ChildNoValue() {
	n := new(pegn.Node)
	s := `[1, [[2],[3]]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// [1,[[2],[3]]]
}

func ExampleJsonParser_NoValue() {
	n := new(pegn.Node)
	s := `[1]     `
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// [1]
}

func ExampleJsonParser_TextNode() {
	n := new(pegn.Node)
	s := `[1, "Test"]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)
	if err != nil {
		fmt.Println(err)
	}

	n.Print()

	// Output:
	// [1,"Test"]
}

func ExampleJsonParser_Emoji() {
	n := new(pegn.Node)
	s := `[1, "Feeling 👺."]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)
	if err != nil {
		fmt.Println(err)
	}

	n.Print()

	// Output:
	// [1,"Feeling 👺."]
}

func ExampleJsonParser_EmptyChildNode() {
	n := new(pegn.Node)
	s := `[1, [[0, ""]]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)
	if err != nil {
		fmt.Println(err)
	}

	n.Print()

	// Output:
	// [1,[[]]]
}

func ExampleJsonParser_SingleChildNode() {
	n := new(pegn.Node)
	s := `[1, [[2, "test1"]]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)
	if err != nil {
		fmt.Println(err)
	}

	n.Print()

	// Output:
	// [1,[[2,"test1"]]]
}

func ExampleJsonParser_TwoChildNodes() {
	n := new(pegn.Node)
	s := `[1, [[2, "test1"], [2, "test2"]]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)
	if err != nil {
		fmt.Println(err)
	}

	n.Print()

	// Output:
	// [1,[[2,"test1"],[2,"test2"]]]
}

func ExampleJsonParser_ManyChildNodes() {
	n := new(pegn.Node)
	s := `[1, [[2, "test1"], [2, "test2"],[2,"test3"]]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)
	if err != nil {
		fmt.Println(err)
	}

	n.Print()

	// Output:
	// [1,[[2,"test1"],[2,"test2"],[2,"test3"]]]
}

func ExampleJsonParser_GrandChildrenNodes() {
	n := new(pegn.Node)
	s := `[1,[[3,[[2,"test4"]]], [2, "test2"],[2,"test3"]]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)
	if err != nil {
		fmt.Println(err)
	}

	n.Print()

	// Output:
	// [1,[[3,[[2,"test4"]]],[2,"test2"],[2,"test3"]]]
}

func ExampleJsonParser_MultipleGrandChildrenNodes() {
	n := new(pegn.Node)
	s := `[1,[[3,[[2,"test4"], [2,"test5"]]], [2, "test2"],[2,"test3"]]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)
	if err != nil {
		fmt.Println(err)
	}

	n.Print()

	// Output:
	// [1,[[3,[[2,"test4"],[2,"test5"]]],[2,"test2"],[2,"test3"]]]
}

func ExampleJsonParser_UnterminatedArray() {
	n := new(pegn.Node)
	s := `[
            0,
            ""
         `
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: expecting ']' after node value line: 3, col: 9
}

func ExampleJsonParser_UnterminatedString() {
	n := new(pegn.Node)
	s := `	[0, "     ]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: unterminated string line: 0, col: 12
}

func ExampleJsonParser_ValueAfterChildError() {
	n := new(pegn.Node)
	s := `[1, [[2, "Test"], "Value"]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: expecting child line: 0, col: 26
}

func ExampleJsonParser_InvalidToken() {
	n := new(pegn.Node)
	s := `[2,&"Test"]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: invalid token line: 0, col: 4
}

func ExampleJsonParser_NoUnicodeEscape() {
	n := new(pegn.Node)
	s := `[2,"\u0312"]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: unicode escape not supported in ast json line: 0, col: 6
}

func ExampleJsonParser_NoNewlineInString() {
	n := new(pegn.Node)
	s := `[2,"Test
                     String"]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: invalid newline in string line: 1, col: 0
}

func ExampleJsonParser_EscapeCodes() {
	n := new(pegn.Node)
	s := `[2,"\b\\\n\r\f\/\""]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// [2,"\b\\\n\r\f/\""]
}

func ExampleJsonParser_ExpectNodeType() {
	n := new(pegn.Node)
	s := `[[],"test"]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: expecting node type line: 0, col: 3
}

func ExampleJsonParser_MalformedAfterChild() {
	n := new(pegn.Node)
	s := `[1, [[2, "Test"] [[2,"Test2"]]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: expecting comma or end of children array line: 0, col: 19
}

func ExampleJsonParser_NoCommaAfterType() {
	n := new(pegn.Node)
	s := `[0 ""]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: invalid input, expecting comma after node type line: 0, col: 6
}

func ExampleJsonParser_UnterminatedSingleLineArray() {
	n := new(pegn.Node)
	s := `[1, [[2, "test1"], [2, "test2",[2,"test3"]]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: expecting ']' after node value line: 0, col: 32
}

func ExampleJsonParser_UnterminatedNode() {
	n := new(pegn.Node)
	s := `[1, [[2, "test1"], [2, "test2"],[2,"test3"]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: missing ending ']' line: 0, col: 44
}

func ExampleJsonParser_DanglingRightBracket() {
	n := new(pegn.Node)
	s := `	[0, "   "  ]]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	err := p.Parse(n)

	fmt.Println(err)

	// Output:
	// jsonparser: extra characters after json string line: 0, col: 14
}

func ExampleJsonParser_EscapeNewline() {
	n := new(pegn.Node)
	s := `[0, "lit\neral"]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// [0,"lit\neral"]
}

func ExampleJsonParser_StartWithEscape() {
	n := new(pegn.Node)
	s := `[0, "\nliteral"]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// [0,"\nliteral"]
}

func ExampleJsonParser_OnlyEscapes() {
	n := new(pegn.Node)
	s := `[0, "\t\n\r"]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// [0,"\t\n\r"]
}

func ExampleJsonParser_EndEscape() {
	n := new(pegn.Node)
	s := `[0, "Test\n"]`
	p := new(pegn.JsonParser)
	p.Init([]byte(s))
	_ = p.Parse(n)

	n.Print()

	// Output:
	// [0,"Test\n"]
}
*/
