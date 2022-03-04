// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package tree

import (
	"encoding/json"
	_json "encoding/json"
	"fmt"
	"log"

	"github.com/rwxrob/bonzai/scan"
	"github.com/rwxrob/bonzai/scan/is"
)

const UNKNOWN = 0

// Tree is an encapsulating struct to contain the Root Node and all
// possible Types for any Node. The tree.New method should be called
// (rather than new(Tree)) to ensure that every Node knows about the
// Tree to which it belongs. Most users of a Tree will make direct use
// of Tree.Root (which can safely be swapped out for a root Node
// reference at any time). Tree implements the PrintAsJSON interface and
// uses custom methods for MarshalJSON and UnmarshalJSON to facilitate
// storage, transfer, and documentation.
type Tree struct {
	Root   *Node
	types  []string
	typesm map[string]int
}

// used to Marshal/Unmarshal
type _Tree struct {
	Root  *Node
	Types []string
}

// SetTypes sets the internal types slice and index map corresponding to
// the given integer. It is normally called from New when creating a new
// Tree.
func (t *Tree) SetTypes(types []string) {
	t.types = []string{"UNKNOWN"}
	t.types = append(t.types, types...)
	t.typesm = map[string]int{"UNKNOWN": 0}
	for n, v := range types {
		t.typesm[v] = n + 1
	}
}

// Seed returns new detached Node from the same Tree
func (t *Tree) Seed(i ...any) *Node {
	leaf := new(Node)
	switch len(i) {
	case 2:
		leaf.V = i[1].(string)
		fallthrough
	case 1:
		leaf.T = i[0].(int)
	}
	leaf.tree = t
	return leaf
}

// New creates and initializes a new tree creating it's Root Node
// reference and assigning it the type of 1 (0 is reserved for UNKNOWN),
// which corresponds to the first Type string (index 0).  Calling this
// constructor is preferred over creating Tree references directly since
// it requires specifying the types used in the tree.  Some uses,
// however, may not know the types in advance and need to assign them
// later to Tree.Types. Most will proceed to use the t.Root after
// calling New.
func New(types []string) *Tree {
	t := new(Tree)
	t.SetTypes(types)
	t.Root = new(Node)
	t.Root.T = 1
	t.Root.tree = t
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

// MarshalJSON implements the json.Marshaler interface to include the
// otherwise private Types list.
func (s *Tree) MarshalJSON() ([]byte, error) {
	return json.Marshal(_Tree{s.Root, s.types})
}

// UnmarshalJSON implements the json.Unmarshaler interface to include
// the otherwise private Types list.
func (s *Tree) UnmarshalJSON(in []byte) error {
	m := new(_Tree)
	if err := json.Unmarshal(in, m); err != nil {
		return err
	}
	s.Root = m.Root
	s.types = m.Types
	return nil
}

// ------------------------------- parse ------------------------------

// Parse takes a string, []byte, or io.Reader of compressed or "pretty"
// JSON data and returns a new tree. See Node.MarshalJSON for more.
func Parse(in any, types []string) (*Tree, error) {
	t := New(types)
	s, err := scan.New(in)
	if err != nil {
		return nil, err
	}

	// shortcuts
	ws := is.Opt{is.WS}
	dq := is.Seq{'\\', '"'}

	// nodes
	jstrc := is.Seq{is.Esc{'\\', '"', is.Ugraphic}}
	jstr := is.Seq{'"', is.Mn1{jstrc}, '"'}
	ntype := is.Seq{is.Not{'"'}, is.Uletter}
	ntyp := is.In{is.Mn1{is.Digit}}
	null := is.Seq{'[', ws, ntyp, ws, ']'}
	leaf := is.Seq{'[', ws, ntyp, ws, jstr, ']'}

	// consume optional initial whitespace
	s.Expect(ws)
	s.Log()

	return t, nil
}
