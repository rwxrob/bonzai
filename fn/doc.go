// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package fn standardizes functional types and traditional map/filter/reduce/each functions so that type switching on them can be done consistently. It also provides a collection of common type implementations.

# func(in any) any

Rather that force casting in order for type switches to work, the func(in any) any is left as is and assumed to be used for any type switches checking for map functions. In other words, the func(in any) any is itself a first-class type in this package and those using it should prefer it for all map function implementations.
*/
package fn
