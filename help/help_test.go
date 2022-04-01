package help

import (
	"fmt"
	"testing"
)

func TestRender_bracketed(t *testing.T) {
	oBracketed = "<>"
	cBracketed = "</>"
	out := render([]rune("<bracket - ed>"))
	fmt.Println(out)
}
