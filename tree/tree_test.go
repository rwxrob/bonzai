package tree_test

import (
	"github.com/rwxrob/bonzai/tree"
)

func ExampleTree() {
	t := tree.New([]string{"foo", "bar"})
	t.Print()
	// Output:
	// {"Trunk":{"T":0,"V":""},"Types":["UNDEFINED","foo","bar"]}
}

// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
import (
	"fmt"

	"github.com/rwxrob/bonzai/tree"
)

func ExampleNode() {
	n := new(tree.Node)
	fmt.Printf(`Node:
%10v  T (type)
%10q  V (value)
%10v  B (upper branch)
%10v  P (previous peer)
%10v  N (next peer)
%10v  F (first sub node )
%10v  L (last sub node)
%10v  Types (index of type names)
	`, n.T, n.V, n.B, n.P, n.N, n.F, n.L, n.Types)
	//Output:
	// Node:
	//          0  T (type)
	//         ""  V (value)
	//      <nil>  B (upper branch)
	//      <nil>  P (previous peer)
	//      <nil>  N (next peer)
	//      <nil>  F (first sub node )
	//      <nil>  L (last sub node)
	//      <nil>  Types (index of type names)
}

// ----------------------- PrintAsJSON interface ----------------------

func ExampleNode_Print() {

	n := new(tree.Node)
	n.Print()
	n.T = 10
	n.Print()
	f := new(tree.Node)
	n.F = f
	n.Print()
	f.T = 20
	n.Print()

	// Output:
	// []
	// [10]
	// [10,[[]]]
	// [10,[[20]]]
}

func ExampleNode_JSONL() {

	n := new(tree.Node)
	n.T = 1
	f := new(tree.Node)
	n.F = f
	f.T = 2
	n.PrintL()
	n.Types = &[]string{"Foo", "Bar"}
	n.PrintL()

	// Output:
	// [1,[[2]]]
	//
}

/*
func ExampleNode_newref() {
	// empty and unknown
	zer := new(pegn.Node)
	zer.Print()
	// Output:
	// []
}

func ExampleNode_newnotype() {
	// empty and unknown
	foo := new(pegn.Node)
	foo.Value = "foo"
	foo.Print()
	// Output:
	// [0,"foo"]
}

func ExampleNode_children() {
	// empty and unknown
	some := new(pegn.Node)
	child := new(pegn.Node)
	child.Value = "thing"
	some.AppendChild(child)
	some.Print()
	// Output:
	// [0,[[0,"thing"]]]
}

func ExampleNode_AppendChild() {

	const (
		UNDEF = iota
		Mum
		Kid
	)

	// Value later ignored because container
	mum := new(pegn.Node)
	mum.Value = "mummy"
	mum.Type = Mum
	mum.Print()            // [1,"mummy"]
	mum.FirstChild.Print() // <nil>

	kid := new(pegn.Node)
	kid.Type = Kid
	kid.Value = "kid"
	mum.AppendChild(kid) // mum now container

	kid.Print()            // [2,"kid"]
	mum.FirstChild.Print() // [2,"kid"]
	mum.LastChild.Print()  // [2,"kid"]
	kid.Parent.Print()     // [1,[[2,"kid"]]]

	another := new(pegn.Node)
	another.Type = Kid
	another.Value = "another"
	mum.AppendChild(another)
	another.Parent.Print()  // [1,[[2,"kid"],[2,"another"]]]
	mum.LastChild.Print()   // [2,"another"]
	kid.NextSib.Print()     // [2,"another"]
	another.PrevSib.Print() // [2,"kid"]

	// Output:
	// [1,"mummy"]
	// <nil>
	// [2,"kid"]
	// [2,"kid"]
	// [2,"kid"]
	// [1,[[2,"kid"]]]
	// [1,[[2,"kid"],[2,"another"]]]
	// [2,"another"]
	// [2,"another"]
	// [2,"kid"]

}

func ExampleNode_AdoptFrom() {

	const (
		_ = iota
		Parent
		Child
	)

	mum := new(pegn.Node)
	mum.Type = Parent

	kid := new(pegn.Node)
	kid.Type = Child
	kid.Value = "kid"
	mum.AppendChild(kid)

	another := new(pegn.Node)
	another.Type = Child
	another.Value = "another"
	mum.AppendChild(another)

	unc := new(pegn.Node)
	unc.Type = Parent

	own := new(pegn.Node)
	own.Type = Child
	own.Value = "own"
	unc.AppendChild(own)

	mum.Print()
	unc.Print()
	unc.AdoptFrom(mum)
	mum.Print()
	unc.Print()

	// Output:
	// [1,[[2,"kid"],[2,"another"]]]
	// [1,[[2,"own"]]]
	// [1]
	// [1,[[2,"own"],[2,"kid"],[2,"another"]]]

}

func ExampleNode_InsertBeforeSelf_noparent() {

	// no parent is fine (but difficult later add parents)
	sib := new(pegn.Node)
	sib.Value = "sib"
	another := new(pegn.Node)
	another.Value = "another"
	sib.InsertBeforeSelf(another)
	sib.PrevSib.Print()         // [0,"another"]
	another.NextSib.Print()     // [0,"sib"]
	fmt.Println(sib.Parent)     // <nil>
	fmt.Println(another.Parent) // <nil>

	// [0,"another"]
	// [0,"sib"]
	// <nil>
	// <nil>

}

func ExampleNode_InsertBeforeSelf_parent() {

	// best to start with parent
	mum := new(pegn.Node)
	sib := new(pegn.Node)
	sib.Value = "sib"
	mum.AppendChild(sib)
	another := new(pegn.Node)
	another.Value = "another"
	sib.InsertBeforeSelf(another)
	sib.PrevSib.Print()     // [0,"another"]
	another.NextSib.Print() // [0,"sib"]
	sib.Parent.Print()      // [0,[[0,"another"],[0,"sib"]]]
	another.Parent.Print()  // [0,[[0,"another"],[0,"sib"]]]
	mum.FirstChild.Print()  // [0,"another"]
	mum.LastChild.Print()   // [0,"sib"]

	// Output:
	// [0,"another"]
	// [0,"sib"]
	// [0,[[0,"another"],[0,"sib"]]]
	// [0,[[0,"another"],[0,"sib"]]]
	// [0,"another"]
	// [0,"sib"]

}

func ExampleNode_AppendAfterSelf() {

	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	one.Value = "one"
	two.Value = "two"
	mum.AppendChild(one)
	one.AppendAfterSelf(two)
	mum.Print()

	// Output:
	// [0,[[0,"one"],[0,"two"]]]

}

func ExampleNode_RemoveSelf_first() {

	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	three := new(pegn.Node)

	two.Value = "two"
	three.Value = "three"

	mum.AppendChild(one)
	mum.AppendChild(two)
	mum.AppendChild(three)
	mum.Print()

	gr := new(pegn.Node)
	gr.Value = "grandchild"
	one.AppendChild(gr)

	orphan := one.RemoveSelf()
	orphan.Print()
	mum.Print()
	fmt.Println(two.PrevSib)
	mum.FirstChild.Print()
	mum.LastChild.Print()

	// Output:
	// [0,[[],[0,"two"],[0,"three"]]]
	// [0,[[0,"grandchild"]]]
	// [0,[[0,"two"],[0,"three"]]]
	// <nil>
	// [0,"two"]
	// [0,"three"]

}

func ExampleNode_RemoveSelf_middle() {

	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	three := new(pegn.Node)

	one.Value = "one"
	two.Value = "two"
	three.Value = "three"

	mum.AppendChild(one)
	mum.AppendChild(two)
	mum.AppendChild(three)
	mum.Print()

	orphan := two.RemoveSelf()
	orphan.Print()
	mum.Print()
	one.NextSib.Print()
	three.PrevSib.Print()
	mum.FirstChild.Print()
	mum.LastChild.Print()

	// Output:
	// [0,[[0,"one"],[0,"two"],[0,"three"]]]
	// [0,"two"]
	// [0,[[0,"one"],[0,"three"]]]
	// [0,"three"]
	// [0,"one"]
	// [0,"one"]
	// [0,"three"]

}

func ExampleNode_RemoveSelf_last() {

	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	three := new(pegn.Node)

	mum.Value = "mummy"
	one.Value = "one"
	two.Value = "two"
	three.Value = "three"

	mum.AppendChild(one)
	mum.AppendChild(two)
	mum.AppendChild(three)
	mum.Print()

	orphan := three.RemoveSelf()
	orphan.Print()
	mum.Print()
	one.NextSib.Print()
	fmt.Println(two.NextSib)
	mum.LastChild.Print()

	// Output:
	// [0,[[0,"one"],[0,"two"],[0,"three"]]]
	// [0,"three"]
	// [0,[[0,"one"],[0,"two"]]]
	// [0,"two"]
	// <nil>
	// [0,"two"]

}

func ExampleNode_ReplaceSelf_first() {

	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	three := new(pegn.Node)

	two.Value = "two"
	three.Value = "three"

	mum.AppendChild(one)
	mum.AppendChild(two)
	mum.AppendChild(three)

	gr := new(pegn.Node)
	gr.Value = "grandchild"
	one.AppendChild(gr)

	mum.Print()

	nw := new(pegn.Node)
	nw.Value = "replacement"

	orphan := one.ReplaceSelf(nw)
	orphan.Print()

	mum.Print()
	mum.FirstChild.Print()
	mum.LastChild.Print()

	fmt.Println(two.PrevSib)
	fmt.Println(nw.NextSib)
	fmt.Println(nw.PrevSib)

	// Output:
	// [0,[[0,[[0,"grandchild"]]],[0,"two"],[0,"three"]]]
	// [0,[[0,"grandchild"]]]
	// [0,[[0,"replacement"],[0,"two"],[0,"three"]]]
	// [0,"replacement"]
	// [0,"three"]
	// [0,"replacement"]
	// [0,"two"]
	// <nil>

}

func ExampleNode_ReplaceSelf_middle() {

	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	three := new(pegn.Node)

	one.Value = "one"
	three.Value = "three"

	mum.AppendChild(one)
	mum.AppendChild(two)
	mum.AppendChild(three)

	gr := new(pegn.Node)
	gr.Value = "grandchild"
	two.AppendChild(gr)

	mum.Print()

	nw := new(pegn.Node)
	nw.Value = "replacement"

	orphan := two.ReplaceSelf(nw)
	orphan.Print()

	mum.Print()
	mum.FirstChild.Print()
	mum.LastChild.Print()

	fmt.Println(one.NextSib)
	fmt.Println(nw.PrevSib)
	fmt.Println(nw.NextSib)
	fmt.Println(three.PrevSib)

	// Output:
	// [0,[[0,"one"],[0,[[0,"grandchild"]]],[0,"three"]]]
	// [0,[[0,"grandchild"]]]
	// [0,[[0,"one"],[0,"replacement"],[0,"three"]]]
	// [0,"one"]
	// [0,"three"]
	// [0,"replacement"]
	// [0,"one"]
	// [0,"three"]
	// [0,"replacement"]

}

func ExampleNode_ReplaceSelf_last() {

	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	three := new(pegn.Node)

	one.Value = "one"
	two.Value = "two"

	mum.AppendChild(one)
	mum.AppendChild(two)
	mum.AppendChild(three)

	gr := new(pegn.Node)
	gr.Value = "grandchild"
	three.AppendChild(gr)

	mum.Print()

	nw := new(pegn.Node)
	nw.Value = "replacement"

	orphan := three.ReplaceSelf(nw)
	orphan.Print()

	mum.Print()
	mum.FirstChild.Print()
	mum.LastChild.Print()

	fmt.Println(two.NextSib)
	fmt.Println(nw.PrevSib)
	fmt.Println(nw.NextSib)

	// Output:
	// [0,[[0,"one"],[0,"two"],[0,[[0,"grandchild"]]]]]
	// [0,[[0,"grandchild"]]]
	// [0,[[0,"one"],[0,"two"],[0,"replacement"]]]
	// [0,"one"]
	// [0,"replacement"]
	// [0,"replacement"]
	// [0,"two"]
	// <nil>

}

func ExampleNode_Children() {
	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	three := new(pegn.Node)

	one.Value = "one"
	two.Value = "two"
	three.Value = "three"

	fmt.Println(mum.Children())

	mum.AppendChild(one)
	mum.AppendChild(two)
	mum.AppendChild(three)

	fmt.Println(mum.Children())

	// Output:
	// []
	// [[0,"one"] [0,"two"] [0,"three"]]

}

func ExampleNode_Visit() {

	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	three := new(pegn.Node)

	one.Value = "one"
	two.Value = "two"
	three.Value = "three"

	mum.AppendChild(one)
	mum.AppendChild(two)
	mum.AppendChild(three)

	mum.Visit(
		func(n *pegn.Node) interface{} {
			n.Print()
			return nil
		}, nil)

	// Output:
	// [0,[[0,"one"],[0,"two"],[0,"three"]]]
	// [0,"one"]
	// [0,"two"]
	// [0,"three"]

}

func ExampleNode_VisitAsync() {

	mum := new(pegn.Node)
	one := new(pegn.Node)
	two := new(pegn.Node)
	three := new(pegn.Node)

	one.Value = "one"
	two.Value = "two"
	three.Value = "three"

	mum.AppendChild(one)
	mum.AppendChild(two)
	mum.AppendChild(three)

	justprint := func(n *pegn.Node) interface{} {
		n.Print()
		return nil
	}

	mum.VisitAsync(justprint, 5, nil)

	// Unordered Output:
	// [0,"three"]
	// [0,"one"]
	// [0,[[0,"one"],[0,"two"],[0,"three"]]]
	// [0,"two"]

}

/*
func ExampleNode_UnmarshalJSON_simple() {
	data := `[1,"something"]`
	n := new(pegn.Node)
	err := n.UnmarshalJSON([]byte(data))
	n.Print()
	fmt.Println(err)
	// Output:
	// [1,"something"]
	// <nil>
}

func ExampleNode_UnmarshalJSON_complex() {

	data := `[1,[[12,[[21,"Grammar"],[24,[[28,[[32,[[50,"some"]]],[32,[[21,"Thing"]]],[32,[[50,"+"]]],[32,[[21,"Else"]]]]],[28,[[32,[[21,"Other"]]]]]]]]],[58,"\n"],[11,[[21,"Thing"],[24,[[28,[[32,[[50,"thing"]]]]]]]]],[58,"\n"],[11,[[21,"Else"],[24,[[28,[[32,[[50,"else"]]]]]]]]],[58,"\n"],[11,[[21,"Other"],[24,[[28,[[32,[[50,"other"]]]]]]]]],[58,"\n"]]]`

	result := `
["Grammar", [
  ["CheckDef", [
    ["CheckId", "Grammar"],
    ["Expression", [
      ["Sequence", [
        ["Plain", [
          ["String", "some"]
        ]],
        ["Plain", [
          ["CheckId", "Thing"]
        ]],
        ["Plain", [
          ["String", "+"]
        ]],
        ["Plain", [
          ["CheckId", "Else"]
        ]]
      ]],
      ["Sequence", [
        ["Plain", [
          ["CheckId", "Other"]
        ]]
      ]]
    ]]
  ]],
  ["EndLine", "\n"],
  ["SchemaDef", [
    ["CheckId", "Thing"],
    ["Expression", [
      ["Sequence", [
        ["Plain", [
          ["String", "thing"]
        ]]
      ]]
    ]]
  ]],
  ["EndLine", "\n"],
  ["SchemaDef", [
    ["CheckId", "Else"],
    ["Expression", [
      ["Sequence", [
        ["Plain", [
          ["String", "else"]
        ]]
      ]]
    ]]
  ]],
  ["EndLine", "\n"],
  ["SchemaDef", [
    ["CheckId", "Other"],
    ["Expression", [
      ["Sequence", [
        ["Plain", [
          ["String", "other"]
        ]]
      ]]
    ]]
  ]],
  ["EndLine", "\n"]
]]
`

	n := new(pegn.Node)
	err := n.UnmarshalJSON([]byte(data))
	fmt.Println(n.String() == result)
	fmt.Println(err)
	// Output:
	// true
	// <nil>
}
*/
