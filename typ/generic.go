// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package generic

// Number combines the primitives generally considered numbers by JSON
// and other high-level structure data representations.
type Number interface {
	int | int64 | int32 | int16 | int8 |
		uint64 | uint32 | uint16 | uint8 |
		float64 | float32
}

// Text combines byte slice and string.
type Text interface {
	[]byte | string
}

// Sharable are the types that have representations in JSON, YAML, TOML
// and other high-level structured data representations.
type Sharable interface {
	int | int64 | int32 | int16 | int8 |
		uint64 | uint32 | uint16 | uint8 |
		float64 | float32 |
		[]byte | string |
		bool
}
