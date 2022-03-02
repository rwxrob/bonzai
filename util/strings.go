package util

import (
	"bufio"
	"fmt"
	"strings"
)

// Lines transforms the input into a string and then divides that string
// up into lines (\r?\n) suitable for functional map operations.
func Lines[T any](in T) []string {
	buf := fmt.Sprintf("%v", in)
	lines := []string{}
	scan := bufio.NewScanner(strings.NewReader(buf))
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}
	return lines
}
