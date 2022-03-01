// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package is contains the pseudo-grammar structs recognized by the scan.Expect and scan.Check methods. These structs are guaranteed never to have their structure order change in any way allowing them to be used in key-less, in-line composable notation (which, despite the many editor warnings, is entirely supported by Go and always will be).
*/
package is

// Not represents the logical inverse of whatever is passed. If This
// were a string, for example, Expect/Check would test that it was *not*
// at a given scan location.
type Not struct {
	This interface{}
}

/*
type Min struct {
	Match interface{}
	Min   int
}

type Count struct {
	Match interface{}
	Count int
}

type Seq []interface{}

type OneOf []interface{}

type MinMax struct {
	Match interface{}
	Min   int
	Max   int
}

type Opt struct {
	This interface{}
}
*/
