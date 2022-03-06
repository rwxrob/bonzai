package util

type stacked struct {
	val   any
	left  *stacked
	right *stacked
}

// Stack implements a simple stack data structure using a linked list.
type Stack struct {
	first *stacked
	last  *stacked
}

// Push will add an item of any type to the stack in front of the
// others.
func (s *Stack) Push(it any) {
	n := new(stacked)
	n.val = it
	if s.first == nil {
		s.first = n
		s.last = n
		return
	}
	s.last.right = n
	n.left = s.last
	s.last = n
}

// Pop will remove the most recently added item and return it.
func (s *Stack) Pop() any {
	if s.last == nil {
		return nil
	}
	popped := s.last.val
	s.last = s.last.left
	return popped
}

func (s *Stack) Peek() any {
	if s.last == nil {
		return nil
	}
	return s.last.val
}
