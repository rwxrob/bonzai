package is

type text interface {
	string | []byte | []rune
}

func AllLatinASCIILower[T text](txt T) bool {
	for _, r := range []rune(string(txt)) {
		if 'a' <= r && r <= 'z' {
			continue
		}
		return false
	}
	return true
}

func AllLatinASCIIUpper[T text](txt T) bool {
	for _, r := range []rune(string(txt)) {
		if 'A' <= r && r <= 'Z' {
			continue
		}
		return false
	}
	return true
}
