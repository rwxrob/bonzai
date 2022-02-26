package maps

import (
	"path/filepath"
	"sort"

	"github.com/rwxrob/bonzai/fn"
)

// Prefix returns a new slice with prefix added to each string.
func Prefix(in []string, pre string) []string {
	return fn.Map(in, func(i string) string { return pre + i })
}

// Keys returns the keys in lexicographically sorted order.
func Keys[T any](m map[string]T) []string {
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
		sort.Strings(keys)
	}
	return keys
}

// CleanPaths runs filepath.Clean on each item in the slice and returns.
func CleanPaths(paths []string) []string {
	return fn.Map(paths, func(i string) string { return filepath.Clean(i) })
}
