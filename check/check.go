package check

import "reflect"

// Blank checks that the passed string, []byte (or slice of either)
// contains something besides an empty string.
func Blank(i interface{}) bool {
	switch v := i.(type) {
	case nil:
		return false
	case string:
		return v != ""
	case []byte:
		return string(v) != ""
	case []string:
		return v != nil && len(v) != 0 && v[0] != ""
	case [][]byte:
		return v != nil && len(v) != 0 && string(v[0]) != ""
	default:
		panic("cannot check if type is blank")
	}
	return false
}

// IsNil is a shortcut for reflect.ValueOf(foo).IsNil() and should only
// be used when foo == nil is in question, such as whenever the value of
// foo is an interface of any kind. In fact, every interface should use
// this check instead just to be sure to avoid surprise (and extremely
// odd) logic errors. Nil is not "nil" in Go.
func IsNil(i interface{}) bool { return reflect.ValueOf(i).IsNil() }
