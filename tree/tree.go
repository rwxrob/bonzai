// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package tree

import (
	"encoding/json"
	"fmt"
	"log"
)

// Tree is an encapsulating struct for containing the Types that go with
// a given rooted node tree of Nodes. The constructor (New) also ensures
// that every Node knows about the Tree to which it belongs.
type Tree struct {
	Trunk Node
	Types []string
}

// New creates and initializes a new try and returns it's root Node.
// Calling this constructor is preferred over creating Tree references
// directly since it requires specifying the types used in the tree.
// Some uses, however, may not know the types in advance and need to
// assign them later to Tree.Types.
func New(types []string) *Tree {
	t := new(Tree)
	t.Types = []string{"UNDEFINED"}
	t.Types = append(t.Types, types...)
	t.Trunk.tree = t
	return t
}

// JSON implements PrintAsJSON multi-line, 2-space indent JSON output.
func (s *Tree) JSON() string { b, _ := json.Marshal(s); return string(b) }

// Print implements PrintAsJSON.
func (s *Tree) Print() { fmt.Println(s.JSON()) }

// Log implements PrintAsJSON.
func (s *Tree) Log() { log.Print(s.JSON()) }

/*
// Node creates and initializes a new Node without attaching it to the
// tree at all. Node returns nil and an error if the Node type is not
// compatible with the current tree or the value is somehow inconsistent
// with the type specified.
func (t *Tree) (typ int, v string) (*Node, error) {
	n := new(Node)
	if !(0 <= typ && typ < len(t.Types)) {
		return nil, fmt.Errorf("bad type, must be 0 <= type <= %v", len(t.Types))
	}
	return n, nil
}
*/

// Nodes are for constructing rooted node trees of typed strings based
// on the tenet of the UNIX philosophy that suggests to focus on
// parsable text above all and converting when needed later.
//
// Branch or Leaf
//
// A Node can either be a "branch" or a "leaf" but not both. Branches
// have other leaves and branches under them. Leaves do not. A leaf can
// transform into a branch if a branch or leaf is added under it.
// For the same of efficiency, any method that transforms a leaf into
// a branch for any reason will automatically discard its value without
// warning.
//
// Types
//
// An empty Node has type of 0 and must display as "[]". Types must have
// both a positive integer and a consistent name or tag to go with it.
// Types as integers are used when printing in short form and provide
// the fastest parsing. Type names are displayed with printing in long
// form.
//
// No nils or panics
//
// No method will ever return a nil value or panic. This simplified the
// calling code using the Node. Empty lists are often returned, however.
type Node struct {
	T int    // type
	V string // value

	tree  *Tree // source of Types, etc.
	up    *Node // branch
	left  *Node // previous
	right *Node // next
	first *Node // first sub
	last  *Node // last sub
	types *[]string
}

// Under returns a slice of references to all the Nodes immediately
// below this node (sometimes called "children"). Never returns nil to
// simplify callers with iteration. Will return an empty list if none.
func (n *Node) Under() []*Node {
	list := []*Node{}
	if n.first == nil {
		return list
	}
	list = append(list, n.first)
	cur := n.first
	for {
		cur = cur.right
		if cur == nil {
			break
		}
		list = append(list, cur)
	}
	return list
}

/*

// GraftLeft adds the passed Node to the left of itself.
func (n *Node) GraftLeft(l *Node) {
	l.B = n.up
	if n.left == nil {
		l.N = n
		n.left = l
		if n.up != nil {
			n.up.F = l
		}
		return
	}
	l.P = n.left
	l.N = n
	n.left.N = l
	n.left = l
}

// GraftRight adds the passed Node to the right of itself.
func (n *Node) GraftRight(r *Node) {
	r.B = n.up
	if n.right == nil {
		r.P = n
		n.right = r
		if n.up != nil {
			n.up.L = r
		}
		return
	}
	r.N = n.right
	r.P = n
	n.right.P = r
	n.right = r
}

// Prune removes itself and grafts everything together to fill void.
func (n *Node) Prune() *Node {
	if n.up != nil {
		if n.up.F == n {
			n.up.F = n.right
		}
		if n.up.L == n {
			n.up.L = n.left
		}
	}
	if n.left != nil {
		n.left.N = n.right
	}
	if n.right != nil {
		n.right.P = n.left
	}
	n.up = nil
	n.right = nil
	n.left = nil
	return n
}

// GraftSelf replace self with a completely new Node
func (n *Node) GraftSelf(c *Node) *Node {
	c.up = n.up
	c.left = n.left
	c.right = n.right
	if n.up.L == n {
		n.up.L = c
	}
	if n.up.F == n {
		n.up.F = c
	}
	if n.left != nil {
		n.left.N = c
	}
	if n.right != nil {
		n.right.P = c
	}
	n.up = nil
	n.right = nil
	n.left = nil
	return n
}

// GraftDown adds the specified branch or leaf Node under itself
// transforming itself into a branch if it was not already.
func (n *Node) GraftDown(c *Node) {
	if n.first == nil {
		c.up = n
		n.first = c
		n.left = c
		return
	}
	n.left.GraftRight(c)
}

// -------------------- PrintAsJSON implementation --------------------

// MarshalJSON fulfills the interface and avoids use of slower
// reflection-based parsing. Nodes must be either branches ([1,[]]) or
// leafs ([1,"foo"]). Branches are allowed to have nothing on them but
// usually have other branches and leaves. This design means that every
// possible Node can be represented by a highly efficient two-element
// array. This MarshalJSON implementation uses the Bonzai json package
// which more closely follows the JSON standard for acceptable string
// data, notably Unicode characters are not escaped and remain readable.
func (n *Node) MarshalJSON() ([]byte, error) {
	list := n.Under()
	if len(list) == 0 {
		if n.V == "" {
			if n.T == 0 {
				return []byte("[]"), nil
			}
			return []byte(fmt.Sprintf(`[%d]`, n.T)), nil
		}
		return []byte(fmt.Sprintf(`[%d,"%v"]`, n.T, json.Escape(n.V))), nil
	}
	byt, _ := list[0].MarshalJSON()
	buf := "[" + string(byt)
	for _, u := range list[1:] {
		byt, _ = u.MarshalJSON() // no error ever returned
		buf += "," + string(byt)
	}
	buf += "]"
	return []byte(fmt.Sprintf(`[%d,%v]`, n.T, buf)), nil
}

// JSONL does not strictly fulfill the PrintAsJSON interface in order to
// preserve the intent of consistent, readable long form JSON. Notably
// different is the use of type names instead of their integer
// equivalents. If the types index is nil then JSON() will be returned
// instead.
func (n *Node) JSONL() string {
	if n.types == nil {
		return n.JSON()
	}
	return n.pretty(0)
}

// called recursively to build the JSONL string
func (n *Node) pretty(depth int) string {
	buf := ""
	indent := strings.Repeat(" ", depth*2)
	depth++
	types := *n.types
	buf += fmt.Sprintf(`%v["%v", `, indent, types[n.T])
	if n.first != nil {
		buf += "[\n"
		under := n.Under()
		for i, c := range under {
			buf += c.pretty(depth)
			if i != len(under)-1 {
				buf += ",\n"
			} else {
				buf += fmt.Sprintf("\n%v]", indent)
			}
		}
		buf += "]"
	} else {
		buf += fmt.Sprintf(`"%v"]`, json.Escape(n.V))
	}
	return buf
}

// JSON implements PrintAsJSON multi-line, 2-space indent JSON output.
func (s *Node) JSON() string { b, _ := s.MarshalJSON(); return string(b) }

// Print implements PrintAsJSON.
func (s *Node) Print() { fmt.Println(s.JSON()) }

// PrintL implements PrintAsJSON.
func (s *Node) PrintL() { fmt.Println(s.JSONL()) }

// Log implements PrintAsJSON.
func (s *Node) Log() { log.Print(s.JSON()) }

// LogL implements PrintAsJSON.
func (s *Node) LogL() { log.Print(s.JSONL()) }

// ----------------------------- old stuff ----------------------------

/*

// IsLeaf returns true if Node has no branch of its own. Note that
// a leaf can transform into a branch once a leaf or branch is added
// under it.
func (n *Node) IsLeaf() bool { return n.first == nil }

// OnBranch returns true if Node is on a branch.
func (n *Node) OnBranch() bool { return n.up == nil }

// IsBranch returns true if Node has any children.
func (n *Node) IsBranch() bool { return n.first != nil }

func (n *Node) AdoptFrom(other *Node) {
	if other.F == nil {
		return
	}
	c := other.F.Prune()
	n.AppendChild(c)
	n.AdoptFrom(other)
}

// Action is a first-class function type used when Visiting each Node.
// The return value will be sent to a channel as each Action completes.
// It can be an error or anything else.
type Action func(n *Node) interface{}

func (n *Node) Visit(act Action, rvals chan interface{}) {
	if rvals == nil {
		act(n)
	} else {
		rvals <- act(n)
	}
	if n.first == nil {
		return
	}
	for _, c := range n.Children() {
		c.Visit(act, rvals)
	}
	return
}

// VisitAsync walks a parent Node and all its Children asynchronously by
// flattening the Node tree into a one-dimensional array and then
// sending each Node its own goroutine Action call. The limit must
// set the maximum number of simultaneous goroutines (which can usually
// be in the thousands) and must be 2 or more or will panic. If the
// channel of return values is not nil it will be sent all return values
// as Actions complete.
func (n *Node) VisitAsync(act Action, lim int, rvals chan interface{}) {
	nodes := []*Node{}

	if lim < 2 {
		panic("visitasync: limit must be 2 or more")
	}

	add := func(node *Node) interface{} {
		nodes = append(nodes, node)
		return nil
	}

	n.Visit(add, nil)

	// use buffered channel to throttle
	sem := make(chan interface{}, lim)
	for _, node := range nodes {
		sem <- true
		if rvals == nil {
			go func(node *Node) {
				defer func() { <-sem }()
				act(node)
			}(node)
			continue
		} else {
			go func(node *Node) {
				defer func() { <-sem }()
				rvals <- act(node)
			}(node)
		}
	}

	// waits for all (keeps filling until full again)
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	// all goroutines have now finished
	if rvals != nil {
		close(rvals)
	}

}

func (n *Node) UnmarshalJSON(b []byte) error {
	p := new(JsonParser)
	p.Init(b)
	err := p.Parse(n)
	if err != nil {
		return err
	}

	return nil
}

*/
