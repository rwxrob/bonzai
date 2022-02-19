/*
Copyright 2022 Robert S. Muhlestein.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package filter

import "strings"

// HasPrefix filters the Text input set and returns only those elements
// that have the give prefix.
func HasPrefix[T Text](set []T, pre string) []T {
	m := []T{}
	for _, i := range set {
		if strings.HasPrefix(string(i), pre) {
			m = append(m, i)
		}
	}
	return m
}

// Minus performs a set "minus" operation by returning a new set with
// the elements of the second set removed from it.
func Minus[T Text, M Text](set []T, min []M) []T {
	m := []T{}
	for _, i := range set {
		var seen bool
		for _, n := range min {
			if string(n) == string(i) {
				seen = true
				break
			}
		}
		if !seen {
			m = append(m, i)
		}
	}
	return m
}
