package check

import "reflect"

// IsNil is a shortcut for reflect.ValueOf(foo).IsNil() and should only
// be used when foo == nil is in question, such as whenever the value of
// foo is an interface of any kind. In fact, every interface should use
// this check instead just to be sure to avoid surprise (and extremely
// odd) logic errors. Nil is not "nil" in Go.
func IsNil(i interface{}) bool { return reflect.ValueOf(i).IsNil() }
