package is

import (
	"slices"
	"strconv"
	"strings"
)

// AllLatinASCIILowerWithDashes checks if the input string [in] contains only
// lowercase Latin ASCII letters and dashes. It returns false if the string
// is empty, starts with a dash, or ends with a dash, or if it contains any
// character outside the range of 'a' to 'z' or the dash character.
func AllLatinASCIILowerWithDashes(in string) bool {
	if len(in) == 0 || in[0] == '-' || in[len(in)-1] == '-' {
		return false
	}
	for _, r := range in {
		if ('a' <= r && r <= 'z') || r == '-' {
			continue
		}
		return false
	}
	return true
}

// AllLatinASCIILower checks if the input string [in] contains only lowercase
// Latin ASCII letters. It returns false if the string is empty or contains
// any character outside the range of 'a' to 'z'.
func AllLatinASCIILower(in string) bool {
	if len(in) == 0 {
		return false
	}
	for _, r := range in {
		if 'a' <= r && r <= 'z' {
			continue
		}
		return false
	}
	return true
}

// AllLatinASCIIUpper checks if the input string [in] contains only uppercase
// Latin ASCII letters. It returns false if the string is empty or contains
// any character outside the range of 'A' to 'Z'.
func AllLatinASCIIUpper(in string) bool {
	if len(in) == 0 {
		return false
	}
	for _, r := range in {
		if 'A' <= r && r <= 'Z' {
			continue
		}
		return false
	}
	return true
}

func Truthy(val string) bool {
	val = strings.ToLower(strings.TrimSpace(val))
	if slices.Contains([]string{"t", "true", "on"}, val) {
		return true
	}
	if slices.Contains([]string{"f", "false", "off"}, val) {
		return false
	}
	if num, err := strconv.Atoi(val); err == nil {
		return num > 0
	}
	return false
}
