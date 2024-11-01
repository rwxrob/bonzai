package qstack_test

import (
	"fmt"
	"log"
	"os"

	"github.com/rwxrob/bonzai/ds/qstack"
)

func ExampleQS_Push() {
	s := qstack.New[any]()
	fmt.Println(s.Len)
	s.Print()
	s.Push("one")
	fmt.Println(s.Len)
	s.Print()
	s.Push("two")
	fmt.Println(s.Len)
	s.Print()
	// Output:
	// 0
	// []
	// 1
	// ["one"]
	// 2
	// ["one","two"]
}

func ExampleQS_Peek() {
	s := qstack.New[any]()
	s.Print()
	s.Push("it")
	s.Print()
	fmt.Println(s.Peek())
	// Output:
	// []
	// ["it"]
	// it
}

func ExampleQS_references() {
	type Some struct {
		Val string
	}

	s := qstack.New[*Some]()

	// just one
	x := &Some{"some"}
	s.Push(x)
	y := s.Peek()
	// log.Printf("Is %p == %p?", x, y)
	fmt.Println(x == y)

	// two
	xx := &Some{"other"}
	s.Push(xx)
	yy := s.Peek()
	// log.Printf("Is %p == %p?", xx, yy)
	fmt.Println(xx == yy)

	// Output:
	// true
	// true
}

func ExampleQS_Pop() {
	s := qstack.New[any]()
	s.Print()
	p := s.Pop()
	fmt.Println(p)
	s.Push("it")
	s.Push("again")
	s.Print()
	fmt.Println(s.Len)
	p = s.Pop()
	s.Print()
	fmt.Println(p)
	fmt.Println(s.Len)
	// Output:
	// []
	// <nil>
	// ["it","again"]
	// 2
	// ["it"]
	// again
	// 1
}

func ExampleQS_Items() {
	s := qstack.New[any]()
	s.Push(1)
	s.Push(true)
	s.Push("foo")
	s.Push(map[string]int{"ten": 10, "twenty": 20})
	s.Print()
	// Output:
	// [1,true,"foo",{"ten":10,"twenty":20}]
}

func ExampleQS_Shift() {
	s := qstack.New[any]()
	s.Push(1)
	s.Push(true)
	s.Push("foo")
	v := s.Shift()
	fmt.Println(v, s.Len)
	s.Print()
	v = s.Shift()
	fmt.Println(v, s.Len)
	s.Print()
	v = s.Shift()
	fmt.Println(v, s.Len)
	s.Print()
	// Output:
	// 1 2
	// [true,"foo"]
	// true 1
	// ["foo"]
	// foo 0
	// []
}

func ExampleQS_Unshift() {
	s := qstack.New[any]()
	s.Push(1)
	s.Push(true)
	s.Push("foo")
	s.Print()
	fmt.Println(s.Len)
	s.Unshift(0, 34, 2)
	s.Print()
	fmt.Println(s.Len)
	// Output:
	// [1,true,"foo"]
	// 3
	// [0,34,2,1,true,"foo"]
	// 6
}

func ExampleQS_invalid_JSON_Types() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	// QS can be used to store any type,
	// but log an error (no panic) when
	// attempting to use the stack item
	// in a string context.

	s := qstack.New[any]()
	s.Push(func() {})
	s.Print()

	// Output:
	// json: unsupported type: func()
}

func ExampleQS_Copy() {
	s := qstack.New[any]()
	s.Push("some")
	s.Push("thing")

	c := s.Copy()
	c.Pop()

	s.Print()
	c.Print()

	// Output:
	// ["some","thing"]
	// ["some"]
}

func ExampleQS_Scan() {
	s := qstack.New[any]()
	s.Push("foo")
	s.Push("bar")
	s.Push("and")
	s.Push("one")

	for s.Scan() {
		fmt.Print(s.Current())
	}

	// Output:
	// foobarandone
}
