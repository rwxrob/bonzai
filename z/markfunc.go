package Z

import (
	"text/template"

	"github.com/rwxrob/to"
)

// This file contains the BonzaiMark builtins that Cmd authors can use
// in their Description and other places where templated BonzaiMark is
// allowed.

var markFuncMap = template.FuncMap{
	"indent": indent,
}

func indent(n int, in string) string {
	return to.Indented(in, IndentBy+n)
}
