package tree_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/pkg/core/ds/tree"
)

func ExampleNode_type_is_Integer() {
	n := new(tree.Node[any])
	n.T = 1
	n.Print()
	//n.T = "one" // BOOM!

	// Output:
	// {"T":1}
}

func ExampleNode_value_Depends_on_Instantiation() {

	z := new(tree.Node[any])
	z.Print() // values that are equiv of zero value for type are omitted

	// flexible
	n := new(tree.Node[any])
	n.V = 1
	n.Print()
	n.V = true
	n.Print()
	n.V = "foo"
	n.Print()

	// strict
	m := new(tree.Node[int])
	m.V = 1
	m.Print()
	//m.V = true // BOOM!

	// Output:
	// {"T":0}
	// {"T":0,"V":1}
	// {"T":0,"V":true}
	// {"T":0,"V":"foo"}
	// {"T":0,"V":1}
}

func ExampleNode_Init() {

	// create and print a brand new one
	n := new(tree.Node[any])
	n.Print()

	// add something to it
	n.V = "something"
	n.Print()

	// initialize it back to "empty"
	n.Init()
	n.Print()

	// now try with something with a tricker zero value
	n.V = func() {}
	// log.Println(n.V) // "0x5019e0" yep it's there
	// n.Print()     // would log error and fail to marshal JSON

	// check it's properties
	fmt.Println(n.P != nil, n.Count)

	// Output:
	// {"T":0}
	// {"T":0,"V":"something"}
	// {"T":0}
	// false 0
}

func ExampleNode_String() {
	t := new(tree.Node[any])
	t.V = `<foo>` // <> not be escaped by encoding/json
	t.Print()
	// Output:
	// {"T":0,"V":"<foo>"}
}

func ExampleNode_properties() {

	// Nodes have these properties updating every time
	// their state is changed so that queries need not
	// to the checks again later.

	// initial state
	n := new(tree.Node[any])
	fmt.Println("n:", n.P == nil, n.V, n.Count)
	u := n.Add(1, nil)
	fmt.Println("n:", n.P == nil, n.V, n.Count)
	fmt.Println("u:", u.P == nil, u.V, u.Count)

	// make an edge node
	u.V = "something"

	// break edge by forcing it to have nodes and a value (discouraged)
	u.Add(9001, "muhaha")
	fmt.Println("u:", u.P == nil, u.V, u.Count)

	// Output:
	// n: true <nil> 0
	// n: true <nil> 1
	// u: false <nil> 0
	// u: false something 1

}

func ExampleNode_Nodes() {

	// create the first tree
	n := new(tree.Node[any])
	n.Add(1, nil)
	n.Add(2, nil)
	fmt.Println(n.Nodes(), n.Count)

	// and another added under it
	m := n.Add(3, nil)
	m.Add(3, nil)
	m.Add(3, nil)
	fmt.Println(m.Nodes(), m.Count)
	fmt.Println(n.Nodes(), n.Count)

	// Output:
	// [{"T":1} {"T":2}] 2
	// [{"T":3} {"T":3}] 2
	// [{"T":1} {"T":2} {"T":3,"N":[{"T":3},{"T":3}]}] 3

}

func ExampleNode_Cut_middle() {
	n := new(tree.Node[any])
	n.Add(1, nil)
	c := n.Add(2, nil)
	n.Add(3, nil)
	n.Print()
	fmt.Println(n.Count)
	x := c.Cut()
	n.Print()
	fmt.Println(n.Count)
	x.Print()
	// Output:
	// {"T":0,"N":[{"T":1},{"T":2},{"T":3}]}
	// 3
	// {"T":0,"N":[{"T":1},{"T":3}]}
	// 2
	// {"T":2}
}

func ExampleNode_Cut_first() {
	n := new(tree.Node[any])
	c := n.Add(1, nil)
	n.Add(2, nil)
	n.Add(3, nil)
	n.Print()
	x := c.Cut()
	n.Print()
	x.Print()
	// Output:
	// {"T":0,"N":[{"T":1},{"T":2},{"T":3}]}
	// {"T":0,"N":[{"T":2},{"T":3}]}
	// {"T":1}
}

func ExampleNode_Cut_last() {
	n := new(tree.Node[any])
	n.Add(1, nil)
	n.Add(2, nil)
	c := n.Add(3, nil)
	n.Print()
	x := c.Cut()
	n.Print()
	x.Print()
	// Output:
	// {"T":0,"N":[{"T":1},{"T":2},{"T":3}]}
	// {"T":0,"N":[{"T":1},{"T":2}]}
	// {"T":3}
}

func ExampleNode_Take() {

	// build up the first
	n := new(tree.Node[any])
	n.T = 10
	n.Add(1, nil)
	n.Add(2, nil)
	n.Add(3, nil)
	n.Print()
	fmt.Println(n.Count)

	// now take them over

	m := new(tree.Node[any])
	m.T = 20
	m.Print()
	fmt.Println(m.Count)
	m.Take(n)
	m.Print()
	fmt.Println(m.Count)
	n.Print()
	fmt.Println(n.Count)

	// Output:
	// {"T":10,"N":[{"T":1},{"T":2},{"T":3}]}
	// 3
	// {"T":20}
	// 0
	// {"T":20,"N":[{"T":1},{"T":2},{"T":3}]}
	// 3
	// {"T":10}
	// 0

}

func ExampleNode_WalkLevels() {
	n := new(tree.Node[any])
	n.Add(1, nil).Add(11, nil)
	n.Add(2, nil).Add(22, nil)
	n.Add(3, nil).Add(33, nil)
	n.WalkLevels(func(c *tree.Node[any]) { fmt.Print(c.T, " ") })
	// Output:
	// 0 1 2 3 11 22 33
}

func ExampleNode_WalkDeepPre() {
	n := new(tree.Node[any])
	n.Add(1, nil).Add(11, nil)
	n.Add(2, nil).Add(22, nil)
	n.Add(3, nil).Add(33, nil)
	n.WalkDeepPre(func(c *tree.Node[any]) { fmt.Print(c.T, " ") })
	// Output:
	// 0 1 11 2 22 3 33
}

func ExampleNode_Morph() {
	n := new(tree.Node[any])
	n.Add(2, "some")
	m := new(tree.Node[any])
	m.Morph(n)
	n.Print()
	m.Print()
	// Output:
	// {"T":0,"N":[{"T":2,"V":"some"}]}
	// {"T":0,"N":[{"T":2,"V":"some"}]}
}

func ExampleNode_Copy() {
	n := new(tree.Node[any])
	n.Add(2, "some")

	c := n.Copy()
	c.Add(3, "new").Add(4, "deep")

	// 	log.Print("Original -------------------------")
	// 	n.LogRefs()
	// 	log.Print("Copy -----------------------------")
	// 	c.LogRefs()

	fmt.Println(&n != &c)
	n.Print()
	c.Print()

	// Output:
	// true
	// {"T":0,"N":[{"T":2,"V":"some"}]}
	// {"T":0,"N":[{"T":2,"V":"some"},{"T":3,"V":"new","N":[{"T":4,"V":"deep"}]}]}

}
