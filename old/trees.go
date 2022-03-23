package scan

import "github.com/rwxrob/structs/tree"

// FIXME These currently do not create the needed Node structures
// properly

// Slice is like "string" but marks the beginning and ending in the
// buffer string data instead of containing the actual data.
type Slice struct {
	Beg *Cur
	End *Cur
}

// Beg begins a new Node containing a Slice at the current location and
// pushes onto Parsing.
func (s *R) Beg(t int) {

	if s.CurNode == nil {
		tr := tree.New[*Slice]([]string{})
		s.CurNode = tr.Root
		s.Trees = append(s.Trees, tr)
	}

	sl := &Slice{Beg: s.Mark()}
	n := s.CurNode.Add(t, sl)
	s.Parsing.Push(n)
}

// End ends a new Node at the current location and pops it off of
// Parsing placing adding it under the node in the current tree. End is
// automatically implied when the end of data is reached for every Beg
// that has not yet been closed.
func (s *R) End() {
	if s.CurNode == nil {
		s.Warn(`no current node open, possibly forgot Beg? at %v`, s.Cur)
		return
	}
	n := s.Parsing.Pop()
	n.V.End = s.Mark()
	if n.P != nil {
		s.CurNode = n.P
	}
}
