package tree_test

import (
	"log"
	"os"

	"github.com/rwxrob/bonzai/tree"
)

func ExampleTree() {
	t := tree.New([]string{"foo", "bar"})
	t.Print()
	// Output:
	// {"Trunk":[1],"Types":["UNKNOWN","foo","bar"]}
}

func ExampleTree_Log() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	t := tree.New([]string{"foo", "bar"})
	t.Log()
	// Output:
	// {"Trunk":[1],"Types":["UNKNOWN","foo","bar"]}
}

func ExampleNode_Log() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Log()
	// Output:
	// [1]
}

func ExampleNode_Info() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Info()

	// Type:       1
	// Value:      ""
	// IsRoot:     true
	// IsDetached: true
	// IsLeaf:     false
	// IsBranch:   false
	// IsNull:     true
}

func ExampleNode_leaf() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Print()
	n.V = "something"
	n.Print()
	n.T = 0
	n.V = ""
	n.Print()
	// Output:
	// [1]
	// [1,"something"]
	// []
}

func ExampleNode_branch() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.NewUnder(2, "something")
	n.Print()
	n.PrettyPrint()
	// Output:
	// [1,[[2,"something"]]]
	// ["foo", [
	//   ["bar", "something"]
	// ]]
}

func ExampleNode_NewRight() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Print()
	t := n.NewRight(2, "some")
	t.Print()
	n.Right().Print()
	//Output:
	// [1]
	// [2,"some"]
	// [2,"some"]
}

func ExampleNode_NewLeft() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Print()
	t := n.NewLeft(2, "some")
	t.Print()
	n.Left().Print()
	//Output:
	// [1]
	// [2,"some"]
	// [2,"some"]
}

func ExampleNode_NewUnder() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Print()
	t := n.NewUnder(2, "some")
	t.Print()
	n.FirstUnder().Print()
	n.LastUnder().Print()
	n.Print()
	//Output:
	// [1]
	// [2,"some"]
	// [2,"some"]
	// [2,"some"]
	// [1,[[2,"some"]]]
}

func ExampleNode_Graft() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Print()
	t := n.NewUnder(2, "some")
	t.Print()
	n.FirstUnder().Print()
	n.LastUnder().Print()
	//Output:
	// [1]
	// [2,"some"]
	// [2,"some"]
	// [2,"some"]
}

/*
func (n *Node) Graft(c *Node) *Node {
func (n *Node) GraftRight(r *Node) *Node {
func (n *Node) GraftLeft(l *Node) *Node {
func (n *Node) GraftUnder(c *Node) *Node {
func (n *Node) Prune() *Node {
func (n *Node) Take(from *Node) {
func (n *Node) under() []*Node {
func (n *Node) Visit(act Action, rvals chan interface{}) {
func (n *Node) VisitAsync(act Action, lim int, rvals chan interface{}) {
*/
