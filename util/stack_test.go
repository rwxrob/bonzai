package util_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/util"
)

func ExampleStack() {
	s := new(util.Stack)
	s.Push("some")
	s.Push("another")
	fmt.Println(s.Peek())
	fmt.Println(s.Pop())
	fmt.Println(s.Peek())
	s.Push(1)
	fmt.Println(s.Peek())
	s.Pop()
	s.Pop()
	fmt.Println(s.Peek())
	// Output:
	// another
	// another
	// some
	// 1
	// <nil>
}
