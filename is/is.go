package is

func AllLatinASCIILower(in string) bool {
	for _, r := range in {
		if 'a' <= r && r <= 'z' {
			continue
		}
		return false
	}
	return true
}

func AllLatinASCIIUpper(in string) bool {
	for _, r := range in {
		if 'A' <= r && r <= 'Z' {
			continue
		}
		return false
	}
	return true
}
