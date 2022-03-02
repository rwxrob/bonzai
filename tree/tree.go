// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package tree

import (
	_json "encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/rwxrob/bonzai/json"
)

// Tree is an encapsulating struct to contain the Trunk of Nodes and
// their Types. The New method also ensures that every Node knows about
// the Tree to which it belongs. Most designs will then fetch the
// t.Trunk (root Node) and use that. As with most Bonzai structs, Tree
// implements the PrintAsJSON interface.
type Tree struct {
	Trunk Node
	Types []string
}

const (
	UNKNOWN = 0
)

// New creates and initializes a new tree creating it's Trunk (root
// Node) and assigning it the type of 1 (0 is reserved for UNKNOWN),
// which corresponds to the first Type string.  Calling this constructor
// is preferred over creating Tree references directly since it requires
// specifying the types used in the tree.  Some uses, however, may not
// know the types in advance and need to assign them later to
// Tree.Types. Most will proceed to use the t.Trunk (root node) after
// calling New.
func New(types []string) *Tree {
	t := new(Tree)
	t.Types = []string{"UNKNOWN"}
	t.Types = append(t.Types, types...)
	t.Trunk.T = 1
	t.Trunk.tree = t
	return t
}

// -------------------- Tree PrintAsJSON interface --------------------

// JSON implements PrintAsJSON multi-line, 2-space indent JSON output.
func (s *Tree) JSON() string { b, _ := _json.Marshal(s); return string(b) }

// String implements PrintAsJSON.
func (s Tree) String() string { return s.JSON() }

// Print implements PrintAsJSON.
func (s *Tree) Print() { fmt.Println(s.JSON()) }

// Log implements PrintAsJSON.
func (s *Tree) Log() { log.Print(s.JSON()) }

// ------------------------------- Node -------------------------------

// Nodes are for constructing rooted node trees of typed strings based
// on the tenet of the UNIX philosophy that suggests to focus on
// parsable text above all and converting when needed later. Usually,
// you will start with the Trunk (root Node) of a new Tree so that you
// can specify the types of your Nodes.
//
// Branch or Leaf
//
// A Node can either be a "branch" or a "leaf" but not both. Branches
// have other leaves and branches under them. Leaves do not. A leaf can
// transform into a branch if a branch or leaf is added under it.  For
// the same of efficiency, any method that transforms a leaf into
// a branch for any reason will automatically discard its value without
// warning.
//
// Types
//
// An empty Node has type of 0 and must display as "[]". Types must have
// both a positive integer and a consistent name or tag to go with it.
// A new Tree will always assign the type 1 to the root Node. Types will
// Print as integers when printing in short form and provide the fastest
// parsing. Type names and whitespace are added when PrettyPrint is
// called.
type Node struct {
	T int    // type
	V string // value, zero-ed out when anything added under

	tree  *Tree // source of Types, etc.
	up    *Node // branch
	left  *Node // previous
	right *Node // next
	first *Node // first sub
	last  *Node // last sub
	types *[]string
}

// Branch returns the current branch this Node is on, or nil.
func (n *Node) Branch() *Node { return n.up }

// Left returns the Node to immediate left or nil.
func (n *Node) Left() *Node { return n.left }

// Right returns the Node to immediate right or nil.
func (n *Node) Right() *Node { return n.right }

// FirstUnder returns the first Node under current Node.
func (n *Node) FirstUnder() *Node { return n.first }

// LastUnder returns the last Node under current Node.
func (n *Node) LastUnder() *Node { return n.last }

// AllUnder returns all Nodes under the current Node or nil.
func (n *Node) AllUnder() []*Node {
	if n.first == nil {
		return nil
	}
	cur := n.first
	c := []*Node{cur}
	for {
		cur = cur.right
		if cur == nil {
			break
		}
		c = append(c, cur)
	}
	return c
}

// IsRoot is not currently on any branch even though it might be
// associated still with a given Tree.
func (n *Node) IsRoot() bool { return n.up == nil }

// IsDetached returns true if Node has no attachments to any other Node.
func (n *Node) IsDetached() bool {
	return n.up == nil && n.first == nil &&
		n.last == nil && n.left == nil && n.right == nil
}

// IsLeaf returns true if Node has no branch of its own but does have
// a value. Note that a leaf can transform into a branch once a leaf or
// branch is added under it.
func (n *Node) IsLeaf() bool { return n.first == nil && n.V != "" }

// IsBranch returns true if Node has anything under it at all
func (n *Node) IsBranch() bool { return n.first != nil }

// IsNull returns true if Node has no value and nothing under it but is
// OnBranch.
func (n *Node) IsNull() bool { return n.first == nil && n.V == "" }

// Info logs a summary of the properties of the Node mostly for
// use when debugging. Remember to log.SetOutput(os.Stdout) and
// log.SetFlags(0) when using this in Go example tests.
func (n *Node) Info() {
	log.Printf(`------
Type:       %v
Value:      %q
IsRoot:     %v
IsDetached: %v
IsLeaf:     %v 
IsBranch:   %v 
IsNull:     %v`,
		n.T, n.V, n.IsRoot(), n.IsDetached(),
		n.IsLeaf(), n.IsBranch(), n.IsNull())
}

// ----------------- Node PrintAsJSON interface (plus) ----------------

// MarshalJSON fulfills the interface and avoids use of slower
// reflection-based parsing. Nodes must be either branches ([1,[]]) or
// leafs ([1,"foo"]). Branches are allowed to have nothing on them but
// usually have other branches and leaves. This design means that every
// possible Node can be represented by a highly efficient two-element
// array. This MarshalJSON implementation uses the Bonzai json package
// which more closely follows the JSON standard for acceptable string
// data, notably Unicode characters are not escaped and remain readable.
func (n *Node) MarshalJSON() ([]byte, error) {
	list := n.AllUnder()
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

// PrettyPrint uses type names instead of their integer
// equivalents and adds indentation and whitespace.
func (n *Node) PrettyPrint() {
	fmt.Println(n.pretty(0))
}

// called recursively to build the JSONL string
func (n *Node) pretty(depth int) string {
	buf := ""
	indent := strings.Repeat(" ", depth*2)
	depth++
	buf += fmt.Sprintf(`%v["%v", `, indent, n.tree.Types[n.T])
	if n.first != nil {
		buf += "[\n"
		under := n.AllUnder()
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

// String implements PrintAsJSON and fmt.Stringer interface as JSON.
func (s Node) String() string { return s.JSON() }

// Print implements PrintAsJSON.
func (s *Node) Print() { fmt.Println(s.JSON()) }

// Log implements PrintAsJSON.
func (s Node) Log() { log.Print(s.JSON()) }

// ------------------------------- Nodes ------------------------------

// Spawn spawns new orphan Node for this same Tree
func (n *Node) spawn(i []any) *Node {
	leaf := new(Node)
	switch len(i) {
	case 2:
		leaf.V = i[1].(string)
		fallthrough
	case 1:
		leaf.T = i[0].(int)
	}
	leaf.tree = n.tree
	return leaf
}

// NewRight creates a new Node and grafts it to the right of the current
// one on the same branch. The type and initial value can optionally be
// passed as arguments.
func (n *Node) NewRight(i ...any) *Node {
	leaf := n.spawn(i)
	n.GraftRight(leaf)
	return leaf
}

// NewLeft creates a new Node and grafts it to the left of current one
// on the same branch. The type and initial value can optionally be
// passed as arguments.
func (n *Node) NewLeft(i ...any) *Node {
	leaf := n.spawn(i)
	n.GraftLeft(leaf)
	return leaf
}

// NewUnder creates a new Node and grafts it down below the current one
// adding it to the left of other branches and leaves below. The type
// and initial value can optionally be passed as arguments.
func (n *Node) NewUnder(i ...any) *Node {
	leaf := n.spawn(i)
	n.GraftUnder(leaf)
	return leaf
}

// Graft replaces current node with a completely new Node and returns it
func (n *Node) Graft(c *Node) *Node {
	c.up = n.up
	c.left = n.left
	c.right = n.right
	if n.up.last == n {
		n.up.last = c
	}
	if n.up.first == n {
		n.up.first = c
	}
	if n.left != nil {
		n.left.right = c
	}
	if n.right != nil {
		n.right.left = c
	}
	n.up = nil
	n.right = nil
	n.left = nil
	return n
}

// GraftRight adds existing Node to the right of itself as a peer and
// returns it.
func (n *Node) GraftRight(r *Node) *Node {
	r.up = n.up
	if n.right == nil {
		r.left = n
		n.right = r
		if n.up != nil {
			n.up.last = r
		}
		return r
	}
	r.right = n.right
	r.left = n
	n.right.left = r
	n.right = r
	return r
}

// GraftLeft adds existing Node to the left of itself and returns it.
func (n *Node) GraftLeft(l *Node) *Node {
	l.up = n.up
	if n.left == nil {
		l.right = n
		n.left = l
		if n.up != nil {
			n.up.first = l
		}
		return l
	}
	l.left = n.left
	l.right = n
	n.left.right = l
	n.left = l
	return l
}

// GraftUnder adds existing node under current node to the right of
// others already underneath and returns it.
func (n *Node) GraftUnder(c *Node) *Node {
	c.up = n
	if n.first == nil {
		n.first = c
		n.last = c
		return c
	}
	return n.first.GraftRight(c)
}

// Prune removes and returns itself and grafts everything together to
// fill void.
func (n *Node) Prune() *Node {
	if n.up != nil {
		if n.up.first == n {
			n.up.first = n.right
		}
		if n.up.last == n {
			n.up.last = n.left
		}
	}
	if n.left != nil {
		n.left.right = n.right
	}
	if n.right != nil {
		n.right.left = n.left
	}
	n.up = nil
	n.right = nil
	n.left = nil
	return n
}

// Take takes everything under target Node and adds underneath itself.
func (n *Node) Take(from *Node) {
	if from.first == nil {
		return
	}
	c := from.first.Prune()
	n.GraftUnder(c)
	n.Take(from)
}

/*

/*
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

*/
