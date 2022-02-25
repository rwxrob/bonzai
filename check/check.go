package check

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
