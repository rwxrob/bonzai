// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filter

// Number combines the primitives generally considered numbers by JSON
// and other high-level structure data representations.
type Number interface {
	int16 | int32 | int64 | float32 | float64
}

// Text combines byte slice and string.
type Text interface {
	[]byte | string
}

// P is for "principle" in this case. These are the types that have
// representations in JSON and other high-level structured data
// representations.
type P interface {
	int16 | int32 | int64 | float32 | float64 |
		[]byte | string | bool
}
