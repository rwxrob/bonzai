// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package fn_test

import (
	"fmt"
	"log"
	"os"

	"github.com/rwxrob/bonzai/pkg/fn"
)

func ExampleA_a_is_for_array() {
	fn.A[int]{1, 2, 3}.Print()
	fmt.Println()
	fn.A[string]{"one", "two", "three"}.Println()
	// Output:
	// 123
	// one
	// two
	// three
}

func ExampleA_Each() {
	fn.A[int]{1, 2, 3}.Each(func(i int) { fmt.Print(i) })
	fn.A[int]{1, 2, 3}.E(func(i int) { fmt.Print(i) })
	// Output:
	// 123123
}

func ExampleA_Map() {
	AddOne := func(i int) int { return i + 1 }
	fn.A[int]{1, 2, 3}.Map(AddOne).Print()
	fn.A[int]{1, 2, 3}.M(AddOne).M(AddOne).Print()
	// Output:
	// 234345
}

func ExampleA_Filter() {
	GtTwo := func(i int) bool { return i > 2 }
	LtFour := func(i int) bool { return i < 4 }
	fn.A[int]{1, 2, 3, 4}.Filter(GtTwo).Print()
	fn.A[int]{1, 2, 3, 4}.F(GtTwo).F(LtFour).Print()
	// Output:
	// 343
}

func ExampleA_Reduce() {
	Sum := func(i int, a *int) { *a += i }
	fmt.Println(*fn.A[int]{1, 2}.Reduce(Sum))
	fmt.Println(*fn.A[int]{1, 2, 3, 4, 5}.R(Sum))
	// Output:
	// 3
	// 15
}

func ExampleA_Print() {
	fn.A[int]{1, 2}.Println()
	fn.A[int]{1, 2}.Printf("some: %v\n")
	fn.A[int]{1, 2}.Print()
	// Output:
	// 1
	// 2
	// some: 1
	// some: 2
	// 12
}

func ExampleA_Log() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	fn.A[int]{1, 2}.Log()
	fn.A[int]{1, 2}.Logf("some: %v")
	// Output:
	// 1
	// 2
	// some: 1
	// some: 2
}

func ExamplePipePrint() {
	thing := func(in any) any { return fmt.Sprintf("%v thing", in) }
	other := func(in any) any { return fmt.Sprintf("%v other", in) }
	fn.PipePrint("some", other, thing)
	// Output:
	// some other thing
}

func ExamplePipePrint_error() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	thing := func(in any) any { return fmt.Sprintf("%v thing", in) }
	other := func(in any) any { return fmt.Errorf("bork") }
	fn.PipePrint("some", other, thing)
	// Output:
	// bork
}
