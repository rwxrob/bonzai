/*
Package mapf contains nothing but map functions suitable for use with
the fn.Map generic function or the equivalent fn.A method. See the maps
package for functions that accept entire generic slices to be
transformed with mapf (or other) functions.
*/
package mapf

// HashComment adds a "# " prefix.
func HashComment(i string) string { return "# " + i }
